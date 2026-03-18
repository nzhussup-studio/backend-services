package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"image-service/internal/model"
	"image-service/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAlbumController(mockSvc *MockAlbumService) (*gin.Engine, *AlbumController) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := &service.Service{AlbumService: mockSvc}
	ctrl := &AlbumController{service: svc}

	r.GET("/v1/album/:id", ctrl.Get)
	r.GET("/v1/album", ctrl.GetPreview)
	r.POST("/v1/album", ctrl.Create)
	r.DELETE("/v1/album/:id", ctrl.Delete)
	r.PUT("/v1/album/:id", ctrl.Update)

	return r, ctrl
}

func TestAlbumController_Get_Success(t *testing.T) {
	mockSvc := new(MockAlbumService)
	router, _ := setupAlbumController(mockSvc)

	expected := &model.Album{ID: "123", Title: "Test", Type: model.Public}
	mockSvc.On("GetAlbum", mock.Anything, "123").Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/v1/album/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAlbumController_Get_Error(t *testing.T) {
	mockSvc := new(MockAlbumService)
	router, _ := setupAlbumController(mockSvc)

	mockSvc.On("GetAlbum", mock.Anything, "999").Return(nil, errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/v1/album/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAlbumController_GetPreview_Success(t *testing.T) {
	mockSvc := new(MockAlbumService)
	router, _ := setupAlbumController(mockSvc)

	expected := []*model.AlbumPreview{{ID: "1", Title: "A", Type: model.Public}}
	mockSvc.On("GetAlbumsPreview", "public").Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/v1/album?type=public", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAlbumController_GetPreview_InvalidType(t *testing.T) {
	mockSvc := new(MockAlbumService)
	router, _ := setupAlbumController(mockSvc)

	req, _ := http.NewRequest(http.MethodGet, "/v1/album?type=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAlbumController_Create_Success(t *testing.T) {
	mockSvc := new(MockAlbumService)
	router, _ := setupAlbumController(mockSvc)

	reqBody := &model.AlbumPreview{Title: "New Album", Desc: "desc", Type: model.Public}
	mockSvc.On("CreateAlbum", reqBody).Return(reqBody, nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/v1/album", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAlbumController_Delete_Success(t *testing.T) {
	mockSvc := new(MockAlbumService)
	router, _ := setupAlbumController(mockSvc)

	mockSvc.On("DeleteAlbum", "123").Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/album/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAlbumController_Update_Success(t *testing.T) {
	mockSvc := new(MockAlbumService)
	router, _ := setupAlbumController(mockSvc)

	input := &model.AlbumPreview{Title: "Updated", Desc: "Updated Desc", Type: model.SemiPublic}
	mockSvc.On("UpdateAlbum", "321", input).Return(input, nil)

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPut, "/v1/album/321", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}
