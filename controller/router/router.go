// Package router is for defining routes and handler functions.
package router

import (
	schemaorg_ctrl "github.com/domahidizoltan/zhero/controller/schema"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/gin-gonic/gin"
)

type Services struct {
	Schema schema.Service
}

func SetRoutes(router *gin.Engine, svc Services) {
	router.Static("/static", "./templates")

	schemaorgCtrl := schemaorg_ctrl.NewController(svc.Schema)
	router.GET("/", schemaorgCtrl.Search)
	router.GET("/schema/:class/edit", schemaorgCtrl.Edit)
	router.POST("/schema/:class/save", schemaorgCtrl.Save)
	router.GET("/class-hierarchy", schemaorgCtrl.GetClassHierarchy)
}
