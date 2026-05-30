//go:build wireinject
// +build wireinject

package di

import (
	"go-gaurd/api/controller/public"
	"go-gaurd/core/config"
	"go-gaurd/database"
	"go-gaurd/feature/auth/domain"
	"go-gaurd/feature/auth/usecase"

	"github.com/google/wire"
)

func InitializeAuthApplication() (*public.AuthController, error) {
	wire.Build(
		config.NewConfig,
		database.NewDatabase,
		domain.NewAuthRepository,
		usecase.NewAuthUseCase,
		database.NewRedisCache,
		public.NewAuthController,
	)
	return nil, nil
}
