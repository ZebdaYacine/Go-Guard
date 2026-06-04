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
	UserRepository domain.AuthRepositoryInterface // Change to pointer to struct
}

type AuthUseCaseInterface interface {
	CreateAccount(ctx context.Context, query Query) Result
	Login(ctx context.Context, query Query) Result
	UpdatePassword(ctx context.Context, query Query) Result
	ForgetPassword(ctx context.Context, query Query) Result
	CheckUserExists(ctx context.Context, query Query) Result
	GetUserByEmail(ctx context.Context, email string) (int, string)
	SendOTP(ctx context.Context, email string, purpose string) Result
}

func NewAuthUseCase(repo domain.AuthRepositoryInterface) AuthUseCaseInterface {
	return &AuthUseCase{UserRepository: repo}
}

func validate_entity(entity interface{}) error {
	err := validate.Struct(entity)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, ve := range validationErrors {
				fmt.Printf("Field: %s, Error: %s, Value: %v\n",
					ve.Field(), ve.Tag(), ve.Value())
			}
		}
	}
	return err
}

// CheckUserExists implements [AuthUseCaseInterface].
func (a *AuthUseCase) CheckUserExists(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// CreateAccount implements [AuthUseCaseInterface].
func (au *AuthUseCase) CreateAccount(ctx context.Context, query Query) Result {
	err := validate_entity(query)
	if err != nil {
		return Result{
			User:    query,
			Message: "Invalid Input",
			Success: false,
		}
	}
	user := query.User.(domain.User_Entity)

	result := au.UserRepository.CreateAccount(ctx, user)
	if !result.Success {
		return Result{
			User:    query,
			Message: "Account creation failed",
			Success: false,
		}
	}

	return Result{
		User:    User_Entity(result.User.(domain.User_Entity)),
		Message: "Account created successfully",
		Success: true,
	}

}

// Login implements [AuthUseCaseInterface].
func (a *AuthUseCase) Login(ctx context.Context, query Query) Result {
	login_entity := query.User.(Login_Entity)
	err := validate_entity(login_entity)
	if err != nil {
		return Result{
			User:    login_entity,
			Message: "Invalid Input",
			Success: false,
		}
	}

	l := domain.Login_Entity{
		Email:    login_entity.Email,
		Password: login_entity.Password,
	}

	result := a.UserRepository.Login(ctx, l)

	return Result{
		User:    result.User,
		Message: "Login successful",
		Success: true,
	}
}

// ForgetPassword implements [AuthUseCaseInterface].
func (a *AuthUseCase) ForgetPassword(ctx context.Context, query Query) Result {
	panic("unimplemented")
}

// GetUserByEmail implements [AuthUseCaseInterface].
func (a *AuthUseCase) GetUserByEmail(ctx context.Context, email string) (int, string) {
	panic("unimplemented")
}

// SendOTP implements [AuthUseCaseInterface].
func (a *AuthUseCase) SendOTP(ctx context.Context, email string, purpose string) Result {
	panic("unimplemented")
}

// UpdatePassword implements [AuthUseCaseInterface].
func (a *AuthUseCase) UpdatePassword(ctx context.Context, query Query) Result {
	panic("unimplemented")
}
