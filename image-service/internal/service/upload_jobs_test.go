package service

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenUploadInputsSurviveMultipartCleanup(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "photo.jpg")
	require.NoError(t, err)
	content := bytes.Repeat([]byte("image-data"), 1024)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	request := httptest.NewRequest("POST", "/", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	require.NoError(t, request.ParseMultipartForm(1))
	files := request.MultipartForm.File["file"]
	files[0].Header.Set("Content-Type", "image/jpeg")

	inputs, err := openUploadInputs(files)
	require.NoError(t, err)
	require.NoError(t, request.MultipartForm.RemoveAll())

	retained, err := inputs[0].open()
	require.NoError(t, err)
	defer retained.Close()
	actual, err := io.ReadAll(retained)
	require.NoError(t, err)
	assert.Equal(t, content, actual)
}
