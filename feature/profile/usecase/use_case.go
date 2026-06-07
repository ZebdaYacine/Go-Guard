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
	GetProfile(ctx context.Context, userId string) Result
	UpdateProfile(ctx context.Context, query Query) Result
	UpdateProfilePicture(ctx context.Context, url string, userId string) Result
	UpdatePassword(ctx context.Context, query Query) Result
	ActiveProfile(ctx context.Context, userId string) Result
	// Logout(ctx context.Context, query Query) Result
	RefreshAccessToken(ctx context.Context, query Query) Result
}

func NewProfileUseCase(repo domain.ProfileRepositoryInterface) ProfileUseCaseInterface {
	return &ProfileUseCase{
		repo: repo,
	}
}

// ActiveProfile implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) ActiveProfile(ctx context.Context, userId string) Result {
	return Result{
		User:    User_Entity{},
		Success: false,
		Message: "FUNCTION NOT IMPEMENTED",
	}
}

// GetProfile implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) GetProfile(ctx context.Context, userId string) Result {
	return Result{
		User:    User_Entity{},
		Success: false,
		Message: "FUNCTION NOT IMPEMENTED",
	}
}

// Logout implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) Logout(ctx context.Context, query Query) Result {
	return Result{
		User:    User_Entity{},
		Success: false,
		Message: "FUNCTION NOT IMPEMENTED",
	}
}

// RefreshAccessToken implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) RefreshAccessToken(ctx context.Context, query Query) Result {
	return Result{
		User:    User_Entity{},
		Success: false,
		Message: "FUNCTION NOT IMPEMENTED",
	}
}

// UpdatePassword implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) UpdatePassword(ctx context.Context, query Query) Result {
	return Result{
		User:    User_Entity{},
		Success: false,
		Message: "FUNCTION NOT IMPEMENTED",
	}
}

// UpdateProfile implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) UpdateProfile(ctx context.Context, query Query) Result {
	return Result{
		User:    User_Entity{},
		Success: false,
		Message: "FUNCTION NOT IMPEMENTED",
	}
}

// UpdateProfilePicture implements [ProfileUseCaseInterface].
func (p *ProfileUseCase) UpdateProfilePicture(ctx context.Context, url string, userId string) Result {
	return Result{
		User:    User_Entity{},
		Success: false,
		Message: "FUNCTION NOT IMPEMENTED",
	}
}
