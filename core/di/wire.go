// wire.go
package di

import (
	"go-gaurd/api/controller/private"
	"go-gaurd/api/controller/public"
	"go-gaurd/core/config"
	"go-gaurd/database"
	"go-gaurd/feature/auth/domain"
	"go-gaurd/feature/auth/usecase"
	domain2 "go-gaurd/feature/profile/domain"
	usecase2 "go-gaurd/feature/profile/usecase"

	"github.com/google/wire"
)

// TODO FIX REDIS INITLAIZER
// Provider set for Redis cache (singleton)
var RedisCacheSet = wire.NewSet(
	database.NewRedisCache,
	wire.Bind(new(database.RedisCache), new(*database.RedisCache)),
)

// Provider set for Database
var DatabaseSet = wire.NewSet(
	database.NewDatabase,
)

func InitializeAuthApplication(redisCache *database.RedisCache) (*public.AuthController, error) {
	wire.Build(
		config.NewConfig,
		DatabaseSet,
		domain.NewAuthRepository,
		usecase.NewAuthUseCase,
		public.NewAuthController,
	)
	return nil, nil
}

func InitializeProfileApplication(redisCache *database.RedisCache) (*private.ProfileController, error) {
	wire.Build(
		config.NewConfig,
		DatabaseSet,
		domain2.NewProfileRepository,
		usecase2.NewProfileUseCase,
		private.NewProfileController,
	)
	return nil, nil
}

// InitializeRedis creates a single Redis instance
func InitializeRedis() (*database.RedisCache, error) {
	wire.Build(
		config.NewConfig,
		database.NewRedisCache,
	)
	return nil, nil
}
