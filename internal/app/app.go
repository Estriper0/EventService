package app

import (
	"database/sql"
	"log/slog"

	rd "github.com/Estriper0/EventService/internal/cache/redis"
	"github.com/Estriper0/EventService/internal/config"
	db "github.com/Estriper0/EventService/internal/repositories/database"
	event_repo "github.com/Estriper0/EventService/internal/repositories/database/event"
	"github.com/Estriper0/EventService/internal/server"
	event_service "github.com/Estriper0/EventService/internal/service/event"
	"github.com/redis/go-redis/v9"
)

type App struct {
	logger     *slog.Logger
	config     *config.Config
	grpcServer *server.GRPCServer
	db         *sql.DB
}

func New(
	logger *slog.Logger,
	config *config.Config,
) *App {
	db := db.GetDB(&config.DB)

	eventRepo := event_repo.New(db)
	redisClient := redis.NewClient(&redis.Options{Addr: config.Redis.Addr, Password: config.Redis.Password})
	cache := rd.New(redisClient)
	eventService := event_service.New(eventRepo, cache, logger, config)
	grpcServer := server.New(logger, config, eventService)

	return &App{
		logger:     logger,
		config:     config,
		grpcServer: grpcServer,
		db:         db,
	}
}

func (a *App) Run() {
	a.logger.Info("Start application")

	a.grpcServer.Run()
}

func (a *App) Stop() {
	a.grpcServer.Stop()
	a.db.Close()

	a.logger.Info("Stop application")
}
