package config

import "github.com/spf13/viper"

type Config struct {
	DBSource            string `mapstructure:"DB_SOURCE"`
	ServerAddr          string `mapstructure:"SERVER_ADDR"`
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	ExchangeRatesAPIKey string `mapstructure:"EXCHANGE_RATES_API_KEY"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("../")
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
