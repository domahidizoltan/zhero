// Package main is for bootstrapping the app
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/domahidizoltan/zhero/config"
	"github.com/domahidizoltan/zhero/controller/preview"
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

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
	adminSrv := createAndStartServer("Admin", cfg.Admin.Server.Port, func(e *gin.Engine) {
		router.SetAdminRoutes(e, services)
	})
	publicSrv := createAndStartServer("Public", cfg.Public.Server.Port, func(e *gin.Engine) {
		router.SetPublicRoutes(e, services)
	})

	<-quit
	log.Info().Msg("Shutting down servers...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := publicSrv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Public server forced to shutdown")
	}

	if err := adminSrv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Admin server forced to shutdown")
	}

	log.Info().Msg("Servers exited")
}

func createAndStartServer(serverName string, port int, setRoutes func(*gin.Engine)) *http.Server {
	ginRouter := gin.New()
	ginRouter.Use(
		gin.Recovery(),
		logging.ZerologMiddleware(log.Logger),
		session.SessionMiddleware(),
	)

	setRoutes(ginRouter)

	srvAddr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:    srvAddr,
		Handler: ginRouter,
	}

	go func() {
		log.Info().Int("port", port).Msg(serverName + " server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg(serverName + " server listen: %s")
		}
	}()

	return srv
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

	previewCtrl := preview.NewController(metaSvc)

	return router.Services{
		Schema:  metaSvc,
		Page:    pageSvc,
		Preview: previewCtrl,
	}
}
