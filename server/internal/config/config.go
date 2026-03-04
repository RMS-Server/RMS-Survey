package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Storage  StorageConfig  `mapstructure:"storage"`
	AI       AIConfig       `mapstructure:"ai"`
	Captcha  CaptchaConfig  `mapstructure:"captcha"`
	Upload   UploadConfig   `mapstructure:"upload"`
	OAuth    OAuthConfig    `mapstructure:"oauth"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	CookieName string `mapstructure:"cookie_name"`
}

type StorageConfig struct {
	LocalPath string `mapstructure:"local_path"`
}

type AIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

type CaptchaConfig struct {
	KeyLong   int `mapstructure:"key_long"`
	ImgWidth  int `mapstructure:"img_width"`
	ImgHeight int `mapstructure:"img_height"`
}

type UploadConfig struct {
	MaxSize      int64    `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
}

type OAuthConfig struct {
	Enabled            bool   `mapstructure:"enabled"`
	ClientID           string `mapstructure:"client_id"`
	AuthURL            string `mapstructure:"auth_url"`
	TokenURL           string `mapstructure:"token_url"`
	UserInfoURL        string `mapstructure:"userinfo_url"`
	RedirectURL        string `mapstructure:"redirect_url"`
	Scopes             string `mapstructure:"scopes"`
	MinPermissionLevel int    `mapstructure:"min_permission_level"`
}

var C Config

// Load reads config from configs/config.yaml and environment variables.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&C); err != nil {
		return nil, err
	}

	return &C, nil
}
