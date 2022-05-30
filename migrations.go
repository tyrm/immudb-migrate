package migrate

import (
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type MigrationsOption func(m *Migrations)

func NewMigrations(opts ...MigrationsOption) *Migrations {
	m := new(Migrations)
	for _, opt := range opts {
		opt(m)
	}
	m.implicitDirectory = filepath.Dir(migrationFile())

	return m
}

type Migrations struct {
	ms MigrationSlice

	implicitDirectory string
}

func (m *Migrations) Add(migration Migration) {
	if migration.Name == "" {
		panic("migrate name is required")
	}
	m.ms = append(m.ms, migration)
}

func (m *Migrations) Register(up MigrationHandler) error {
	fpath := migrationFile()
	name, err := extractMigrationName(fpath)
	if err != nil {
		return err
	}

	m.Add(Migration{
		Name: name,
		Up:   up,
	})

	return nil
}

func (m *Migrations) Sorted() MigrationSlice {
	migrations := make(MigrationSlice, len(m.ms))
	copy(migrations, m.ms)
	sortAsc(migrations)

	return migrations
}

func migrationFile() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	for {
		f, ok := frames.Next()
		if !ok {
			break
		}
		if !strings.Contains(f.Function, "/immudb-migrate.") {
			return f.File
		}
	}

	return ""
}

var fnameRE = regexp.MustCompile(`^(\d{14})_[0-9a-z_\-]+\.`)

func extractMigrationName(fpath string) (string, error) {
	fname := filepath.Base(fpath)

	matches := fnameRE.FindStringSubmatch(fname)
	if matches == nil {
		return "", NewMigrationNameError("unsupported migrate name format", fname)
	}

	return matches[1], nil
}
