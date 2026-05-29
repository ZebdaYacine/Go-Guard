//go:build wireinject
// +build wireinject

package di

import (
	"go-gaurd/core/config"
	"go-gaurd/database"
	"go-gaurd/feature/auth/domain"
	"go-gaurd/feature/auth/usecase"

	"github.com/google/wire"
)

func InitializeAuthApplication() (*usecase.AuthUseCase, error) {
	wire.Build(
		config.NewConfig,
		database.NewDatabase,
		domain.NewAuthRepository,
		usecase.NewAuthUseCase,
	)
	return nil, nil
}
