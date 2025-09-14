// Package controller is for defining routes and handler functions.
package controller

import (
	"github.com/domahidizoltan/zhero/config"
	schemaorgCtrl "github.com/domahidizoltan/zhero/controller/schemaorg"
	schemaorgSvc "github.com/domahidizoltan/zhero/service/schemaorg"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SetRoutes(router *gin.Engine) {
	router.Static("/static", "./templates")

	svc, err := schemaorgSvc.New(config.RdfConfig{File: "rdf_schema.jsonld"})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Schema.org service")
	}
	schemaorgCtrl := schemaorgCtrl.New(*svc)

	router.GET("/", schemaorgCtrl.SearchSchema)
	router.GET("/class-hierarchy", schemaorgCtrl.GetClassHierarchy)
}
