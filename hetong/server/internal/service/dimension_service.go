package service

import (
	"dashboard/internal/repository"
)

type DimensionService struct {
	dimRepo *repository.DimensionRepository
}

func NewDimensionService() *DimensionService {
	return &DimensionService{
		dimRepo: repository.NewDimensionRepository(),
	}
}

func (s *DimensionService) GetAllDimensions() (
	[]map[string]interface{},
	[]map[string]interface{},
	[]map[string]interface{},
	error,
) {
	timeOpts, err := s.dimRepo.GetTimeOptions()
	if err != nil {
		return nil, nil, nil, err
	}

	regionOpts, err := s.dimRepo.GetRegionTree()
	if err != nil {
		return nil, nil, nil, err
	}

	bizOpts, err := s.dimRepo.GetBusinessTree()
	if err != nil {
		return nil, nil, nil, err
	}

	return timeOpts, regionOpts, bizOpts, nil
}
