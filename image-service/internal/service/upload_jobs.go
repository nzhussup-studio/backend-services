package service

import (
	"fmt"
	"image-service/internal/cache"
	custom_errors "image-service/internal/errors"
	"image-service/internal/model"
	"image-service/internal/repository"
	"mime/multipart"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UploadJob struct {
	ID        string         `json:"id"`
	AlbumID   string         `json:"-"`
	Status    string         `json:"status"`
	Total     int            `json:"total"`
	Completed int            `json:"completed"`
	Images    []*model.Image `json:"images,omitempty"`
	Error     string         `json:"error,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}

type uploadJobService struct {
	imageService *ImageService
	jobs         map[string]*UploadJob
	mu           sync.RWMutex
	queue        chan struct{}
}

func NewUploadJobService(storage *repository.Storage, redis cache.RedisClientInterface, validate *validator.Validate, cfg ImageConfig) *uploadJobService {
	return &uploadJobService{
		imageService: &ImageService{storage: storage, redis: redis, validate: validate, cfg: cfg},
		jobs:         make(map[string]*UploadJob),
		queue:        make(chan struct{}, 1),
	}
}

func (s *uploadJobService) StartUpload(albumID string, files []*multipart.FileHeader) (*UploadJob, error) {
	if err := validateUpload(files, s.imageService.cfg); err != nil {
		return nil, err
	}

	select {
	case s.queue <- struct{}{}:
	default:
		return nil, custom_errors.NewError(custom_errors.ErrConflict, "another upload is already processing; try again when it finishes")
	}

	job := &UploadJob{
		ID:        uuid.NewString(),
		AlbumID:   albumID,
		Status:    "queued",
		Total:     len(files),
		CreatedAt: time.Now().UTC(),
	}
	s.mu.Lock()
	s.jobs[job.ID] = job
	s.mu.Unlock()

	go func() {
		defer func() { <-s.queue }()
		s.update(job.ID, func(current *UploadJob) {
			current.Status = "processing"
		})
		images, err := s.imageService.uploadImage(albumID, files, func() {
			s.update(job.ID, func(current *UploadJob) {
				current.Completed++
			})
		})
		if err != nil {
			s.update(job.ID, func(current *UploadJob) {
				current.Status = "failed"
				current.Error = errorMessage(err)
			})
			return
		}
		s.update(job.ID, func(current *UploadJob) {
			current.Status = "completed"
			current.Images = images
		})
	}()

	return job, nil
}

func (s *uploadJobService) GetUploadStatus(albumID string, jobID string) (*UploadJob, error) {
	s.mu.RLock()
	job, ok := s.jobs[jobID]
	if !ok || job.AlbumID != albumID {
		s.mu.RUnlock()
		return nil, custom_errors.NewError(custom_errors.ErrNotFound, "upload job not found")
	}
	copy := *job
	s.mu.RUnlock()
	return &copy, nil
}

func (s *uploadJobService) update(jobID string, update func(*UploadJob)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[jobID]; ok {
		update(job)
	}
}

func validateUpload(files []*multipart.FileHeader, cfg ImageConfig) error {
	if len(files) == 0 {
		return custom_errors.NewError(custom_errors.ErrBadRequest, "no files uploaded")
	}
	if cfg.MaxUploadCount > 0 && len(files) > cfg.MaxUploadCount {
		return custom_errors.NewError(custom_errors.ErrBadRequest, fmt.Sprintf("too many images: maximum is %d", cfg.MaxUploadCount))
	}
	var total int64
	for _, file := range files {
		if file.Size > 0 && file.Size > cfg.MaxUploadBytes {
			return custom_errors.NewError(custom_errors.ErrBadRequest, fmt.Sprintf("image file too large: maximum is %d MB", cfg.MaxUploadBytes/1024/1024))
		}
		total += file.Size
	}
	if cfg.MaxTotalUploadBytes > 0 && total > cfg.MaxTotalUploadBytes {
		return custom_errors.NewError(custom_errors.ErrBadRequest, fmt.Sprintf("total upload too large: maximum is %d MB", cfg.MaxTotalUploadBytes/1024/1024))
	}
	return nil
}

func errorMessage(err error) string {
	return err.Error()
}
