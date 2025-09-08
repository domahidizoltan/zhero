// Package controller is for defining routes and handler functions.
package controller

import (
	"github.com/domahidizoltan/zhero/controller/schemaorg"
	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.Engine) {
	router.Static("/static", "./templates")

	router.GET("/", schemaorg.Create)
	router.GET("/search", schemaorg.Search)
}
