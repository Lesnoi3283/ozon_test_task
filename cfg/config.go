package cfg

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Cfg struct {
	LogLevel             string
	ServerAddress        string
	ServerPort           string
	DefaultCommentsLimit int
	MaxCommentsLimit     int
	DefaultPostsLimit    int
	MaxPostsLimit        int
	DBConnectionString   string
	InMemoryStorage      bool
	RedisAddress         string
	RedisPort            string
	RedisPassword        string
	MaxCommentTextLength int
	DebugMode            bool
}

// Configure reads values from env and command line args into a Cfg structure.
func Configure() (*Cfg, error) {
	cfg := &Cfg{}

	flag.BoolVar(&cfg.InMemoryStorage, "m", false, "switch to in-memory storage (Redis)")
	flag.BoolVar(&cfg.DebugMode, "d", false, "enable debug mode")
	flag.Parse()

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	} else {
		cfg.LogLevel = "debug"
	}

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

	if redisAddr := os.Getenv("REDIS_ADDRESS"); redisAddr != "" {
		cfg.RedisAddress = redisAddr
	} else {
		cfg.RedisAddress = "localhost"
	}

	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		cfg.RedisPort = redisPort
	} else {
		cfg.RedisPort = "6379"
	}

	if redisPass := os.Getenv("REDIS_PASSWORD"); redisPass != "" {
		cfg.RedisPassword = redisPass
	} else {
		cfg.RedisPassword = ""
	}

	if val := os.Getenv("MAX_COMMENT_TEXT_LENGTH"); val != "" {
		maxLength, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_COMMENT_TEXT_LENGTH: %w", err)
		}
		cfg.MaxCommentTextLength = maxLength
	} else {
		cfg.MaxCommentTextLength = 2000
	}

	return cfg, nil
}
