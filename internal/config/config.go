package config

import (
	"time"

	"github.com/spf13/viper"
)

type Cfg struct {
	Env  string   `mapstructure:"Env"`
	GRPC GRPC     `mapstructure:"grpc"`
	Db   Database `mapstructure:"database"`
}

type GRPC struct {
	Timeout time.Duration `mapstructure:"timeout"`
}

type Database struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbName"`
}

func Load() *Cfg {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("config/")

	var cfg Cfg
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
