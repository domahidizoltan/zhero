// Package controller is for defining routes and handler functions.
package controller

import (
	"github.com/domahidizoltan/zhero/config"
	schemaorgCtrl "github.com/domahidizoltan/zhero/controller/schemaorg"
	metaRepo "github.com/domahidizoltan/zhero/repository/schemametadata"
	metaSvc "github.com/domahidizoltan/zhero/service/schemametadata"
	schemaorgSvc "github.com/domahidizoltan/zhero/service/schemaorg"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SetRoutes(router *gin.Engine) {
	router.Static("/static", "./templates")

	schemaorgSvc, err := schemaorgSvc.New(config.RdfConfig{File: "rdf_schema.jsonld"})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Schema.org service")
	}

	mRepo, err := metaRepo.New("zhero.db")
	// defer mRepo.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create schema metadata repository")
	}
	metaSvc := metaSvc.New(mRepo)
	schemaorgCtrl := schemaorgCtrl.New(*schemaorgSvc, metaSvc)

	router.GET("/", schemaorgCtrl.Search)
	router.GET("/schema/:class/edit", schemaorgCtrl.Edit)
	router.POST("/schema/:class/save", schemaorgCtrl.Save)
	router.GET("/class-hierarchy", schemaorgCtrl.GetClassHierarchy)
}
