package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBHost              string        `mapstructure:"DB_HOST"`
	DBPort              string        `mapstructure:"DB_PORT"`
	DBPortGa            string        `mapstructure:"DB_PORT_GA"`
	DBUser              string        `mapstructure:"DB_USER"`
	DBPass              string        `mapstructure:"DB_PASS"`
	DBName              string        `mapstructure:"DB_NAME"`
	ServerAddr          string        `mapstructure:"SERVER_ADDR"`
	SecretKey           string        `mapstructure:"SECRET_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
