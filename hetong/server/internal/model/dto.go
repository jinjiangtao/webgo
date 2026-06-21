package model

import "time"

type DimensionLevel string

const (
	DimTimeLevel     DimensionLevel = "time"
	DimRegionLevel   DimensionLevel = "region"
	DimBusinessLevel DimensionLevel = "business"
)

type DrillLevel struct {
	Dimension DimensionLevel `json:"dimension"`
	Value     string         `json:"value"`
	Label     string         `json:"label"`
}

type Filters struct {
	TimeRange     []string `json:"timeRange,omitempty"`
	Regions       []string `json:"regions,omitempty"`
	BusinessTypes []string `json:"businessTypes,omitempty"`
}

type AggregateRequest struct {
	Dimensions []DimensionLevel `json:"dimensions"`
	Metrics    []string         `json:"metrics"`
	Filters    Filters          `json:"filters"`
	DrillPath  []DrillLevel     `json:"drillPath"`
}

type Comparison struct {
	YoY        float64 `json:"yoy"`
	MoM        float64 `json:"mom"`
	YoYPercent float64 `json:"yoyPercent"`
	MoMPercent float64 `json:"momPercent"`
}

type Anomaly struct {
	IsAnomaly bool   `json:"isAnomaly"`
	Severity  string `json:"severity"`
	Reason    string `json:"reason"`
}

type AggRecord struct {
	Dimensions    map[DimensionLevel]string `json:"dimensions"`
	Metrics       map[string]float64        `json:"metrics"`
	Comparison    *Comparison               `json:"comparison,omitempty"`
	Anomaly       *Anomaly                  `json:"anomaly,omitempty"`
	CanDrillDown  bool                      `json:"canDrillDown"`
	NextDimension DimensionLevel            `json:"nextDimension"`
}

type AggregateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Records []AggRecord `json:"records"`
		Summary struct {
			TotalRecords       int              `json:"totalRecords"`
			DrillPath          []DrillLevel     `json:"drillPath"`
			AvailableDimensions []DimensionLevel `json:"availableDimensions"`
		} `json:"summary"`
	} `json:"data"`
}

type DrillRequest struct {
	DrillAction    string       `json:"drillAction"`
	CurrentPath    []DrillLevel `json:"currentPath"`
	DrillDimension string       `json:"drillDimension"`
	DrillValue     string       `json:"drillValue"`
	Metrics        []string     `json:"metrics"`
}

type DrillResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		DrillPath     []DrillLevel  `json:"drillPath"`
		NextDimension string        `json:"nextDimension"`
		Records       []AggRecord   `json:"records"`
	} `json:"data"`
}

type TraceRequest struct {
	DrillPath []DrillLevel `json:"drillPath"`
	Filters   Filters      `json:"filters"`
	Page      int          `json:"page"`
	PageSize  int          `json:"pageSize"`
}

type TraceResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TraceID       string                   `json:"traceId"`
		AggregatePath []map[string]interface{} `json:"aggregatePath"`
		RawData       []map[string]interface{} `json:"rawData"`
		TotalRaw      int64                    `json:"totalRaw"`
		Page          int                      `json:"page"`
		PageSize      int                      `json:"pageSize"`
	} `json:"data"`
}

type DimensionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TimeOptions     []map[string]interface{} `json:"timeOptions"`
		RegionOptions   []map[string]interface{} `json:"regionOptions"`
		BusinessOptions []map[string]interface{} `json:"businessOptions"`
	} `json:"data"`
}

type SnapshotRequest struct {
	Name      string                 `json:"name"`
	State     map[string]interface{} `json:"state"`
	DrillPath []DrillLevel           `json:"drillPath"`
	Filters   Filters                `json:"filters"`
	CreatedBy string                 `json:"createdBy"`
}

type SnapshotItem struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
}

type SnapshotListResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []SnapshotItem `json:"data"`
}

type SnapshotDetailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID        uint                   `json:"id"`
		Name      string                 `json:"name"`
		State     map[string]interface{} `json:"state"`
		DrillPath []DrillLevel           `json:"drillPath"`
		Filters   Filters                `json:"filters"`
		CreatedAt time.Time              `json:"createdAt"`
		CreatedBy string                 `json:"createdBy"`
	} `json:"data"`
}

type ExportRequest struct {
	DrillPath []DrillLevel `json:"drillPath"`
	Filters   Filters      `json:"filters"`
	Metrics   []string     `json:"metrics"`
	Format    string       `json:"format"`
}
