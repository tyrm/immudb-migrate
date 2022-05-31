package migrate

type migrationConfig struct {
	nop bool
}

func newMigrationConfig(opts []MigrationOption) *migrationConfig {
	cfg := new(migrationConfig)
	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

type MigrationOption func(cfg *migrationConfig)

func WithNopMigration() MigrationOption {
	return func(cfg *migrationConfig) {
		cfg.nop = true
	}
}
