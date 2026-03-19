// Package router is for defining routes and handler functions.
package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/controller"
	page_ctrl "github.com/domahidizoltan/zhero/controller/adminpage"
	schemaorg_ctrl "github.com/domahidizoltan/zhero/controller/adminschema"
	dynamicpage_ctrl "github.com/domahidizoltan/zhero/controller/dynamicpage"
	"github.com/domahidizoltan/zhero/controller/pagerenderer"
	preview_ctrl "github.com/domahidizoltan/zhero/controller/preview"
	template_ctrl "github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/route"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/template"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Services struct {
	Schema              schema.Service
	Page                page.Service
	DynamicPageRenderer pagerenderer.DynamicPageRenderer
	Route               route.Service
}

var mimeTypes = map[string]string{
	"js":  "text/javascript",
	"css": "text/css",
}

func addCommonHandlers(router *gin.Engine, isAdmin bool) {
	staticRoot := "./template"
	assets := template.Assets
	if isAdmin {
		staticRoot += "/admin"
		assets = template.AdminAssets
	}

	router.Static("/static", staticRoot)

	router.GET("/asset/*path", func(ctx *gin.Context) {
		assetPath := ctx.Param("path")
		mimeType := "text/plain"

		if content, found := assets[assetPath]; found {
			if extIdx := strings.LastIndex(assetPath, "."); extIdx > -1 {
				ext := strings.ToLower(assetPath[extIdx+1:])
				if mt, found := mimeTypes[ext]; found {
					mimeType = mt
				}
			}
			ctx.Data(http.StatusOK, mimeType, content)
			return
		}
		ctx.Data(http.StatusNotFound, mimeType, nil)
	})
}

func SetPublicRoutes(router *gin.Engine, svc Services) {
	addCommonHandlers(router, false)
	registerPublicPageHelpers(svc)

	router.GET("/", func(c *gin.Context) {
		schemaNames, err := svc.Page.GetEnabledSchemaNames(context.Background())
		if err != nil {
			controller.InternalServerError(c, "failed to load page", err)
			return
		}

		if len(schemaNames) == 0 {
			template_ctrl.WithLayout(c, "empty")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/"+schemaNames[0])
	})

	dynamicPageCtrl := dynamicpage_ctrl.NewController(svc.DynamicPageRenderer, svc.Schema, svc.Page)
	previewCtrl := preview_ctrl.NewController(dynamicPageCtrl)

	router.POST("/preview/:class", previewCtrl.InFlightPage)
	router.GET("/preview/:class/:identifier", previewCtrl.LoadPage)

	router.Use(CustomRouteMiddleware(svc, dynamicPageCtrl))

	router.NoRoute(func(c *gin.Context) {
		dynamicPageCtrl.LoadPage(c, true)
	})
}

func registerPublicPageHelpers(svc Services) {
	raymond.RegisterHelper("eachMenuItem", func(options *raymond.Options) raymond.SafeString {
		b := strings.Builder{}
		names, err := svc.Page.GetEnabledSchemaNames(context.Background())
		if err != nil {
			log.Err(err).Msg("failed to get menu items")
		}
		for _, name := range names {
			frame := options.NewDataFrame()
			frame.Set("menu", name)
			b.WriteString(options.FnData(frame))
		}
		return raymond.SafeString(b.String())
	})
}

func SetAdminRoutes(router *gin.Engine, svc Services) {
	addCommonHandlers(router, true)

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

		pageCtrl := page_ctrl.NewController(svc.Schema, svc.Page, svc.Route)
		admin.GET("/page/list", pageCtrl.Main)
		admin.GET("/page/list/:class", pageCtrl.List)
		admin.GET("/page/create/:class", pageCtrl.Create)
		admin.POST("/page/edit/:class", pageCtrl.EditAction)
		admin.GET("/page/edit/:class/:identifier", pageCtrl.Edit)
		admin.POST("/page/save/:class", pageCtrl.Save)
		admin.POST("/page/get-valid-slug", pageCtrl.GetValidSlug)
	}
}
