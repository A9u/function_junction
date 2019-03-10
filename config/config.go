package config

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type config struct {
	appName    string
	appPort    int
	db         databaseConfig
	smtpApiKey string
	url        string
	allEmail   string
}

var appConfig config

func Load() {
	viper.SetDefault("APP_NAME", "boilerplate")
	viper.SetDefault("APP_PORT", "8000")

	viper.SetConfigName("application")

	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	appConfig = config{
		appName:    readEnvString("APP_NAME"),
		appPort:    readEnvInt("APP_PORT"),
		smtpApiKey: readEnvString("SMTP_API_KEY"),
		db:         newDatabaseConfig(),
		url:        readEnvString("URL"),
		allEmail:   readEnvString("ALL_EMAIL"),
	}
}

func AppName() string {
	return appConfig.appName
}

func AppPort() int {
	return appConfig.appPort
}

func readEnvInt(key string) int {
	checkIfSet(key)
	v, err := strconv.Atoi(viper.GetString(key))
	if err != nil {
		panic(fmt.Sprintf("key %s is not a valid integer", key))
	}
	return v
}

func readEnvString(key string) string {
	checkIfSet(key)
	return viper.GetString(key)
}

func checkIfSet(key string) {
	if !viper.IsSet(key) {
		err := errors.New(fmt.Sprintf("Key %s is not set", key))
		panic(err)
	}
}

func SmtpApiKey() string {
	return appConfig.smtpApiKey
}

func URL() string {
	return appConfig.url
}

func AllEmail() string {
	return appConfig.allEmail
}
