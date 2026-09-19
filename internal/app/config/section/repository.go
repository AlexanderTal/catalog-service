package section

import "time"

type (
	Repository struct {
		Postgres RepositoryPostgres
	}

	RepositoryPostgres struct {
		Address        string        `required:"true"`
		Username       string        `required:"true"`
		Password       string        `required:"true"`
		Name           string        `required:"true"`
		ReadTimeout    time.Duration `split_words:"true" default:"30s"`
		WriteTimeout   time.Duration `split_words:"true" default:"30s"`
		MigrationTable string        `split_words:"true" default:"schema_migrations"`
	}
)
