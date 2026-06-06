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
	GetProfile(ctx context.Context, query Query) Result
	UpdateProfile(ctx context.Context, query Query) Result
	UpdateProfilePicture(ctx context.Context, query Query) Result
	UpdatePassword(ctx context.Context, query Query) Result
	ActiveProfile(ctx context.Context, query Query) Result
	Logout(ctx context.Context, query Query) Result
	RefreshAccessToken(ctx context.Context, query Query) Result
}

func NewProfileRepository(db *database.Database) ProfileRepositoryInterface {
	return &ProfileRepository{DB: db.DB}
}

// ActiveProfile implements [ProfileRepositoryInterface].
func (p *ProfileRepository) ActiveProfile(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// GetProfile implements [ProfileRepositoryInterface].
func (p *ProfileRepository) GetProfile(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// Logout implements [ProfileRepositoryInterface].
func (p *ProfileRepository) Logout(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// RefreshAccessToken implements [ProfileRepositoryInterface].
func (p *ProfileRepository) RefreshAccessToken(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// UpdatePassword implements [ProfileRepositoryInterface].
func (p *ProfileRepository) UpdatePassword(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// UpdateProfile implements [ProfileRepositoryInterface].
func (p *ProfileRepository) UpdateProfile(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// UpdateProfilePicture implements [ProfileRepositoryInterface].
func (p *ProfileRepository) UpdateProfilePicture(ctx context.Context, query Query) Result {
	panic("unimplemented")
}
