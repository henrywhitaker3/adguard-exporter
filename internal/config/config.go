package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig/v2"
)

type Global struct {
	Server  Server
	Configs []Config
}

type Server struct {
	Interval time.Duration `env:"INTERVAL, default=30s"`
	Debug    bool          `env:"DEBUG, default=false"`
	BindAddr string        `env:"BIND_ADDR"`
}

type Config struct {
	Url      string
	Username string
	Password string
}

type EnvConfig struct {
	Urls      []string `env:"ADGUARD_SERVERS"`
	Usernames []string `env:"ADGUARD_USERNAMES"`
	Passwords []string `env:"ADGUARD_PASSWORDS"`
}

func FromEnv() (*Global, error) {
	godotenv.Load()

	env := &EnvConfig{}
	if err := envconfig.Process(context.Background(), env); err != nil {
		return nil, err
	}
	serv := &Server{}
	if err := envconfig.Process(context.Background(), serv); err != nil {
		return nil, err
	}
	if len(serv.BindAddr) == 0 {
		serv.BindAddr = ":9618"
	}

	if val, ok := os.LookupEnv("ADGUARD_PASSWORDS_FILE"); ok {
		passwords, err := os.ReadFile(val)
		if err != nil {
			return nil, fmt.Errorf("read passwords file: %w", err)
		}
		passwordValues := strings.TrimSuffix(string(passwords), "\n")
		passwordValues = strings.TrimSuffix(passwordValues, "\r")
		env.Passwords = strings.Split(passwordValues, ",")
	}

	if err := env.Validate(); err != nil {
		return nil, err
	}

	configs := []Config{}
	for i := range env.Urls {
		configs = append(configs, buildConfig(i, env.Urls, env.Usernames, env.Passwords))
	}

	return &Global{
		Server:  *serv,
		Configs: configs,
	}, nil
}

func buildConfig(i int, urls, usernames, password []string) Config {
	conf := Config{
		Url: urls[i],
	}
	if i < len(usernames) {
		conf.Username = usernames[i]
	}
	if i < len(password) {
		conf.Password = password[i]
	}

	return conf
}

func (e *EnvConfig) Validate() error {
	if len(e.Urls) == 0 {
		return errors.New("no urls supplied")
	}
	if len(e.Usernames) > 0 && len(e.Urls) != len(e.Usernames) {
		return errors.New("number of urls does not match number of usernames")
	}
	if len(e.Passwords) > 0 && len(e.Urls) != len(e.Passwords) {
		return errors.New("number of urls does not match number of passwords")
	}

	return nil
}
