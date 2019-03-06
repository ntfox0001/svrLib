package database

import (
	"database/sql"
	"errors"
)

type Transcation struct {
	db *sql.DB
	tx *sql.Tx
}

func newTranscation(db *sql.DB) *Transcation {
	return &Transcation{
		db: db,
	}
}

func (t *Transcation) Start() error {
	if tx, err := t.db.Begin(); err == nil {
		t.tx = tx
		return nil
	} else {
		return err
	}
}

func (t *Transcation) Commit() error {
	if t.tx == nil {
		return errors.New("Need start transcation.")
	}
	err := t.tx.Commit()
	t.tx = nil
	return err
}

func (t *Transcation) Rollback() error {
	if t.tx == nil {
		return errors.New("Need start transcation.")
	}
	err := t.tx.Rollback()
	t.tx = nil
	return err
}

func (t *Transcation) SyncExecOperation(opt IOperation) (*DataResult, error) {
	if t.tx == nil {
		return nil, errors.New("Need start transcation.")
	}
	rt := opt.callData(t.tx)

	return rt, nil
}
