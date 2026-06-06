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
	repo domain.ProfileRepositoryInterface // Change to pointer to struct
}

type ProfileUseCaseInterface interface {
	GetProfile(ctx context.Context, query Query) Result
	UpdateProfile(ctx context.Context, query Query) Result
	UpdateProfilePicture(ctx context.Context, query Query) Result
	UpdatePassword(ctx context.Context, query Query) Result
	ActiveProfile(ctx context.Context, query Query) Result
	Logout(ctx context.Context, query Query) Result
	RefreshAccessToken(ctx context.Context, query Query) Result
}

func NewProfileUseCase(repo domain.ProfileRepositoryInterface) ProfileUseCaseInterface {
	return &ProfileUseCase{
		repo: repo,
	}
}

// ActiveProfile implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) ActiveProfile(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// GetProfile implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) GetProfile(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// Logout implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) Logout(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// RefreshAccessToken implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) RefreshAccessToken(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// UpdatePassword implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) UpdatePassword(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// UpdateProfile implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) UpdateProfile(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// UpdateProfilePicture implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) UpdateProfilePicture(ctx context.Context, query Query) Result {
	panic("unimplemented")
}
