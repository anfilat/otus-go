package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func newConfig(configFile string) (Config, error) {
	config := Config{}

	v := viper.New()

	configure(v)

	if configFile != "" {
		v.SetConfigFile(configFile)
		err := v.ReadInConfig()
		if err != nil {
			return config, fmt.Errorf("failed to read configuration: %w", err)
		}
	}

	if err := v.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("failed to validate configuration: %w", err)
	}

	return config, nil
}

func configure(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("logger.level", "INFO")

	v.SetDefault("server.host", "127.0.0.1")
	v.SetDefault("server.httpPort", "8080")
	v.SetDefault("server.grpcPort", "8081")

	v.SetDefault("database.inmem", true)
}

type Config struct {
	Logger   LoggerConf
	Server   ServerConf
	Database DatabaseConf
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return err
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	return nil
}

type LoggerConf struct {
	Level string
	File  string
}

type ServerConf struct {
	Host     string
	HTTPPort string
	GrpcPort string
}

func (c ServerConf) Validate() error {
	if c.Host == "" {
		return errors.New("http app server host is required")
	}

	if c.HTTPPort == "" {
		return errors.New("http app server port is required")
	}

	if c.GrpcPort == "" {
		return errors.New("grpc app server port is required")
	}

	return nil
}

type DatabaseConf struct {
	Inmem   bool
	Connect string
}

func (c DatabaseConf) Validate() error {
	if !c.Inmem && c.Connect == "" {
		return errors.New("database connect is required")
	}

	return nil
}
