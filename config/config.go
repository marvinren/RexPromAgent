package config

import (
	"fmt"
	"log"
)
import "github.com/spf13/viper"

type Config struct {
	server     Server
	database   Database
	prometheus Prometheus
}

type Database struct {
	driver          string
	dsn             string
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime int
}

type Prometheus struct {
	prometheusConfigPath string
	alertRulesConfigPath string
}

type Server struct {
	address string
}

// global config object
var config Config

func Initialize() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../.")
	viper.AddConfigPath("../../.")

	viper.SetDefault("server.address", ":9093")
	viper.SetDefault("database.maxIdleConns", 3)
	viper.SetDefault("database.maxOpenConns", 64)
	viper.SetDefault("database.connMaxLifetime", 10)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("parser config file error, %s", err)
		return nil
	}
	return &config
}
