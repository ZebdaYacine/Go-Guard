package usecase

import (
	"context"
	"go-gaurd/feature/profile/domain"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ProfileUseCase struct {
	repo *domain.ProfileRepository // Change to pointer to struct
}

type ProfileUseCaseInterface interface {
	GetAccount(ctx context.Context, query Query) Result
	UpdateAccount(ctx context.Context, query Query) Result
	ActiveAccount(ctx context.Context, query Query) Result
}

func NewProfileUseCase(repo *domain.ProfileRepository) *ProfileUseCase {
	return &ProfileUseCase{
		repo: repo,
	}
}

// ActiveAccount implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) ActiveAccount(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// GetAccount implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) GetAccount(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// UpdateAccount implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) UpdateAccount(ctx context.Context, query Query) Result {
	panic("unimplemented")
}
