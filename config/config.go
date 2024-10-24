package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	RedisHost  string `mapstructure:"REDIS_HOST"`
	RedisPort  string `mapstructure:"REDIS_PORT"`

	AuthPath              string `mapstructure:"AUTH_PATH"`
	AuthRegisterPath      string `mapstructure:"AUTH_REGISTER_PATH"`
	AuthLoginPath         string `mapstructure:"AUTH_LOGIN_PATH"`
	ReferralPath          string `mapstructure:"REFERRAL_PATH"`
	ReferralCreatePath    string `mapstructure:"REFERRAL_CREATE_PATH"`
	ReferralDeletePath    string `mapstructure:"REFERRAL_DELETE_PATH"`
	ReferralGetPath       string `mapstructure:"REFERRAL_GET_PATH"`
	ReferralRegisterPath  string `mapstructure:"REFERRAL_REGISTER_PATH"`
	ReferralReferralsPath string `mapstructure:"REFERRAL_REFERRALS_PATH"`
}

var C Config

func InitConfig() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	viper.AutomaticEnv()

	C = Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		RedisHost:  viper.GetString("REDIS_HOST"),
		RedisPort:  viper.GetString("REDIS_PORT"),

		AuthPath:              viper.GetString("AUTH_PATH"),
		AuthRegisterPath:      viper.GetString("AUTH_REGISTER_PATH"),
		AuthLoginPath:         viper.GetString("AUTH_LOGIN_PATH"),
		ReferralPath:          viper.GetString("REFERRAL_PATH"),
		ReferralCreatePath:    viper.GetString("REFERRAL_CREATE_PATH"),
		ReferralDeletePath:    viper.GetString("REFERRAL_DELETE_PATH"),
		ReferralGetPath:       viper.GetString("REFERRAL_GET_PATH"),
		ReferralRegisterPath:  viper.GetString("REFERRAL_REGISTER_PATH"),
		ReferralReferralsPath: viper.GetString("REFERRAL_REFERRALS_PATH"),
	}
}
