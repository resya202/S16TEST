package config

import "github.com/spf13/viper"

var (
	PostgresURL   string
	CosmosAPIBase string
)

func Load() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	PostgresURL = viper.GetString("POSTGRES_URL")
	CosmosAPIBase = viper.GetString("COSMOS_API_BASE")
	return nil
}
