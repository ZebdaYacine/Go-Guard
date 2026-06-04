package domain

import (
	"context"
	"database/sql"
	"go-gaurd/database"
)

type ProfileRepository struct {
	DB *sql.DB
}

type ProfileRepositoryInterface interface {
	GetAccount(ctx context.Context, query Query) Result
	UpdateAccount(ctx context.Context, query Query) Result
	ActiveAccount(ctx context.Context, query Query) Result
}

func NewProfileRepository(db *database.Database) *ProfileRepository {
	return &ProfileRepository{DB: db.DB}
}

// ActiveAccount implements [ProfileRepositoryInterface].
func (p *ProfileRepository) ActiveAccount(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// GetAccount implements [ProfileRepositoryInterface].
func (p *ProfileRepository) GetAccount(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// UpdateAccount implements [ProfileRepositoryInterface].
func (p *ProfileRepository) UpdateAccount(ctx context.Context, query Query) Result {
	panic("unimplemented")
}
