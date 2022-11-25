package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/user"
	"time"
)

type Cookie struct {
	Auth struct {
		Exp time.Duration `mapstructure:"exp"`
	} `mapstructure:"auth"`
	Refresh struct {
		Exp time.Duration `mapstructure:"exp"`
	} `mapstructure:"refresh"`
}

type Config struct {
	AuthSignSecret string `mapstructure:"auth_secret"`
	Cookie         Cookie `mapstructure:"cookie"`
}

type UserUseCase struct {
	UserConfig Config
	UserRepo   user.Repository
	Logger     hardlogger.Logger
}

func New(
	userRepository user.Repository,
	userConfig Config,
	logger hardlogger.Logger,
) user.UseCase {
	return &UserUseCase{
		UserConfig: userConfig,
		UserRepo:   userRepository,
		Logger:     logger,
	}
}
