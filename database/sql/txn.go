package goxSql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base/v2/errors"
	"github.com/google/uuid"
	"sync"
)

const txnKey = "__SQLCX_TXN_KEY__"

// TxnBeginner is used to Begin a transaction
type TxnBeginner interface {
	Begin() (Tx, error)
}

//go:generate mockgen -source=txn.go -destination=./mock_txn.go -package=goxSql
type Tx interface {
	// Commit commits the transaction.
	// It returns an error if the commit fails.
	Commit() error

	// Rollback rolls back the transaction.
	// It returns an error if the rollback fails.
	Rollback() error

	// Exec executes a SQL statement within the transaction.
	// The first argument is the SQL query, and the second argument (optional) is any parameters to be used in the query.
	// It returns a sql.Result object and an error.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// Prepare prepares a SQL statement for execution within the transaction.
	// It returns a sql.Stmt object and an error.
	Prepare(query string) (*sql.Stmt, error)

	// PrepareContext prepares a SQL statement for execution within the transaction, using the provided context.
	// It returns a sql.Stmt object and an error.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// Query executes a SQL query within the transaction.
	// The first argument is the SQL query, and the second argument (optional) is any parameters to be used in the query.
	// It returns a sql.Rows object and an error.
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// QueryContext executes a SQL query within the transaction, using the provided context.
	// It returns a sql.Rows object and an error.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow executes a SQL query within the transaction and returns a sql.Row object.
	QueryRow(query string, args ...interface{}) *sql.Row

	// QueryRowContext executes a SQL query within the transaction, using the provided context, and returns a sql.Row object.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Stmt prepares a SQL statement for execution within the transaction using an existing sql.Stmt object.
	// It returns a sql.Stmt object.
	Stmt(stmt *sql.Stmt) *sql.Stmt
}

type txImpl struct {
	Tx

	name     string
	uniqueId string

	topLevelTx *txImpl
	isChild    bool

	mu             *sync.Mutex
	rollbackTx     *txImpl
	isCommitCalled bool
}

func NewTxExt(tx Tx) Tx {
	return &txImpl{
		Tx:       tx,
		uniqueId: uuid.NewString(),
		mu:       &sync.Mutex{},
	}
}

func (tx *txImpl) WithName(name string) Tx {
	tx.name = name
	return tx
}

func (tx *txImpl) Commit() (err error) {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.isChild {
		tx.isCommitCalled = true
	} else {
		// If this is a parent then we do the final commit
		// However, if we have seen a rollback in a child then we return error
		if tx.rollbackTx == nil {
			err = tx.Tx.Commit()
		} else {
			err = &ErrCommitFailedDueToChildTxnFailed{
				Tx:             tx,
				ChildFailedTxn: tx.rollbackTx,
			}
		}
	}
	return err
}

func (tx *txImpl) Rollback() (err error) {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.isChild {
		if !tx.isCommitCalled {
			tx.topLevelTx.rollbackTx = tx
		}
	} else {
		err = tx.Tx.Rollback()
	}
	return err
}

func (tx *txImpl) String() string {
	if tx == nil {
		return "{nil transaction}"
	}
	return fmt.Sprintf("{Name=%s, Child=%t UniqueId=%s}", tx.name, tx.isChild, tx.uniqueId)
}

type TxBeginOptions struct {
	TxnBeginner                 TxnBeginner
	Name                        string
	ContinueExistingTxnIfExists bool
}

func Begin(ctx context.Context, options TxBeginOptions) (context.Context, Tx, error) {
	if options.ContinueExistingTxnIfExists && ctx != nil && ctx.Value(txnKey) != nil {
		if tx, ok := ctx.Value(txnKey).(*txImpl); ok && !tx.isCommitCalled {
			return ctx,
				&txImpl{
					Tx:         tx.Tx,
					uniqueId:   tx.uniqueId,
					name:       options.Name,
					isChild:    true,
					mu:         tx.mu,
					topLevelTx: tx,
				},
				nil
		}
	}

	if options.TxnBeginner == nil {
		return ctx, nil, errors.New("missing TxnBeginner in options")
	}

	if tx, err := options.TxnBeginner.Begin(); err == nil {
		t := NewTxExt(tx).(*txImpl)
		t.isChild = false
		txWrapper := t.WithName(options.Name)
		ctx = context.WithValue(ctx, txnKey, txWrapper)
		return ctx, txWrapper, nil
	} else {
		return ctx, nil, err
	}
}

type ErrCommitFailedDueToChildTxnFailed struct {
	Tx             *txImpl
	ChildFailedTxn *txImpl
}

func (e ErrCommitFailedDueToChildTxnFailed) Error() string {
	return fmt.Sprintf("commit failed - some child txn has failed: parentTxn=%s, failedTxn=%s", e.Tx.String(), e.ChildFailedTxn.String())
}
