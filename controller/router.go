// Package controller is for defining routes and handler functions.
package controller

import (
	schemaorg_ctrl "github.com/domahidizoltan/zhero/controller/schemaorg"
	"github.com/domahidizoltan/zhero/domain/schemametadata"
	"github.com/domahidizoltan/zhero/domain/schemaorg"
	"github.com/gin-gonic/gin"
)

type Services struct {
	Schemaorg      schemaorg.Service
	SchemaMetadata schemametadata.Service
}

func SetRoutes(router *gin.Engine, svc Services) {
	router.Static("/static", "./templates")

	schemaorgCtrl := schemaorg_ctrl.NewController(svc.Schemaorg, svc.SchemaMetadata)
	router.GET("/", schemaorgCtrl.Search)
	router.GET("/schema/:class/edit", schemaorgCtrl.Edit)
	router.POST("/schema/:class/save", schemaorgCtrl.Save)
	router.GET("/class-hierarchy", schemaorgCtrl.GetClassHierarchy)
}
