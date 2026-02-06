package config

import (
	"time"

	"github.com/spf13/viper"
)

type Cfg struct {
	Env  string   `mapstructure:"Env"`
	GRPC GRPC     `mapstructure:"grpc"`
	Db   Database `mapstructure:"database"`
	Path string   `mapstructure:"Path"`
}

type GRPC struct {
	Timeout time.Duration `mapstructure:"timeout"`
	Port    int           `mapstructure:"port"`
}

type Database struct {
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DbName          string        `mapstructure:"dbName"`
	Port            int           `mapstructure:"port"`
	Host            string        `mapstructure:"host"`
	MaxOpenConns    int           `mapstructure:"MaxOpenConns"`
	MaxIdleConns    int           `mapstructure:"MaxIdleConns"`
	ConnMaxIdleTime time.Duration `mapstructure:"ConnMaxIdleTime"`
	ConnMaxLifeTime time.Duration `mapstructure:"ConnMaxLifeTime"`
}

func Load() *Cfg {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

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
