package resources

import (
	"database/sql"
)

type Ptx struct {
	tx *sql.Tx
}

func (p *Ptx) Commit() error {
	return p.tx.Commit()
}

func (p *Ptx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.tx.Exec(query, args...)
}

func (p *Ptx) Prepare(query string) (*sql.Stmt, error) {
	return p.tx.Prepare(query)
}

func (p *Ptx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.tx.Query(query, args...)
}

func (p *Ptx) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.tx.QueryRow(query, args...)
}

func (p *Ptx) Rollback() error {
	return p.tx.Rollback()
}

func (p *Ptx) Stmt(stmt *sql.Stmt) *sql.Stmt {
	return p.tx.Stmt(stmt)
}

func (p *Ptx) TearDown() error {
	// Stub method for tests
	return nil
}
