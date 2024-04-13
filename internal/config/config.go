package config

import (
	"context"
	"errors"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Global struct {
	Server  Server
	Configs []Config
}

type Server struct {
	Interval time.Duration `env:"INTERVAL, default=30s"`
	Debug    bool          `env:"DEBUG, default=false"`
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

	if err := env.Validate(); err != nil {
		return nil, err
	}

	configs := []Config{}
	for i := range env.Urls {
		configs = append(configs, Config{
			Url:      env.Urls[i],
			Username: env.Usernames[i],
			Password: env.Passwords[i],
		})
	}

	return &Global{
		Server:  *serv,
		Configs: configs,
	}, nil
}

func (e *EnvConfig) Validate() error {
	if len(e.Urls) == 0 {
		return errors.New("no urls supplied")
	}
	if len(e.Usernames) == 0 {
		return errors.New("no usernames supplied")
	}
	if len(e.Passwords) == 0 {
		return errors.New("no passwords supplied")
	}

	if len(e.Urls) != len(e.Usernames) {
		return errors.New("number of urls does not match number of usernames")
	}
	if len(e.Urls) != len(e.Passwords) {
		return errors.New("number of urls does not match number of passwords")
	}

	return nil
}
