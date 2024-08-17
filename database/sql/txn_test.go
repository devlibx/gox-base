package goxSql

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base/v2/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommit(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("Run when commit is called with no error", func(t *testing.T) {
		txMock := NewMockTx(ctrl)
		txMock.EXPECT().Commit().Return(nil).Times(1)
		tx := NewTxExt(txMock)
		assert.NoError(t, tx.Commit())
	})

	t.Run("Run when commit is called with error", func(t *testing.T) {
		txMock := NewMockTx(ctrl)
		txMock.EXPECT().Commit().Return(errors.New("bad error")).Times(1)
		tx := NewTxExt(txMock)
		assert.Error(t, tx.Commit())
	})
}

func TestCommitForRecursiveTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	txMock := NewMockTx(ctrl)
	rollbackCallCount, commitCallCount := 0, 0
	txMock.EXPECT().Commit().DoAndReturn(func() error {
		commitCallCount++
		return nil
	})
	txMock.EXPECT().Rollback().DoAndReturn(func() error {
		rollbackCallCount++
		return nil
	})

	// Child function to test transaction
	child := func(ctx context.Context, t *testing.T, txMock Tx) {
		// Begin a transaction
		ctx, tx, err := Begin(ctx, TxBeginOptions{TxnBeginner: &tb{tx: txMock, err: nil}, Name: "child", ContinueExistingTxnIfExists: true})
		assert.NoError(t, err)
		defer tx.Rollback()

		// Commit a transaction
		err = tx.Commit()
		assert.NoError(t, err)
	}

	// Parent function to test transaction
	parent := func(ctx context.Context, t *testing.T, txMock Tx) {
		// Begin a transaction
		ctx, tx, err := Begin(ctx, TxBeginOptions{TxnBeginner: &tb{tx: txMock, err: nil}, Name: "parent", ContinueExistingTxnIfExists: true})
		assert.NoError(t, err)
		defer tx.Rollback()

		// Make a recursive call to child to test a child call
		child(ctx, t, txMock)

		// Commit a transaction
		err = tx.Commit()
		assert.NoError(t, err)
	}

	ctx := context.Background()
	parent(ctx, t, txMock)
	assert.Equal(t, 1, commitCallCount, "Commit should be called only once")
	assert.Equal(t, 1, rollbackCallCount, "Rollback should be called only once - don't worry we do it it deffer so no impact")
}

func TestCommitForRecursiveTxWithErrorInChild(t *testing.T) {
	ctrl := gomock.NewController(t)
	txMock := NewMockTx(ctrl)

	// We should not get a commit call, only one rollback call
	rollbackCallCount := 0
	txMock.EXPECT().Commit().Times(0)
	txMock.EXPECT().Rollback().DoAndReturn(func() error {
		rollbackCallCount++
		return nil
	})

	// Child function to test transaction
	child := func(ctx context.Context, t *testing.T, txMock Tx) {
		// Begin a transaction
		ctx, tx, err := Begin(ctx, TxBeginOptions{TxnBeginner: &tb{tx: txMock, err: nil}, Name: "child", ContinueExistingTxnIfExists: true})
		assert.NoError(t, err)
		defer tx.Rollback()

		// NOTE - here child did not committed
		// Commit a transaction
		// err = tx.Commit()
		// assert.NoError(t, err)
	}

	// Parent function to test transaction
	parent := func(ctx context.Context, t *testing.T, txMock Tx) {
		// Begin a transaction
		ctx, tx, err := Begin(ctx, TxBeginOptions{TxnBeginner: &tb{tx: txMock, err: nil}, Name: "parent", ContinueExistingTxnIfExists: true})
		assert.NoError(t, err)
		defer tx.Rollback()

		// Make a recursive call to child to test a child call
		child(ctx, t, txMock)

		// Commit a transaction
		err = tx.Commit()
		assert.Error(t, err)
		e, ok := err.(*ErrCommitFailedDueToChildTxnFailed)
		assert.True(t, ok)
		fmt.Println(e.Error())
	}

	ctx := context.Background()
	parent(ctx, t, txMock)
	assert.Equal(t, 1, rollbackCallCount, "Rollback should be called only once - don't worry we do it it deffer so no impact")
}

func TestCommitForRecursiveTxWithErrorInChildButChildOptedForNewTxn(t *testing.T) {
	ctrl := gomock.NewController(t)
	txMock := NewMockTx(ctrl)

	// We should not get a commit call, only one rollback call
	txMock.EXPECT().Commit().Return(nil).Times(1)
	txMock.EXPECT().Rollback().Return(nil).Times(2)

	// Child function to test transaction
	child := func(ctx context.Context, t *testing.T, txMock Tx) {
		// Begin a transaction
		ctx, tx, err := Begin(ctx, TxBeginOptions{TxnBeginner: &tb{tx: txMock, err: nil}, Name: "child", ContinueExistingTxnIfExists: false})
		assert.NoError(t, err)
		defer tx.Rollback()

		// NOTE - here child did not committed
		// Commit a transaction
		// err = tx.Commit()
		// assert.NoError(t, err)
	}

	// Parent function to test transaction
	// Note - child opted to "Not continue with existing Txn" so parent will not fail even if child failed
	parentContinueToCommitEvenIfChildFailed := func(ctx context.Context, t *testing.T, txMock Tx) {
		// Begin a transaction
		ctx, tx, err := Begin(ctx, TxBeginOptions{TxnBeginner: &tb{tx: txMock, err: nil}, Name: "parent", ContinueExistingTxnIfExists: true})
		assert.NoError(t, err)
		defer tx.Rollback()

		// Make a recursive call to child to test a child call
		child(ctx, t, txMock)

		// Commit a transaction
		err = tx.Commit()
		assert.NoError(t, err)
	}

	ctx := context.Background()
	parentContinueToCommitEvenIfChildFailed(ctx, t, txMock)
}

type tb struct {
	tx  Tx
	err error
}

func (t *tb) Begin() (Tx, error) {
	return t.tx, t.err
}
