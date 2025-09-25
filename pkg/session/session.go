// Package session is for session management
package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

const (
	flashKey = "_flash"
)

func SessionMiddleware() gin.HandlerFunc {
	store := memstore.NewStore([]byte("zhero"))
	return sessions.Sessions("zheroSession", store)
}

func SetFlash(c *gin.Context, flashMsg string) error {
	s := sessions.Default(c)
	s.AddFlash(flashMsg)
	return s.Save()
}

func GetFlash(c *gin.Context) (string, error) {
	s := sessions.Default(c)
	var flashes []string
	for _, f := range s.Flashes() {
		flashes = append(flashes, f.(string))
	}

	if len(flashes) > 0 {
		return flashes[0], s.Save()
	}
	return "", s.Save()
}
