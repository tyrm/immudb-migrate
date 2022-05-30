package migrate

import (
	"errors"
	"fmt"
	"testing"
)

//revive:disable:add-constant

var (
	testErr1         = errors.New("test error 1")          // nolint
	testRollbackErr1 = errors.New("rollback test error 1") // nolint
	testTopErr1      = errors.New("top test error 1")      // nolint
	testTopErr2      = errors.New("top test error 2")      // nolint
)

func TestNewLockError(t *testing.T) {
	t.Parallel()

	tables := []struct {
		message string
		err     error

		output LockError
	}{
		{"test message 1", testErr1, LockError{message: "test message 1", err: testErr1}},
		{"test message 2", nil, LockError{message: "test message 2", err: nil}},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running NewLockError", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			newErr := NewLockError(table.message, table.err)

			if newErr.message != table.message {
				t.Errorf("[%d] invalid message, got: '%s', want: '%s'", i, newErr.message, table.message)
			}
			if !errors.Is(newErr.err, table.err) {
				t.Errorf("[%d] invalid error, got: '%v', want: '%v'", i, newErr.err, table.err)
			}
		})
	}
}

func TestLockError_Error(t *testing.T) {
	t.Parallel()

	tables := []struct {
		message string
		err     error

		output string
	}{
		{"test message 1", testErr1, "test message 1 (test error 1)"},
		{"test message 2", nil, "test message 2"},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running LockError Error", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			newErr := NewLockError(table.message, table.err)

			if newErr.Error() != table.output {
				t.Errorf("[%d] invalid error string, got: '%s', want: '%s'", i, newErr.Error(), table.output)
			}
		})
	}
}

func TestNewMigrationNameError(t *testing.T) {
	t.Parallel()

	tables := []struct {
		message string
		name    string

		output MigrationNameError
	}{
		{"test message 1", "name 1", MigrationNameError{message: "test message 1", name: "name 1"}},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running NewMigrationNameError", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			newErr := NewMigrationNameError(table.message, table.name)

			if newErr.message != table.message {
				t.Errorf("[%d] invalid message, got: '%s', want: '%s'", i, newErr.message, table.message)
			}
			if newErr.name != table.name {
				t.Errorf("[%d] invalid name, got: '%s', want: '%s'", i, newErr.message, table.message)
			}
		})
	}
}

func TestMigrationNameError_Error(t *testing.T) {
	t.Parallel()

	tables := []struct {
		message string
		name    string

		output string
	}{
		{"test message 1", "name 1", "test message 1: \"name 1\""},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running MigrationNameError Error", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			newErr := NewMigrationNameError(table.message, table.name)

			if newErr.Error() != table.output {
				t.Errorf("[%d] invalid error string, got: '%s', want: '%s'", i, newErr.Error(), table.output)
			}
		})
	}
}

func TestNewRollbackError(t *testing.T) {
	t.Parallel()

	tables := []struct {
		rollback error
		top      error

		output RollbackError
	}{
		{testRollbackErr1, testTopErr1, RollbackError{rollback: testRollbackErr1, top: testTopErr1}},
		{nil, testTopErr2, RollbackError{rollback: nil, top: testTopErr2}},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running NewMigrationNameError", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			newErr := NewRollbackError(table.rollback, table.top)

			if !errors.Is(newErr.rollback, table.rollback) {
				t.Errorf("[%d] invalid rollback error, got: '%v', want: '%v'", i, newErr.rollback, table.rollback)
			}
			if !errors.Is(newErr.top, table.top) {
				t.Errorf("[%d] invalid top error, got: '%v', want: '%v'", i, newErr.top, table.top)
			}
		})
	}
}

func TestRollbackError_Error(t *testing.T) {
	t.Parallel()

	tables := []struct {
		rollback error
		top      error

		output string
	}{
		{testRollbackErr1, testTopErr1, "rollback test error 1; top test error 1"},
		{nil, testTopErr2, "top test error 2"},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running MigrationNameError Error", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			newErr := NewRollbackError(table.rollback, table.top)

			if newErr.Error() != table.output {
				t.Errorf("[%d] invalid error string, got: '%s', want: '%s'", i, newErr.Error(), table.output)
			}
		})
	}
}

//revive:disable:add-constant
