package services

import "errors"

var ErrStorageDisabled = errors.New("storage service is not configured")

type StorageService interface {
	Enabled() bool
	PutObject(key string, data []byte, contentType string) (string, error)
	DeleteObject(key string) error
}

type noopStorageService struct{}

func NewNoopStorageService() StorageService {
	return noopStorageService{}
}

func (noopStorageService) Enabled() bool {
	return false
}

func (noopStorageService) PutObject(_ string, _ []byte, _ string) (string, error) {
	return "", ErrStorageDisabled
}

func (noopStorageService) DeleteObject(_ string) error {
	return ErrStorageDisabled
}
