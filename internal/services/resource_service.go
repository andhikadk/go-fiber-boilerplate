package services

import (
	"errors"

	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/models"
	"gorm.io/gorm"
)

var ErrResourceNotFound = errors.New("resource not found")

type ResourceService interface {
	ListResources(page, limit int) ([]dto.ResourceResponse, int64, error)
	GetResource(id uint) (*dto.ResourceResponse, error)
	CreateResource(userID uint, req *dto.CreateResourceRequest) (*dto.ResourceResponse, error)
	UpdateResource(id uint, req *dto.UpdateResourceRequest) (*dto.ResourceResponse, error)
	DeleteResource(id uint) error
}

type resourceService struct {
	db *gorm.DB
}

func NewResourceService(db *gorm.DB) ResourceService {
	return &resourceService{db: db}
}

func (s *resourceService) ListResources(page, limit int) ([]dto.ResourceResponse, int64, error) {
	var resources []models.Resource
	var total int64
	if err := s.db.Model(&models.Resource{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if err := s.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&resources).Error; err != nil {
		return nil, 0, err
	}
	out := make([]dto.ResourceResponse, 0, len(resources))
	for _, resource := range resources {
		out = append(out, toResourceResponse(&resource))
	}
	return out, total, nil
}

func (s *resourceService) GetResource(id uint) (*dto.ResourceResponse, error) {
	resource, err := s.findResource(id)
	if err != nil {
		return nil, err
	}
	resp := toResourceResponse(resource)
	return &resp, nil
}

func (s *resourceService) CreateResource(userID uint, req *dto.CreateResourceRequest) (*dto.ResourceResponse, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}
	resource := &models.Resource{
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
		CreatedByID: userID,
	}
	if err := s.db.Create(resource).Error; err != nil {
		return nil, err
	}
	resp := toResourceResponse(resource)
	return &resp, nil
}

func (s *resourceService) UpdateResource(id uint, req *dto.UpdateResourceRequest) (*dto.ResourceResponse, error) {
	resource, err := s.findResource(id)
	if err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) > 0 {
		if err := s.db.Model(resource).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetResource(id)
}

func (s *resourceService) DeleteResource(id uint) error {
	result := s.db.Delete(&models.Resource{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrResourceNotFound
	}
	return nil
}

func (s *resourceService) findResource(id uint) (*models.Resource, error) {
	var resource models.Resource
	if err := s.db.First(&resource, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	return &resource, nil
}

func toResourceResponse(resource *models.Resource) dto.ResourceResponse {
	return dto.ResourceResponse{
		ID:          resource.ID,
		Name:        resource.Name,
		Description: resource.Description,
		Status:      resource.Status,
		CreatedByID: resource.CreatedByID,
		CreatedAt:   resource.CreatedAt,
		UpdatedAt:   resource.UpdatedAt,
	}
}
