package api

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func createProxy(c *gin.Context, container string, application string, port int) *httputil.ReverseProxy {
	address := fmt.Sprintf("http://%v:%v", container, port)
	applicationPart := fmt.Sprintf("/%v", application)

	trimmedPath := strings.Replace(c.Request.URL.Path, applicationPart, "", 1)
	c.Request.URL.Path = trimmedPath

	proxiedURL, _ := url.Parse(address)
	return httputil.NewSingleHostReverseProxy(proxiedURL)
}
