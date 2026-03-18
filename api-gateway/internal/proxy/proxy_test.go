package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type responseWriterWithCloseNotify struct {
	*httptest.ResponseRecorder
}

func (rw *responseWriterWithCloseNotify) CloseNotify() <-chan bool {
	return make(chan bool)
}

func TestNewServiceProxy(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("proxied response"))
	}))
	defer backend.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/proxy/*path", NewServiceProxy(backend.URL))

	req := httptest.NewRequest(http.MethodGet, "/proxy/some/path", nil)
	w := &responseWriterWithCloseNotify{httptest.NewRecorder()}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "proxied response", w.Body.String())
}
