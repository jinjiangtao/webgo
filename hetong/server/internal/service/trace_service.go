package service

import (
	"dashboard/internal/model"
	"dashboard/internal/repository"
	"dashboard/pkg/drill"

	"github.com/google/uuid"
)

type TraceService struct {
	rawRepo  *repository.RawDataRepository
	drillEng *drill.DrillEngine
}

func NewTraceService() *TraceService {
	return &TraceService{
		rawRepo:  repository.NewRawDataRepository(),
		drillEng: drill.NewDrillEngine(),
	}
}

func (s *TraceService) GetTraceData(req model.TraceRequest) (*model.TraceResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	rawData, total, err := s.rawRepo.QueryRawData(
		req.Filters,
		req.DrillPath,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		return nil, err
	}

	rawDataList := make([]map[string]interface{}, 0)
	for _, raw := range rawData {
		rawDataList = append(rawDataList, map[string]interface{}{
			"ID":           raw.ID,
			"TraceNo":      raw.TraceNo,
			"Date":         raw.Date,
			"RegionCode":   raw.RegionCode,
			"BusinessCode": raw.BusinessCode,
			"OrderNo":      raw.OrderNo,
			"Amount":       raw.Amount,
			"UserID":       raw.UserID,
			"ProductName":  raw.ProductName,
			"Quantity":     raw.Quantity,
			"TradeTime":    raw.TradeTime.Format("2006-01-02 15:04:05"),
		})
	}

	aggregatePath := s.drillEng.BuildAggregatePath(req.DrillPath)

	resp := &model.TraceResponse{}
	resp.Code = 0
	resp.Message = "success"
	resp.Data.TraceID = uuid.New().String()
	resp.Data.AggregatePath = aggregatePath
	resp.Data.RawData = rawDataList
	resp.Data.TotalRaw = total
	resp.Data.Page = req.Page
	resp.Data.PageSize = req.PageSize

	return resp, nil
}

func (s *TraceService) GetRawDataByTraceNo(traceNo string) (*model.RawData, error) {
	return s.rawRepo.GetRawDataByTraceNo(traceNo)
}
