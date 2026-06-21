package service

import (
	"dashboard/internal/model"
	"dashboard/pkg/export"
)

type ExportService struct {
	excelExporter *export.ExcelExporter
	traceService  *TraceService
	aggService    *AggregateService
}

func NewExportService() *ExportService {
	return &ExportService{
		excelExporter: export.NewExcelExporter(),
		traceService:  NewTraceService(),
		aggService:    NewAggregateService(),
	}
}

func (s *ExportService) ExportExcel(req model.ExportRequest) (string, error) {
	traceReq := model.TraceRequest{
		DrillPath: req.DrillPath,
		Filters:   req.Filters,
		Page:      1,
		PageSize:  1000,
	}

	traceData, err := s.traceService.GetTraceData(traceReq)
	if err != nil {
		return "", err
	}

	filename, err := s.excelExporter.ExportTraceReport(
		traceData,
		req.DrillPath,
		req.Filters,
		req.Metrics,
	)

	return filename, err
}

func (s *ExportService) ExportAggregateExcel(req model.ExportRequest) (string, error) {
	aggReq := model.AggregateRequest{
		Dimensions: []model.DimensionLevel{model.DimTimeLevel, model.DimRegionLevel, model.DimBusinessLevel},
		Metrics:    req.Metrics,
		Filters:    req.Filters,
		DrillPath:  req.DrillPath,
	}

	aggData, err := s.aggService.GetAggregatedData(aggReq)
	if err != nil {
		return "", err
	}

	filename, err := s.excelExporter.ExportAggregateReport(
		aggData.Data.Records,
		req.DrillPath,
		req.Filters,
		req.Metrics,
	)

	return filename, err
}
