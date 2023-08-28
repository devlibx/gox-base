package queue

import (
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

type mySqlStore struct {
	db     *sql.DB
	config queue.MySqlBackedStoreBackendConfig

	initOnce  *sync.Once
	initErr   error
	closeOnce *sync.Once
	closed    bool
}

func (m *mySqlStore) RewriteQuery(input string) string {
	return input
}

func (m *mySqlStore) Init() error {
	m.initOnce.Do(func() {

		// Open connection to DB
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.config.User, m.config.Password, m.config.Host, m.config.Port, m.config.Database)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			err = errors.Wrap(err, "failed to open the db", err)
			m.initErr = err
			return
		}

		m.db = db
		m.db.SetMaxOpenConns(m.config.MaxConnection)
	})

	return m.initErr
}

func (m *mySqlStore) GetSqlDb() (*sql.DB, error) {
	if m.closed {
		return nil, errors.New("store is not initialized or it is already closed")
	}
	return m.db, nil
}

func (m *mySqlStore) Close() error {
	m.closeOnce.Do(func() {

	})
	return nil
}

func NewMySqlBackedStore(config queue.MySqlBackedStoreBackendConfig, init bool) (*mySqlStore, error) {
	m := &mySqlStore{
		config:    config,
		initOnce:  &sync.Once{},
		closeOnce: &sync.Once{},
	}

	if init {
		if err := m.Init(); err != nil {
			return nil, err
		}
	}

	return m, nil
}
