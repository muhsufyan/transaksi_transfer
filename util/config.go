package util

import (
	"time"

	"github.com/spf13/viper"
)

/*
semua config untuk development disimpan di root/app.env
*/

// config stores all configuration of the app (env). the value are read from file environment variabel (app.env)
type Config struct {
	// get the value from env variable (app.env) gunakan unmashal(). viper use mapstructure untuk unmarshaling values
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// read configurasi from file / env variable (our case is app.env)
func LoadConfig(path string) (config Config, err error) {
	// lokasi dari file config
	viper.AddConfigPath(path)
	// mencari config file dg nama sesuai dg isi param (our case is app.env so app is the name)
	viper.SetConfigName("app")
	// tipe config file (our case is .env at app.env). selain env bisa juga JSON,XML, atau format lainnya
	viper.SetConfigType("env")

	// viper read value from env varia
	// scra otomatis meng-override nilai yg dimiliki saat read config file dg nilai dari correspoding env var if exist
	viper.AutomaticEnv()
	// start reading config value
	err = viper.ReadInConfig()
	if err != nil {
		// kembalikan err
		return
	}
	// if no error maka unmarshal nilai kedlm objek target config
	viper.Unmarshal(&config)
	// kembalikan config obj
	return
}
