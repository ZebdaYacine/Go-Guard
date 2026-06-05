package domain

import (
	"context"
	"database/sql"
	"fmt"
	"go-gaurd/core/utils"
	"go-gaurd/database"
	"go-gaurd/models"
	"log"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/gofrs/uuid"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *database.Database) AuthRepositoryInterface {
	return &AuthRepository{DB: db.DB}
}

type AuthRepositoryInterface interface {
	CreateAccount(ctx context.Context, u User_Entity) Result
	Login(ctx context.Context, l Login_Entity) Result
}

func (ua *AuthRepository) CreateAccount(ctx context.Context, u User_Entity) Result {
	boil.SetDB(ua.DB)
	boil.DebugMode = true

	log.Printf("User data: Username='%s', Email='%s', Phone='%s', Password='%s', Role=%s, Sex='%s', Picture='%s'",
		u.User_name, u.Email, u.Phone, u.Password, u.Role, u.Sex, u.Picture)

	newID := uuid.Must(uuid.NewV7())
	user := models.User{
		ID:       newID.String(),
		Username: null.StringFrom(u.User_name),
		Email:    null.StringFrom(u.Email),
		Phone:    null.StringFrom(u.Phone),
		Password: null.StringFrom(utils.HashPasswordSHA256(u.Password)),
		Role:     null.IntFrom(utils.GetValidRoles(u.Role)),
		Sex:      null.StringFrom(utils.GetCodeGender(u.Sex)),
		Picture:  null.StringFrom(u.Picture),
	}

	err := user.Insert(ctx, ua.DB, boil.Infer())

	if err != nil {
		fmt.Printf("error inserting user: %v \n", err)
		alert := utils.HandleMysqlError(err)
		return Result{
			Error:   alert,
			User:    u,
			Success: false,
		}
	}

	user_result := User_Entity{
		User_name: user.Username.String,
		Email:     user.Email.String,
		Phone:     user.Phone.String,
		Password:  "",
		Role:      "", // You might want to map this back from user.Role.Int
		Sex:       user.Sex.String,
		Picture:   user.Picture.String,
	}

	return Result{
		Id:      user.ID,
		User:    user_result,
		Success: true,
	}
}

// Login implements [AuthRepositoryInterface].
// Login authenticates a user and returns the user entity if credentials are valid
func (ua *AuthRepository) Login(ctx context.Context, l Login_Entity) Result {
	boil.SetDB(ua.DB)
	boil.DebugMode = true

	user, err := models.Users(
		models.UserWhere.Email.EQ(null.StringFrom(l.Email)),
	).One(ctx, ua.DB)

	if err != nil {

		if err == sql.ErrNoRows {
			return Result{
				User:    nil,
				Success: false,
				Error:   "User not found",
			}
		}
		fmt.Printf("error querying user: %v \n", err)
		return Result{
			User:    nil,
			Success: false,
			Error:   "database error",
		}

	}

	// Verify password exists
	if !user.Password.Valid {
		return Result{
			User:    nil,
			Success: false,
			Error:   "password incorrect",
		}
	}

	// IMPORTANT: Verify the password matches
	if !utils.CheckPasswordHashSHA256(l.Password, user.Password.String) {
		return Result{
			User:    nil,
			Success: false,
			Error:   "password incorrect",
		}
	}

	// Convert models.User to User_Entity
	userEntity := User_Entity{
		User_name: user.Username.String,
		Email:     user.Email.String,
		Phone:     user.Phone.String,
		Password:  "", // Don't return the password hash for security
		Role:      "", // Convert int64 to int
		Sex:       user.Sex.String,
		Picture:   user.Picture.String,
	}

	return Result{
		User:    userEntity,
		Success: true,
		Id:      user.ID,
		Error:   "",
	}
}
