package service

import (
	"dashboard/internal/model"
	"dashboard/internal/repository"
	"dashboard/pkg/drill"
)

type DrillService struct {
	drillEng   *drill.DrillEngine
	aggRepo    *repository.AggregateRepository
	dimRepo    *repository.DimensionRepository
	aggService *AggregateService
}

func NewDrillService() *DrillService {
	return &DrillService{
		drillEng:   drill.NewDrillEngine(),
		aggRepo:    repository.NewAggregateRepository(),
		dimRepo:    repository.NewDimensionRepository(),
		aggService: NewAggregateService(),
	}
}

func (s *DrillService) ExecuteDrill(req model.DrillRequest) (*model.DrillResponse, error) {
	var newPath []model.DrillLevel
	var nextDim model.DimensionLevel

	if req.DrillAction == "down" {
		label := s.getDimensionLabel(model.DimensionLevel(req.DrillDimension), req.DrillValue)
		newPath = s.drillEng.DrillDown(
			req.CurrentPath,
			model.DimensionLevel(req.DrillDimension),
			req.DrillValue,
			label,
		)
		nextDim, _ = s.drillEng.GetNextDimension(newPath)
	} else if req.DrillAction == "up" {
		newPath = s.drillEng.RollUpOne(req.CurrentPath)
		nextDim, _ = s.drillEng.GetNextDimension(newPath)
	}

	aggReq := model.AggregateRequest{
		Dimensions: s.getDimensionsForDrill(newPath),
		Metrics:    req.Metrics,
		DrillPath:  newPath,
	}

	aggResp, err := s.aggService.GetAggregatedData(aggReq)
	if err != nil {
		return nil, err
	}

	resp := &model.DrillResponse{}
	resp.Code = 0
	resp.Message = "success"
	resp.Data.DrillPath = newPath
	resp.Data.NextDimension = string(nextDim)
	resp.Data.Records = aggResp.Data.Records

	return resp, nil
}

func (s *DrillService) getDimensionLabel(dim model.DimensionLevel, value string) string {
	switch dim {
	case model.DimRegionLevel:
		if region, err := s.dimRepo.GetRegionByCode(value); err == nil {
			return region.RegionName
		}
	case model.DimBusinessLevel:
		if biz, err := s.dimRepo.GetBusinessByCode(value); err == nil {
			return biz.BusinessName
		}
	case model.DimTimeLevel:
		return value + "年"
	}
	return value
}

func (s *DrillService) getDimensionsForDrill(drillPath []model.DrillLevel) []model.DimensionLevel {
	usedDims := make(map[model.DimensionLevel]bool)
	for _, level := range drillPath {
		usedDims[level.Dimension] = true
	}

	allDims := []model.DimensionLevel{model.DimTimeLevel, model.DimRegionLevel, model.DimBusinessLevel}
	result := make([]model.DimensionLevel, 0)

	for _, dim := range allDims {
		if !usedDims[dim] {
			result = append(result, dim)
		}
	}

	if len(result) == 0 {
		return allDims
	}

	return result
}

func (s *DrillService) CanDrillDown(drillPath []model.DrillLevel) bool {
	return s.drillEng.CanDrillDown(drillPath)
}

func (s *DrillService) GetDrillPathSummary(drillPath []model.DrillLevel) string {
	return s.drillEng.GetPathSummary(drillPath)
}
