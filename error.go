package migrate

import (
	"errors"
	"fmt"
)

var (
	ErrNoNewMigrations = errors.New("there are no any migrations")
)

// lock error

// NewLockError wraps a message in a LockError object.
func NewLockError(m string, err error) *LockError {
	return &LockError{
		err:     err,
		message: m,
	}
}

// LockError represents an error returned when a lock is unable to be obtained or freed.
type LockError struct {
	err     error
	message string
}

// Error returns the error message as a string.
func (e *LockError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s (%s)", e.message, e.err.Error())
	}

	return e.message
}

// migration name error

// NewMigrationNameError wraps a message in a MigrationNameError object.
func NewMigrationNameError(m, n string) *MigrationNameError {
	return &MigrationNameError{
		message: m,
		name:    n,
	}
}

// MigrationNameError represents an error returned when a migration has an invalid name.
type MigrationNameError struct {
	message string
	name    string
}

// Error returns the error message as a string.
func (e *MigrationNameError) Error() string {
	return fmt.Sprintf("%s: %q", e.message, e.name)
}

// rollback error

// NewRollbackError wraps a message in a RollbackError object.
func NewRollbackError(rollback, top error) *RollbackError {
	return &RollbackError{
		rollback: rollback,
		top:      top,
	}
}

// RollbackError represents an error returned when a migration has an invalid name.
type RollbackError struct {
	rollback error
	top      error
}

// Error returns the error message as a string.
func (e *RollbackError) Error() string {
	if e.rollback != nil {
		return fmt.Sprintf("%s; %s", e.rollback.Error(), e.top.Error())
	}

	return e.top.Error()
}
