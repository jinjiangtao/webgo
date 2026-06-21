package model

import "time"

type DimTime struct {
	Date      string
	Year      int
	Quarter   int
	Month     int
	Week      int
	MonthName string
}

type DimRegion struct {
	RegionCode string
	RegionName string
	ParentCode *string
	Level      int
	Path       string
	Children   []DimRegion
}

type DimBusiness struct {
	BusinessCode string
	BusinessName string
	ParentCode   *string
	Level        int
	Category     string
	Children     []DimBusiness
}

type AggData struct {
	ID           uint
	Date         string
	RegionCode   string
	BusinessCode string
	Sales        float64
	Orders       int
	Users        int
	Amount       float64
	AggLevel     int
	CreatedAt    time.Time
}

type RawData struct {
	ID           uint
	TraceNo      string
	Date         string
	RegionCode   string
	BusinessCode string
	OrderNo      string
	Amount       float64
	UserID       string
	ProductName  string
	Quantity     int
	TradeTime    time.Time
	Extra        string
	CreatedAt    time.Time
}

type Snapshot struct {
	ID        uint
	Name      string
	State     string
	DrillPath string
	Filters   string
	CreatedAt time.Time
	CreatedBy string
}
