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
		c.Redirect(http.StatusTemporaryRedirect, "/admin/page/list")
	})

	admin := router.Group("/admin")
	{
		schemaorgCtrl := schemaorg_ctrl.NewController(svc.Schema)
		admin.GET("/schema/search", schemaorgCtrl.Search)
		admin.GET("/schema/edit/:class", schemaorgCtrl.Edit)
		admin.POST("/schema/save/:class", schemaorgCtrl.Save)
		admin.GET("/schema/class-hierarchy", schemaorgCtrl.GetClassHierarchy)

		pageCtrl := page_ctrl.NewController(svc.Schema, svc.Page)
		admin.GET("/page/list", pageCtrl.Main)
		admin.GET("/page/list/:class", pageCtrl.List)
		admin.GET("/page/create/:class", pageCtrl.Create)
		admin.POST("/page/edit/:class", pageCtrl.EditAction)
		admin.GET("/page/edit/:class/:identifier", pageCtrl.Edit)
		admin.POST("/page/save/:class", pageCtrl.Save)
	}
}
