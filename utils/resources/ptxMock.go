package resources

import (
	"database/sql"
)

type PtxMock struct {
	tx *sql.Tx
}

func (p *PtxMock) Commit() error {
	return nil
}

func (p *PtxMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.tx.Exec(query, args...)
}

func (p *PtxMock) Prepare(query string) (*sql.Stmt, error) {
	return p.tx.Prepare(query)
}

func (p *PtxMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.tx.Query(query, args...)
}

func (p *PtxMock) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.tx.QueryRow(query, args...)
}

func (p *PtxMock) Rollback() error {
	return nil
}

func (p *PtxMock) Stmt(stmt *sql.Stmt) *sql.Stmt {
	return p.tx.Stmt(stmt)
}

func (p *PtxMock) TearDown() error {
	return p.tx.Rollback()
}
