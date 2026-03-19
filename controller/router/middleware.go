package router

import (
	"net/http"
	"slices"
	"strings"

	"github.com/domahidizoltan/zhero/controller/dynamicpage"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var skipPrefixes = []string{"static", "asset", "preview", "favicon.ico"}

func CustomRouteMiddleware(svc Services, dynamicPageCtrl dynamicpage.Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		prefix, _, _ := strings.Cut(requestPath[1:], "/")
		if slices.Contains(skipPrefixes, prefix) {
			setParams(c, map[string]string{"skipLoadPage": "true"})
			c.Next()
			return
		}

		customRoute, err := svc.Route.GetByRoute(c.Request.Context(), requestPath)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", requestPath).
				Msg("failed to query custom route")
			c.Next()
			return
		}

		if customRoute == nil {
			route, err := svc.Route.GetLatestVersion(c.Request.Context(), requestPath[1:])
			if err != nil {
				log.Error().
					Err(err).
					Str("page", requestPath).
					Msg("failed to get route by request path")
			}
			if route != nil {
				log.Info().
					Str("page", route.Page).
					Str("requested", requestPath).
					Str("redirect", route.Route).
					Msg("redirecting page to assigned route")
				c.Redirect(http.StatusMovedPermanently, route.Route)
				c.Abort()
				return
			}

			c.Next()
			return
		}

		latestRoute, err := svc.Route.GetLatestVersion(c.Request.Context(), customRoute.Page)
		if err != nil {
			log.Error().
				Err(err).
				Str("page", customRoute.Page).
				Msg("failed to get latest route version")
			c.Next()
			return
		}

		if latestRoute != nil && latestRoute.Version > customRoute.Version {
			log.Info().
				Str("page", customRoute.Page).
				Str("requested", customRoute.Route).
				Str("redirect", latestRoute.Route).
				Msg("redirecting outdated page")
			c.Redirect(http.StatusMovedPermanently, latestRoute.Route)
			c.Abort()
			return
		}

		parts := strings.Split(latestRoute.Page, "/")
		schemaName, identifier := parts[0], parts[1]
		setParams(c, map[string]string{"class": schemaName, "identifier": identifier})
		c.Next()
	}
}

func setParams(c *gin.Context, kv map[string]string) {
	p := gin.Params{}
	for k, v := range kv {
		p = append(p, gin.Param{Key: k, Value: v})
	}
	c.Params = mergeParams(c.Params, p)
}

func mergeParams(existing gin.Params, new gin.Params) gin.Params {
	result := make(gin.Params, 0, len(existing)+len(new))
	result = append(result, existing...)

	for _, newParam := range new {
		if idx := slices.IndexFunc(result, func(p gin.Param) bool {
			return p.Key == newParam.Key
		}); idx > -1 {
			result[idx] = newParam
		} else {
			result = append(result, newParam)
		}
	}

	return result
}
