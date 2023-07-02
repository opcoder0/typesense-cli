package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/pelletier/go-toml"
)

// Config consists of client connection parameters to Typesense server
type Config struct {
	Host   string
	Port   int64
	APIKey string
}

// ErrNotFound indicates the specified config instance was not found
var ErrNotFound = errors.New("config entry not found")

// ErrInvalidKey indicates the config key has invalid type or value
var ErrInvalidKey = errors.New("invalid config key")

// Load loads typesense client config for the specified instance
func Load(instance string) (Config, error) {
	usr, _ := user.Current()
	cfgPath := fmt.Sprintf("%s/.typesense-cli/config.toml", usr.HomeDir)
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatal("typesense-cli config:", err)
	}
	config, err := toml.LoadBytes(data)
	if err != nil {
		log.Fatal("typesense-cli config load:", err)
	}
	host := config.Get(fmt.Sprintf("%s.host", instance))
	port := config.Get(fmt.Sprintf("%s.port", instance))
	apiKey := config.Get(fmt.Sprintf("%s.apikey", instance))
	if host == nil || port == nil || apiKey == nil {
		return Config{}, ErrNotFound
	}
	portInt, ok := port.(int64)
	if !ok {
		return Config{}, fmt.Errorf("%q: %w", "port", ErrInvalidKey)
	}
	if host.(string) == "" {
		return Config{}, fmt.Errorf("%q: %w", "host", ErrInvalidKey)
	}
	if apiKey.(string) == "" {
		return Config{}, fmt.Errorf("%q: %w", "apikey", ErrInvalidKey)
	}
	return Config{
		Host:   host.(string),
		Port:   portInt,
		APIKey: apiKey.(string),
	}, nil
}
