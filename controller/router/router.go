// Package router is for defining routes and handler functions.
package router

import (
	"net/http"

	page_ctrl "github.com/domahidizoltan/zhero/controller/page"
	schemaorg_ctrl "github.com/domahidizoltan/zhero/controller/schema"
	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/gin-gonic/gin"
)

type Services struct {
	Schema schema.Service
	Page   page.Service
}

func SetRoutes(router *gin.Engine, svc Services) {
	router.Static("/static", "./template")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/page/list")
	})

	schemaorgCtrl := schemaorg_ctrl.NewController(svc.Schema)
	router.GET("/schema/search", schemaorgCtrl.Search)
	router.GET("/schema/edit/:class", schemaorgCtrl.Edit)
	router.POST("/schema/save/:class", schemaorgCtrl.Save)
	router.GET("/schema/class-hierarchy", schemaorgCtrl.GetClassHierarchy)

	pageCtrl := page_ctrl.NewController(svc.Schema, svc.Page)
	router.GET("/page/list", pageCtrl.Main)
	router.GET("/page/list/:class", pageCtrl.List)
	router.GET("/page/create/:class", pageCtrl.Create)
	router.POST("/page/edit/:class", pageCtrl.EditAction)
	router.GET("/page/edit/:class/:identifier", pageCtrl.Edit)
	router.POST("/page/save/:class", pageCtrl.Save)
}
