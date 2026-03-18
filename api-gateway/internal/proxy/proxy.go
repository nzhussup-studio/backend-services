package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func NewServiceProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid target"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		c.Request.URL.Scheme = remote.Scheme
		c.Request.URL.Host = remote.Host
		c.Request.Host = remote.Host
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
