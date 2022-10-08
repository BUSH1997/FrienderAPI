package config

import (
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client/timepad"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/syncer"
	postgreslib "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	httplib "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/http"
	"github.com/spf13/viper"
)

type TransportConfig struct {
	TimePad timepad.TimePadTransportConfig `mapstructure:"timepad"`
	HTTP    httplib.Config                 `mapstructure:"http"`
}

type Config struct {
	Syncer    syncer.SyncerConfig  `mapstructure:"syncer"`
	Transport TransportConfig      `mapstructure:"transport"`
	Postgres  postgreslib.Postgres `mapstructure:"postgres"`
}

func LoadConfig(config *Config, path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("/home/ubuntu/testfriender/FrienderAPI/config/public_sync")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	return nil
}
