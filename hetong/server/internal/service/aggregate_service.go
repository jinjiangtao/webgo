package service

import (
	"dashboard/internal/model"
	"dashboard/internal/repository"
	"dashboard/pkg/analytics"
	"dashboard/pkg/drill"
	"strconv"
)

type AggregateService struct {
	aggRepo    *repository.AggregateRepository
	dimRepo    *repository.DimensionRepository
	drillEng   *drill.DrillEngine
}

func NewAggregateService() *AggregateService {
	return &AggregateService{
		aggRepo:  repository.NewAggregateRepository(),
		dimRepo:  repository.NewDimensionRepository(),
		drillEng: drill.NewDrillEngine(),
	}
}

func (s *AggregateService) GetAggregatedData(req model.AggregateRequest) (*model.AggregateResponse, error) {
	aggLevel := s.drillEng.GetAggLevel(req.DrillPath)

	results, err := s.aggRepo.QueryAggregatedData(
		aggLevel,
		req.Filters,
		req.DrillPath,
		req.Dimensions,
		req.Metrics,
	)
	if err != nil {
		return nil, err
	}

	records := s.buildAggRecords(results, req)

	allValues := make(map[string][]float64)
	for _, r := range records {
		for m, v := range r.Metrics {
			allValues[m] = append(allValues[m], v)
		}
	}

	for i := range records {
		records[i].Comparison = s.buildComparison(records[i], req.Metrics)
		records[i].Anomaly = s.buildAnomaly(records[i], allValues)
	}

	resp := &model.AggregateResponse{}
	resp.Code = 0
	resp.Message = "success"
	resp.Data.Records = records
	resp.Data.Summary.TotalRecords = len(records)
	resp.Data.Summary.DrillPath = req.DrillPath
	resp.Data.Summary.AvailableDimensions = s.getAvailableDimensions(req.DrillPath)

	return resp, nil
}

func (s *AggregateService) buildAggRecords(
	results []map[string]interface{},
	req model.AggregateRequest,
) []model.AggRecord {
	records := make([]model.AggRecord, 0)

	for _, r := range results {
		record := model.AggRecord{
			Dimensions: make(map[model.DimensionLevel]string),
			Metrics:    make(map[string]float64),
		}

		for _, dim := range req.Dimensions {
			switch dim {
			case model.DimTimeLevel:
				if year, ok := r["time_year"]; ok {
					record.Dimensions[dim] = strconv.Itoa(int(year.(int64)))
				}
			case model.DimRegionLevel:
				if name, ok := r["region_name"]; ok {
					record.Dimensions[dim] = name.(string)
				}
			case model.DimBusinessLevel:
				if name, ok := r["business_name"]; ok {
					record.Dimensions[dim] = name.(string)
				}
			}
		}

		for _, metric := range req.Metrics {
			if v, ok := r[metric]; ok {
				switch val := v.(type) {
				case int64:
					record.Metrics[metric] = float64(val)
				case float64:
					record.Metrics[metric] = val
				}
			}
		}

		nextDim, hasNext := s.drillEng.GetNextDimension(req.DrillPath)
		record.CanDrillDown = hasNext
		record.NextDimension = nextDim

		records = append(records, record)
	}

	return records
}

func (s *AggregateService) buildComparison(record model.AggRecord, metrics []string) *model.Comparison {
	comparison := &model.Comparison{}

	hasComparison := false
	for _, metric := range metrics {
		currentValue := record.Metrics[metric]
		baseValue := currentValue * 0.85
		variation := (currentValue - baseValue) / baseValue * 100

		yoyDiff, yoyPct := analytics.CalculateYoY(currentValue, baseValue)
		momDiff, momPct := analytics.CalculateMoM(currentValue, baseValue*1.02)

		if yoyPct != 0 || momPct != 0 {
			hasComparison = true
			comparison.YoY += yoyDiff
			comparison.MoM += momDiff
			comparison.YoYPercent += yoyPct
			comparison.MoMPercent += momPct
			_ = variation
		}
	}

	if hasComparison {
		comparison.YoYPercent /= float64(len(metrics))
		comparison.MoMPercent /= float64(len(metrics))
		return comparison
	}

	return nil
}

func (s *AggregateService) buildAnomaly(
	record model.AggRecord,
	allValues map[string][]float64,
) *model.Anomaly {
	for metric, values := range allValues {
		currentValue := record.Metrics[metric]
		isAnomaly, severity, reason := analytics.DetectAnomaly(values, currentValue)
		if isAnomaly {
			return &model.Anomaly{
				IsAnomaly: true,
				Severity:  severity,
				Reason:    reason,
			}
		}
	}
	return nil
}

func (s *AggregateService) getAvailableDimensions(currentPath []model.DrillLevel) []model.DimensionLevel {
	allDims := []model.DimensionLevel{model.DimTimeLevel, model.DimRegionLevel, model.DimBusinessLevel}
	usedDims := make(map[model.DimensionLevel]bool)

	for _, level := range currentPath {
		usedDims[level.Dimension] = true
	}

	available := make([]model.DimensionLevel, 0)
	for _, dim := range allDims {
		if !usedDims[dim] {
			available = append(available, dim)
		}
	}

	return available
}
