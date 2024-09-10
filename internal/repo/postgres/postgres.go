package postgres

import (
	"github.com/Slava02/Involvio/internal/interface"
	"github.com/Slava02/Involvio/pkg/postgres"
)

// TranslationRepo -.
type Repo struct {
	*postgres.Postgres
}

var _ _interface.Repo = (*Repo)(nil)

// New -.
func New(pg *postgres.Postgres) *Repo {
	return &Repo{pg}
}

// https://github.com/Masterminds/squirrel
