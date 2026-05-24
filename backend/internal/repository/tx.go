package repository

import (
	"errors"

	"gorm.io/gorm"
)

var defaultDB *gorm.DB

// ErrTxNotReady is returned by Tx when SetDefaultDB hasn't been called.
var ErrTxNotReady = errors.New("repository: Tx called before SetDefaultDB")

// SetDefaultDB stores the shared DB connection for use by package-level Tx.
// Called once at startup.
func SetDefaultDB(db *gorm.DB) {
	defaultDB = db
}

// Tx executes fn within a database transaction, passing a Repository whose
// repos all share the same tx. Returns ErrTxNotReady if SetDefaultDB wasn't
// called (e.g. in tests). Any error from fn triggers a rollback.
func Tx(fn func(txRepo *Repository) error) error {
	if defaultDB == nil {
		return ErrTxNotReady
	}
	tx := defaultDB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	txRepo := New(tx)
	if err := fn(txRepo); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
