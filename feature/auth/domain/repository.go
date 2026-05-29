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

func NewAuthRepository(db *database.Database) *AuthRepository {
	return &AuthRepository{DB: db.DB}
}

type AuthRepositoryInterface interface {
	CreateAccount(ctx context.Context, query Query) Result
}

func (ua *AuthRepository) CreateAccount(ctx context.Context, query Query) Result {
	boil.SetDB(ua.DB)
	boil.DebugMode = true

	columns := boil.Columns{
		Cols: []string{"USERNAME", "EMAIL", "PHONE", "PASSWORD", "ROLE", "SEX", "PICTURE"},
	}

	var role null.Int
	if query.User.Role == utils.RoleAdmin {
		role = null.IntFrom(utils.GetValidRoles(query.User.Role))
	}

	user := models.User{
		USERNAME: null.StringFrom(query.User.User_name),
		EMAIL:    null.StringFrom(query.User.Email),
		PHONE:    null.StringFrom(query.User.Phone),
		PASSWORD: null.StringFrom(query.User.Password),
		ROLE:     role,
		SEX:      null.StringFrom(query.User.Sex),
		PICTURE:  null.StringFrom(query.User.Picture),
	}
	err := user.Insert(ctx, ua.DB, columns)

	if err != nil {
		log.Fatalf("error inserting user: %v", err)
		return Result{
			User:    query.User,
			Success: false,
		}
	}

	return Result{
		User:    query.User,
		Success: true,
	}

}
