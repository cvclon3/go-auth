// cmd/api/main.go
package main

import (
	"goauthbackend.ggvp.dev/internal/data"
	"os"
	"sync"
	_ "time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"goauthbackend.ggvp.dev/internal/jsonlog"
)

const version = "1.0.0"

// `config` type to house all our app's configurations
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	redisURL string
	debug    bool
}

// Main `application` type
type application struct {
	config      config
	logger      *jsonlog.Logger
	redisClient *redis.Client
	wg          sync.WaitGroup
	models      data.Models
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := updateConfigWithEnvVariables()
	if err != nil {
		logger.PrintFatal(err, nil, cfg.debug)
	}

	db, err := openDB(*cfg)
	if err != nil {
		logger.PrintFatal(err, nil, cfg.debug)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil, cfg.debug)

	opt, err := redis.ParseURL(cfg.redisURL)
	if err != nil {
		logger.PrintFatal(err, nil, cfg.debug)
	}
	cliend := redis.NewClient(opt)

	logger.PrintInfo("redis connection pool established", nil, cfg.debug)

	app := &application{
		config:      *cfg,
		logger:      logger,
		redisClient: cliend,
		models:      data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil, cfg.debug)
	}
}
