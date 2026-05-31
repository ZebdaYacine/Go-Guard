package usecase

import (
	"context"
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
	Login(ctx context.Context, query Query) Result
	UpdatePassword(ctx context.Context, query Query) Result
	CheckUserExists(ctx context.Context, query Query) Result
	GetUserByEmail(ctx context.Context, email string) (int, string)
	SendOTP(ctx context.Context, email string, purpose string) Result
}



func (au *AuthUseCase) GetUserByEmail(ctx context.Context, email string) (int, string) {
	// Implementation for getting user by email logic
	return 0, ""
}

func (au *AuthUseCase) UpdatePassword(ctx context.Context, query Query) Result {
	// Implementation for update password logic
	return Result{
		Message: "Password updated successfully",
		Success: true,
	}
}

func (au *AuthUseCase) CheckUserExists(ctx context.Context, query Query) Result {
	// Implementation for checking user existence logic
	return Result{
		Message: "User found",
		Success: true,
	}
}

func (au *AuthUseCase) Login(ctx context.Context, query Query) Result {
	// Implementation for login logic
	return Result{
		Message: "Account created successfully",
		Success: true,
	}
}

func (au *AuthUseCase) CreateAccount(ctx context.Context, query Query) Result {
	// err := validate.Struct(query)
	// if err != nil {
	// 	if validationErrors, ok := err.(validator.ValidationErrors); ok {
	// 		for _, ve := range validationErrors {
	// 			fmt.Printf("Field: %s, Error: %s, Value: %v\n",
	// 				ve.Field(), ve.Tag(), ve.Value())
	// 		}
	// 	}
	// 	return Result{
	// 		User:    User_Entity(query.User),
	// 		Message: "Invalid Input",
	// 		Success: false,
	// 	}
	// }
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
