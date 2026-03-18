package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"image-service/internal/model"
	"image-service/internal/service"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMultipartFile(fieldName, filename, contentType, content string) (*multipart.FileHeader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	form, err := req.MultipartReader()
	if err != nil {
		return nil, err
	}
	mf, err := form.ReadForm(1024)
	if err != nil {
		return nil, err
	}
	files := mf.File[fieldName]
	if len(files) == 0 {
		return nil, errors.New("no files found")
	}
	return files[0], nil
}

func TestImageController_Upload_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockImageService)
	mockProd := new(MockProducer)

	ctrl := &ImageController{
		service:  &service.Service{ImageService: mockSvc},
		producer: mockProd,
	}

	fileHeader, err := createMultipartFile("file", "test.jpg", "image/jpeg", "dummy image content")
	assert.NoError(t, err)

	albumID := "album123"
	mockSvc.On("UploadImage", albumID, mock.Anything).Return([]*model.Image{
		{ID: "img1", Type: model.JPEG},
	}, nil)

	// mockProd.On("SendMessage", mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: albumID}}
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	// Manually set multipart form with our file
	c.Request.MultipartForm = &multipart.Form{
		File: map[string][]*multipart.FileHeader{"file": {fileHeader}},
	}

	ctrl.Upload(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
	mockProd.AssertExpectations(t)

	var respBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Image(s) uploaded successfully", respBody["message"])
}

func TestImageController_Upload_NoFiles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockImageService)
	mockProd := new(MockProducer)

	ctrl := &ImageController{
		service:  &service.Service{ImageService: mockSvc},
		producer: mockProd,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "album123"}}
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Request.MultipartForm = &multipart.Form{File: map[string][]*multipart.FileHeader{}}

	ctrl.Upload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "UploadImage", mock.Anything, mock.Anything)
	mockProd.AssertNotCalled(t, "SendMessage", mock.Anything)
}

func TestImageController_Upload_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockImageService)
	mockProd := new(MockProducer)

	ctrl := &ImageController{
		service:  &service.Service{ImageService: mockSvc},
		producer: mockProd,
	}

	fileHeader, err := createMultipartFile("file", "test.jpg", "image/jpeg", "dummy content")
	assert.NoError(t, err)

	albumID := "album123"
	mockSvc.On("UploadImage", mock.Anything, mock.Anything).Return([]*model.Image(nil), errors.New("mock service error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: albumID}}
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Request.MultipartForm = &multipart.Form{
		File: map[string][]*multipart.FileHeader{"file": {fileHeader}},
	}

	ctrl.Upload(c)

	assert.NotEqual(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
	mockProd.AssertNotCalled(t, "SendMessage", mock.Anything)
}

func TestImageController_Delete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockImageService)

	ctrl := &ImageController{
		service: &service.Service{ImageService: mockSvc},
	}

	albumID := "album123"
	imageID := "img123"

	mockSvc.On("DeleteImage", albumID, imageID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: albumID},
		{Key: "imageID", Value: imageID},
	}
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

	ctrl.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestImageController_Delete_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockImageService)

	ctrl := &ImageController{
		service: &service.Service{ImageService: mockSvc},
	}

	albumID := "album123"
	imageID := "img123"

	mockSvc.On("DeleteImage", albumID, imageID).Return(errors.New("delete error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: albumID},
		{Key: "imageID", Value: imageID},
	}
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

	ctrl.Delete(c)

	assert.NotEqual(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestImageController_Serve_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockImageService)

	ctrl := &ImageController{
		service: &service.Service{ImageService: mockSvc},
	}

	albumID := "album123"
	imageID := "img123"
	imagePath := "/path/to/image.jpg"

	mockSvc.On("ServeImage", albumID, imageID).Return(imagePath, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: albumID},
		{Key: "imageID", Value: imageID},
	}
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	ctrl.Serve(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestImageController_Serve_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockImageService)

	ctrl := &ImageController{
		service: &service.Service{ImageService: mockSvc},
	}

	albumID := "album123"
	imageID := "img123"

	mockSvc.On("ServeImage", albumID, imageID).Return("", errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: albumID},
		{Key: "imageID", Value: imageID},
	}
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	ctrl.Serve(c)

	assert.NotEqual(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}
