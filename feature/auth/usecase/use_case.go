package usecase

import (
	"context"
	"fmt"
	"go-gaurd/feature/auth/domain"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type AuthUseCase struct {
	UserRepository *domain.AuthRepository // Change to pointer to struct
}

func NewAuthUseCase(repo *domain.AuthRepository) *AuthUseCase { // Change parameter type
	return &AuthUseCase{UserRepository: repo}
}

type AuthUseCaseInterface interface {
	CreateAccount(ctx context.Context, query Query) Result
}

func (au *AuthUseCase) CreateAccount(ctx context.Context, query Query) Result {
	err := validate.Struct(query)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, ve := range validationErrors {
				fmt.Printf("Field: %s, Error: %s, Value: %v\n",
					ve.Field(), ve.Tag(), ve.Value())
			}
		}
		return Result{
			User:    User_Entity(query.User),
			Message: "Invalid Input",
			Success: false,
		}
	}
	result := au.UserRepository.CreateAccount(ctx, domain.Query{User: domain.User_Entity(query.User)})
	if !result.Success {
		return Result{
			User:    User_Entity(result.User),
			Message: "Account creation failed",
			Success: false,
		}
	}
	return Result{
		User:    User_Entity(result.User),
		Message: "Account created successfully",
		Success: true,
	}
}
