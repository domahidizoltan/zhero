// Package router is for defining routes and handler functions.
package router

import (
	"net/http"

	"github.com/domahidizoltan/zhero/controller/page"
	schemaorg_ctrl "github.com/domahidizoltan/zhero/controller/schema"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/gin-gonic/gin"
)

type Services struct {
	Schema schema.Service
}

func SetRoutes(router *gin.Engine, svc Services) {
	router.Static("/static", "./template")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/page/list")
	})

	schemaorgCtrl := schemaorg_ctrl.NewController(svc.Schema)
	router.GET("/schema/search", schemaorgCtrl.Search)
	router.GET("/schema/:class/edit", schemaorgCtrl.Edit)
	router.POST("/schema/:class/save", schemaorgCtrl.Save)
	router.GET("/schema/class-hierarchy", schemaorgCtrl.GetClassHierarchy)

	pageCtrl := page.NewController(svc.Schema)
	router.GET("/page/list", pageCtrl.Main)
	router.GET("/page/list/:class", pageCtrl.List)
}
