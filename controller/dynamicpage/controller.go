package dynamicpage

import (
	"net/http"
	"os"

	"github.com/domahidizoltan/zhero/controller"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	dynamicPageRdr controller.UserFacingPageRenderer
}

func NewController(pageRenderer controller.UserFacingPageRenderer) Controller {
	return Controller{
		dynamicPageRdr: pageRenderer,
	}
}

func (ctrl *Controller) Index(c *gin.Context) {
	body, err := os.ReadFile("template/temp_body.html")
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html", []byte(err.Error()))
		return
	}
	_ = body
	// content, err := ctrl.dynamicPageRdr.Render(string(body))
	// if err != nil {
	// 	c.Data(http.StatusInternalServerError, "text/html", []byte(err.Error()))
	// 	return
	// }
	// c.Data(http.StatusOK, "text/html", []byte(content))
}
