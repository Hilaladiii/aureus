package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Env struct {
	AppName          string `mapstructure:"APP_NAME"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBName           string `mapstructure:"DB_NAME"`
	DBUser           string `mapstructure:"DB_USER"`
	DBPort           string `mapstructure:"DB_PORT"`
	DBPassword       string `mapstructure:"DB_PASSWORD"`
	JwtSecret        string `mapstructure:"JWT_SECRET"`
	JwtExpire        string `mapstructure:"JWT_EXPIRE"`
	StripeSecretKey  string `mapstructure:"STRIPE_SECRET_KEY"`
	StripeWebHookKey string `mapstructure:"STRIPE_WEBHOOK_KEY"`
	StripeSuccessUrl string `mapstructure:"STRIPE_SUCCESS_URL"`
	StripeCancelUrl  string `mapstructure:"STRIPE_CANCEL_URL"`
	S3Endpoint       string `mapstructure:"S3_ENDPOINT"`
	S3SecretKey      string `mapstructure:"S3_SECRET_KEY"`
	S3AccessKey      string `mapstructure:"S3_ACCESS_KEY"`
}

var envs = []string{
	"APP_NAME", "DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD",
	"JWT_SECRET", "JWT_EXPIRE",
	"STRIPE_SECRET_KEY", "STRIPE_WEBHOOK_KEY",
	"S3_ENDPOINT", "S3_SECRET_KEY", "S3_ACCESS_KEY",
}

func LoadEnv() (Env, error) {
	var config Env

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}
