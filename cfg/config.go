package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress        string
	ServerPort           string
	DefaultCommentsLimit int
	MaxCommentsLimit     int
	DefaultPostsLimit    int
	MaxPostsLimit        int
	InMemoryStorage      bool
}

// Configure reads values from env and command line args into a Config structure.
func (c *Config) Configure() (*Config, error) {
	cfg := &Config{}

	flag.BoolVar(&c.InMemoryStorage, "m", false, "switch to in-memory storage (Redis)")

	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.ServerAddress = addr
	} else {
		cfg.ServerAddress = "127.0.0.1"
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.ServerPort = port
	} else {
		cfg.ServerPort = "8080"
	}

	if val := os.Getenv("DEFAULT_COMMENTS_LIMIT"); val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid DEFAULT_COMMENTS_LIMIT: %w", err)
		}
		cfg.DefaultCommentsLimit = limit
	} else {
		cfg.DefaultCommentsLimit = 10
	}

	if val := os.Getenv("MAX_COMMENTS_LIMIT"); val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_COMMENTS_LIMIT: %w", err)
		}
		cfg.MaxCommentsLimit = limit
	} else {
		cfg.MaxCommentsLimit = 50
	}

	if val := os.Getenv("DEFAULT_POSTS_LIMIT"); val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid DEFAULT_POSTS_LIMIT: %w", err)
		}
		cfg.DefaultPostsLimit = limit
	} else {
		cfg.DefaultPostsLimit = 10
	}

	if val := os.Getenv("MAX_POSTS_LIMIT"); val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_POSTS_LIMIT: %w", err)
		}
		cfg.MaxPostsLimit = limit
	} else {
		cfg.MaxPostsLimit = 100
	}

	return cfg, nil
}
