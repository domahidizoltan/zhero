// Package server is a server abstraction for the app exposed to GoMobile
package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/domahidizoltan/zhero/config"
	"github.com/domahidizoltan/zhero/controller/preview"
	"github.com/domahidizoltan/zhero/controller/router"
	"github.com/domahidizoltan/zhero/data/db/sqlite"
	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/domain/schemaorg"
	"github.com/domahidizoltan/zhero/pkg/database"
	"github.com/domahidizoltan/zhero/pkg/handlebars"
	"github.com/domahidizoltan/zhero/pkg/logging"
	"github.com/domahidizoltan/zhero/pkg/session"
	page_repo "github.com/domahidizoltan/zhero/repository/page"
	meta_repo "github.com/domahidizoltan/zhero/repository/schema"
	"github.com/domahidizoltan/zhero/template"
	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
)

type Server struct {
	adminSrv, publicSrv *http.Server
	db                  *sql.DB
	absolutePath        string
}

func New() *Server {
	return &Server{}
}

func (s *Server) SetAbsolutePath(absolutePath string) {
	s.absolutePath = absolutePath
}

func (s *Server) Start() {
	gin.SetMode(gin.ReleaseMode)

	cfg, err := config.LoadConfig(s.absolutePath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	handlebars.SetAbsolutePath(cfg.Env.AbsolutePath)
	template.InitTemplates()
	logging.ConfigureLogging(cfg)

	dbFile := s.absolutePath + cfg.DB.SQLite.File
	if err := database.InitSqliteDB(dbFile); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	s.db = database.GetDB()

	if err := database.Migrate(s.db, sqlite.Scripts); err != nil {
		log.Fatal().Err(err).Msg("failed to run database migrations")
	}
	services := getRouterServices(s.db, *cfg)
	s.adminSrv = createAndStartServer("Admin", cfg.Admin.Server.Port, func(e *gin.Engine) {
		router.SetAdminRoutes(e, services)
	})
	s.publicSrv = createAndStartServer("Public", cfg.Public.Server.Port, func(e *gin.Engine) {
		router.SetPublicRoutes(e, services)
	})
}

func (s *Server) Stop() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.publicSrv != nil {
		if err := s.publicSrv.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Public server forced to shutdown")
		}
	} else {
		log.Warn().Msg("Public server instance is nil, skipping shutdown")
	}

	if s.adminSrv != nil {
		if err := s.adminSrv.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Admin server forced to shutdown")
		}
	} else {
		log.Warn().Msg("Admin server instance is nil, skipping shutdown")
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	} else {
		log.Warn().Msg("Database instance is nil, skipping close")
	}
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
	schemaorgSvc, err := schemaorg.NewService(cfg.Env.AbsolutePath, cfg.Admin.RDF)
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
