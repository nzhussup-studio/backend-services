package controller

import (
	"image-service/internal/config/messaging"
	custom_errors "image-service/internal/errors"
	customJson "image-service/internal/json"
	"image-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	service  *service.Service
	producer messaging.ProducerInterface
}

// Upload godoc
// @Summary Upload image(s) to an album
// @Description Upload one or more image files to the specified album
// @Tags Image
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Album ID"
// @Param file formData file true "Image file(s) to upload"
// @Success 201 {object} model.SuccessResponse{data=[]model.Image} "Image uploaded successfully"
// @Failure 400 {object} model.ErrorDetails "Bad Request"
// @Failure 404 {object} model.ErrorDetails "Album Not Found"
// @Failure 409 {object} model.ErrorDetails "Conflict"
// @Failure 500 {object} model.ErrorDetails "Internal Server Error"
// @Router /v1/album/{id}/upload [post]
// @Security ApiKeyAuth
func (ctrl *ImageController) Upload(c *gin.Context) {
	albumID := c.Param("id")

	form, err := c.MultipartForm()
	if err != nil {
		customJson.ConstructJsonResponseError(c, custom_errors.NewError(custom_errors.ErrBadRequest, "Failed to get multipart form"), http.StatusBadRequest)
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		customJson.ConstructJsonResponseError(c, custom_errors.NewError(custom_errors.ErrBadRequest, "No files uploaded"), http.StatusBadRequest)
		return
	}

	savedImage, err := ctrl.service.ImageService.UploadImage(albumID, files)
	if err != nil {
		custom_errors.MapErrors(c, err)
		return
	}
	customJson.ConstructJsonResponseSuccess(c, savedImage, http.StatusCreated, "Image(s) uploaded successfully")

	// uploadMeta := map[string]interface{}{
	// 	"status":   "uploaded",
	// 	"album_id": albumID,
	// 	"images":   savedImage,
	// }
	// message, err := json.Marshal(uploadMeta)
	// if err != nil {
	// 	log.Printf("Error marshaling file metadata: %v", err)
	// 	return
	// }
	// err = ctrl.producer.SendMessage(message)
	// if err != nil {
	// 	log.Printf("Error sending message to Kafka: %v", err)
	// 	return
	// }
	// log.Printf("Message sent to Kafka: %s", message)
}

// Delete godoc
// @Summary Delete an image from an album
// @Description Deletes an image by ID from the given album
// @Tags Image
// @Produce json
// @Param id path string true "Album ID"
// @Param imageID path string true "Image ID"
// @Success 200 {object} model.SuccessResponse{data=nil} "Image deleted successfully"
// @Failure 400 {object} model.ErrorDetails "Bad Request"
// @Failure 404 {object} model.ErrorDetails "Image Not Found"
// @Failure 500 {object} model.ErrorDetails "Internal Server Error"
// @Router /v1/album/{id}/{imageID} [delete]
// @Security ApiKeyAuth
func (ctrl *ImageController) Delete(c *gin.Context) {
	albumID := c.Param("id")
	imageID := c.Param("imageID")

	err := ctrl.service.ImageService.DeleteImage(albumID, imageID)
	if err != nil {
		custom_errors.MapErrors(c, err)
		return
	}
	customJson.ConstructJsonResponseSuccess(c, nil, "Image deleted successfully")
}

// Serve godoc
// @Summary Serve an image file
// @Description Returns the raw image file for preview or download
// @Tags Image
// @Produce octet-stream
// @Param id path string true "Album ID"
// @Param imageID path string true "Image ID"
// @Success 200 {file} file "Image file"
// @Failure 400 {object} model.ErrorDetails "Bad Request"
// @Router /v1/album/{id}/{imageID} [get]
// @Security ApiKeyAuth
func (ctrl *ImageController) Serve(c *gin.Context) {
	albumID := c.Param("id")
	imageID := c.Param("imageID")
	imagePath, err := ctrl.service.ImageService.ServeImage(albumID, imageID)
	if err != nil {
		custom_errors.MapErrors(c, err)
		return
	}

	c.File(imagePath)
}

// Rename godoc
// @Summary Rename an image in an album
// @Description Renames a specific image in an album to a new valid name (alphanumeric, no path components)
// @Tags Image
// @Produce json
// @Param id path string true "Album ID"
// @Param imageID path string true "Image ID"
// @Param newName query string true "New name for the image (without extension)"
// @Success 200 {object} model.SuccessResponse{data=model.Image} "Image renamed successfully"
// @Failure 400 {object} model.ErrorDetails "Bad Request"
// @Failure 404 {object} model.ErrorDetails "Image Not Found"
// @Failure 409 {object} model.ErrorDetails "Conflict - Duplicate image name"
// @Failure 500 {object} model.ErrorDetails "Internal Server Error"
// @Router /v1/album/{id}/{imageID}/rename [patch]
// @Security ApiKeyAuth
func (ctrl *ImageController) Rename(c *gin.Context) {
	albumID := c.Param("id")
	imageID := c.Param("imageID")
	newName := c.DefaultQuery("newName", imageID)

	if newName == "" {
		customJson.ConstructJsonResponseError(c, custom_errors.NewError(custom_errors.ErrBadRequest, "New name is required"), http.StatusBadRequest)
		return
	}

	image, err := ctrl.service.ImageService.RenameImage(albumID, imageID, newName)
	if err != nil {
		custom_errors.MapErrors(c, err)
		return
	}
	customJson.ConstructJsonResponseSuccess(c, image, "Image renamed successfully")
}
