// Package main is for bootstrapping the app
package main

import (
	"database/sql"
	"fmt"

	"github.com/domahidizoltan/zhero/config"
	"github.com/domahidizoltan/zhero/controller/router"
	"github.com/domahidizoltan/zhero/data/db/sqlite"
	"github.com/domahidizoltan/zhero/pkg/database"
	"github.com/domahidizoltan/zhero/pkg/logging"
	"github.com/domahidizoltan/zhero/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/domain/schemaorg"
	page_repo "github.com/domahidizoltan/zhero/repository/page"
	meta_repo "github.com/domahidizoltan/zhero/repository/schema"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg, err := config.LoadConfig()
	logging.ConfigureLogging(cfg)

	r := gin.New()
	r.Use(
		gin.Recovery(),
		logging.ZerologMiddleware(log.Logger),
		session.SessionMiddleware(),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	if err := database.InitSqliteDB(cfg.DB.SQLite.File); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	defer func() {
		if err := database.GetDB().Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()

	if err := database.Migrate(database.GetDB(), sqlite.Scripts); err != nil {
		log.Fatal().Err(err).Msg("failed to run database migrations")
	}

	services := getRouterServices(database.GetDB(), *cfg)
	router.SetRoutes(r, services)

	serverAddr := fmt.Sprintf(":%d", cfg.Admin.Server.Port)
	log.Info().Int("port", cfg.Admin.Server.Port).Msg("server started on port")
	if err := r.Run(serverAddr); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func getRouterServices(db *sql.DB, cfg config.Config) router.Services {
	schemaorgSvc, err := schemaorg.NewService(cfg.Admin.RDF)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Schema.org service")
	}

	metaRepo := meta_repo.NewRepo(db)
	metaSvc := schema.NewService(metaRepo, schemaorgSvc)
	pageRepo := page_repo.NewRepo(db)
	pageSvc := page.NewService(pageRepo)

	return router.Services{
		Schema: metaSvc,
		Page:   pageSvc,
	}
}
