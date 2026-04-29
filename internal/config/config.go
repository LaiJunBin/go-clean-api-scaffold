package config

import (
	"log"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string // development | release | production
	Server   *ServerConfig
	Database *DatabaseConfig
	Log      *LogConfig
	Test     *TestConfig
}

type TestConfig struct {
	Server *ServerConfig
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

type ServerConfig struct {
	Port           int
	AllowedOrigins []string
}

type DatabaseConfig struct {
	Sqlite *SqliteConfig
}

type SqliteConfig struct {
	Path string
}

type LogConfig struct {
	Level string // debug | info | warn | error  (default: info)
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Join(filepath.Dir(d), "..")
}

var configRootDir = rootDir()
var configName = "config"

func NewConfig() *Config {
	viper.AddConfigPath(configRootDir)
	viper.SetConfigName(configName)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("read config failed, err:%v\n", err)
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Panicf("unmarshal config failed, err:%v\n", err)
	}

	return &config
}

func NewTestConfig() *Config {
	v := viper.New()
	v.AddConfigPath(configRootDir)
	v.SetConfigName("config_test")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		log.Panicf("read test config failed, err:%v\n", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Panicf("unmarshal test config failed, err:%v\n", err)
	}
	return &config
}
