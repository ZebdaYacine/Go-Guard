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
	RestPassword(ctx context.Context, query Query) Result
	CheckUserExists(ctx context.Context, email string) Result
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

	user_entity := query.User.(User_Entity)
	user := domain.User_Entity{
		User_name: user_entity.User_name,
		Email:     user_entity.Email,
		Phone:     user_entity.Phone,
		Password:  user_entity.Password,
		Role:      user_entity.Role,
		Sex:       user_entity.Sex,
		Picture:   user_entity.Picture,
	}

	result := au.UserRepository.CreateAccount(ctx, user)
	if !result.Success {
		return Result{
			User:    query,
			Message: result.Error,
			Success: false,
		}
	}
	return Result{
		Id:      result.Id,
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

	if !result.Success {
		return Result{
			User:    login_entity,
			Message: result.Error,
			Success: false,
		}
	}

	return Result{
		User:    result.User,
		Id:      result.Id,
		Message: "Login successful",
		Success: true,
	}
}

// ForgetPassword implements [AuthUseCaseInterface].
func (a *AuthUseCase) CheckUserExists(ctx context.Context, email string) Result {
	result := a.UserRepository.GetUserByEmail(ctx, email)
	if user := result.User; user == nil {
		return Result{
			User:    nil,
			Success: false,
			Message: result.Error,
		}
	}
	return Result{
		User:    nil,
		Success: true,
		Message: "User found",
	}
}

// SendOTP implements [AuthUseCaseInterface].
func (a *AuthUseCase) SendOTP(ctx context.Context, email string, purpose string) Result {
	panic("unimplemented")
}

// RestPassword implements [AuthUseCaseInterface].
func (a *AuthUseCase) RestPassword(ctx context.Context, query Query) Result {

	resetPassword_entity := query.User.(ResetPassword_Entity)
	err := validate_entity(resetPassword_entity)

	if err != nil {
		return Result{
			User:    resetPassword_entity,
			Message: "Invalid Input",
			Success: false,
		}
	}

	if resetPassword_entity.ConfirmePassword != resetPassword_entity.NewPassword {
		return Result{
			User:    resetPassword_entity,
			Message: "Passwords do not match",
			Success: false,
		}
	}

	q := domain.ResetPassword_Entity{
		Email:       resetPassword_entity.Email,
		NewPassword: resetPassword_entity.NewPassword,
	}

	result := a.UserRepository.RestPassword(ctx, q)

	if !result.Success {
		return Result{
			User:    resetPassword_entity,
			Message: "Error processing password",
			Success: false,
		}
	}

	return Result{
		User:    resetPassword_entity,
		Message: "Password reset successfully",
		Success: true,
	}
}
