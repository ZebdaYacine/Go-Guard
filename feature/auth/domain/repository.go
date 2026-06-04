package domain

import (
	"context"
	"database/sql"
	"go-gaurd/core/utils"
	"go-gaurd/database"
	"go-gaurd/models"
	"log"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
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

	columns := boil.Columns{
		Cols: []string{"USERNAME", "EMAIL", "PHONE", "PASSWORD", "ROLE", "SEX", "PICTURE"},
	}

	var role null.Int
	if u.Role == utils.RoleAdmin {
		role = null.IntFrom(utils.GetValidRoles(u.Role))
	}

	user := models.User{
		USERNAME: null.StringFrom(u.User_name),
		EMAIL:    null.StringFrom(u.Email),
		PHONE:    null.StringFrom(u.Phone),
		PASSWORD: null.StringFrom(u.Password),
		ROLE:     role,
		SEX:      null.StringFrom(u.Sex),
		PICTURE:  null.StringFrom(u.Picture),
	}
	err := user.Insert(ctx, ua.DB, columns)

	if err != nil {
		log.Fatalf("error inserting user: %v", err)
		return Result{
			User:    u,
			Success: false,
		}
	}

	return Result{
		User:    user,
		Success: true,
	}

}

// Login implements [AuthRepositoryInterface].
// Login authenticates a user and returns the user entity if credentials are valid
func (ua *AuthRepository) Login(ctx context.Context, l Login_Entity) Result {
	boil.SetDB(ua.DB)
	boil.DebugMode = true

	// Try to find user by username OR email (assuming Login_Entity has Username field)
	user, err := models.Users(
		models.UserWhere.USERNAME.EQ(null.StringFrom(l.Email)),
	).One(ctx, ua.DB)

	if err != nil {
		// If not found by username, try email
		user, err = models.Users(
			models.UserWhere.EMAIL.EQ(null.StringFrom(l.Email)),
		).One(ctx, ua.DB)

		if err != nil {
			if err == sql.ErrNoRows {
				return Result{
					User:    nil,
					Success: false,
					Error:   "User not found",
				}
			}
			log.Printf("error querying user: %v", err)
			return Result{
				User:    nil,
				Success: false,
				Error:   "database error",
			}
		}
	}

	// Verify password exists
	if !user.PASSWORD.Valid {
		return Result{
			User:    nil,
			Success: false,
			Error:   "password incorrect",
		}
	}

	// IMPORTANT: Verify the password matches
	if !utils.CheckPasswordHashSHA256(l.Password, user.PASSWORD.String) {
		return Result{
			User:    nil,
			Success: false,
			Error:   "password incorrect",
		}
	}

	// Convert models.User to User_Entity
	userEntity := User_Entity{
		User_name: user.USERNAME.String,
		Email:     user.EMAIL.String,
		Phone:     user.PHONE.String,
		Password:  "", // Don't return the password hash for security
		Role:      "", // Convert int64 to int
		Sex:       user.SEX.String,
		Picture:   user.PICTURE.String,
	}

	return Result{
		User:    userEntity,
		Success: true,
		Error:   "",
	}
}
