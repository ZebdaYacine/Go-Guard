//go:build wireinject
// +build wireinject

package gen

import (
	"go-gaurd/api/controller/private"
	"go-gaurd/api/controller/public"
	"go-gaurd/core/config"
	"go-gaurd/database"
	"go-gaurd/feature/auth/domain"
	"go-gaurd/feature/auth/usecase"
	profiledomain "go-gaurd/feature/profile/domain"
	profileusecase "go-gaurd/feature/profile/usecase"

	"github.com/google/wire"
)

var (
	ConfigSet   = wire.NewSet(config.NewConfig)
	DatabaseSet = wire.NewSet(database.NewDatabase)
	RedisSet    = wire.NewSet(database.NewRedisCache)
	MinioSet    = wire.NewSet(database.NewMinioClient)
	CoreSet     = wire.NewSet(ConfigSet, DatabaseSet, RedisSet, MinioSet)
)

// InitializeAll creates all dependencies once and returns them
func InitializeAll() (*AppDependencies, error) {
	wire.Build(
		CoreSet,
		domain.NewAuthRepository,
		usecase.NewAuthUseCase,
		public.NewAuthController,
		profiledomain.NewProfileRepository,
		profileusecase.NewProfileUseCase,
		private.NewProfileController,
		wire.Struct(new(AppDependencies), "*"),
	)
	return nil, nil
}

type AppDependencies struct {
	Config            *config.Config
	Redis             *database.RedisCache
	Database          *database.Database
	Minio             *database.MinioClient
	AuthController    public.AuthControllerInterface
	ProfileController private.ProfileControllerInterface
}