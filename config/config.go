package config

import (
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client/timepad"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client/vk"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/syncer"
	postgreslib "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	httplib "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/http"
	"github.com/spf13/viper"
)

type TransportConfig struct {
	HTTP    httplib.Config                 `mapstructure:"http"`
	VK      vk.VKTransportConfig           `mapstructure:"vk"`
	TimePad timepad.TimePadTransportConfig `mapstructure:"timepad"`
}

type VKConfig struct {
	AccessToken string `mapstructure:"access_token"`
	GroupId     string `mapstructure:"group_id"`
	AlbumId     string `mapstructure:"album_id"`
	Version     string `mapstructure:"version"`
}

type Config struct {
	Syncer    syncer.SyncerConfig  `mapstructure:"syncer"`
	Transport TransportConfig      `mapstructure:"transport"`
	Postgres  postgreslib.Postgres `mapstructure:"postgres"`
	Vk        VKConfig             `mapstructure:"vk"`
	BlackList []string             `mapstructure:"blacklist"`
	SkipList  []string             `mapstructure:"skiplist"`
}

func LoadConfig(config *Config, path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("public_sync")
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
