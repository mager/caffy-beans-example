package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/mager/caffy-beans-example/config"
	"github.com/mager/caffy-beans-example/database"
	"github.com/mager/caffy-beans-example/handler"
	"github.com/mager/caffy-beans-example/logger"
	"github.com/mager/caffy-beans-example/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Options(
			config.Module,
			database.Module,
			logger.Module,
			router.Module,
		),
		fx.Invoke(Register),
	).Run()
}

func Register(
	lc fx.Lifecycle,
	cfg *config.Config,
	database *firestore.Client,
	logger *zap.SugaredLogger,
	router *mux.Router,
) {
	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("Listening on ", cfg.Application.Address)
				go http.ListenAndServe(cfg.Application.Address, router)
				return nil
			},
			OnStop: func(context.Context) error {
				defer logger.Sync()
				defer database.Close()
				return nil
			},
		},
	)

	handler.New(logger, router, database)
}
