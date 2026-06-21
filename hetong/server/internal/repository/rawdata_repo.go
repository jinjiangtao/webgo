package repository

import (
	"dashboard/internal/model"
	"fmt"
	"strings"
	"time"
)

type RawDataRepository struct{}

func NewRawDataRepository() *RawDataRepository {
	return &RawDataRepository{}
}

func (r *RawDataRepository) QueryRawData(
	filters model.Filters,
	drillPath []model.DrillLevel,
	page int,
	pageSize int,
) ([]model.RawData, int64, error) {
	var args []interface{}

	whereClause, args := r.buildWhereClause(filters, drillPath, args)

	countSQL := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM raw_data r
		JOIN dim_time t ON r.date = t.date
		JOIN dim_region reg ON r.region_code = reg.region_code
		JOIN dim_business biz ON r.business_code = biz.business_code
		%s
	`, whereClause)

	var total int64
	err := DB.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	querySQL := fmt.Sprintf(`
		SELECT r.id, r.trace_no, r.date, r.region_code, r.business_code, 
			   r.order_no, r.amount, r.user_id, r.product_name, 
			   r.quantity, r.trade_time, r.extra, r.created_at
		FROM raw_data r
		JOIN dim_time t ON r.date = t.date
		JOIN dim_region reg ON r.region_code = reg.region_code
		JOIN dim_business biz ON r.business_code = biz.business_code
		%s
		ORDER BY r.trade_time DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	queryArgs := make([]interface{}, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := DB.Query(querySQL, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rawDataList []model.RawData
	for rows.Next() {
		var raw model.RawData
		var tradeTimeStr string
		var createdAtStr string
		err := rows.Scan(
			&raw.ID, &raw.TraceNo, &raw.Date, &raw.RegionCode, &raw.BusinessCode,
			&raw.OrderNo, &raw.Amount, &raw.UserID, &raw.ProductName,
			&raw.Quantity, &tradeTimeStr, &raw.Extra, &createdAtStr,
		)
		if err != nil {
			return nil, 0, err
		}

		raw.TradeTime, _ = time.Parse("2006-01-02 15:04:05", tradeTimeStr)
		raw.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)

		rawDataList = append(rawDataList, raw)
	}

	return rawDataList, total, nil
}

func (r *RawDataRepository) buildWhereClause(
	filters model.Filters,
	drillPath []model.DrillLevel,
	args []interface{},
) (string, []interface{}) {
	var conditions []string

	if len(filters.TimeRange) == 2 {
		conditions = append(conditions, "r.date BETWEEN ? AND ?")
		args = append(args, filters.TimeRange[0], filters.TimeRange[1])
	}
	if len(filters.Regions) > 0 {
		conditions = append(conditions, "reg.path LIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", filters.Regions[0]))
	}
	if len(filters.BusinessTypes) > 0 {
		placeholders := make([]string, len(filters.BusinessTypes))
		for i := range filters.BusinessTypes {
			placeholders[i] = "?"
			args = append(args, filters.BusinessTypes[i])
		}
		conditions = append(conditions, "biz.business_code IN ("+strings.Join(placeholders, ", ")+")")
	}

	for _, level := range drillPath {
		switch level.Dimension {
		case model.DimTimeLevel:
			conditions = append(conditions, "t.year = ?")
			args = append(args, level.Value)
		case model.DimRegionLevel:
			conditions = append(conditions, "reg.path LIKE ?")
			args = append(args, fmt.Sprintf("%%%s%%", level.Value))
		case model.DimBusinessLevel:
			conditions = append(conditions, "(biz.business_code = ? OR biz.parent_code = ?)")
			args = append(args, level.Value, level.Value)
		}
	}

	if len(conditions) > 0 {
		return "WHERE " + strings.Join(conditions, " AND "), args
	}
	return "", args
}

func (r *RawDataRepository) GetRawDataByTraceNo(traceNo string) (*model.RawData, error) {
	var raw model.RawData
	var tradeTimeStr string
	var createdAtStr string

	sql := `
		SELECT id, trace_no, date, region_code, business_code,
			   order_no, amount, user_id, product_name,
			   quantity, trade_time, extra, created_at
		FROM raw_data
		WHERE trace_no = ?
	`

	err := DB.QueryRow(sql, traceNo).Scan(
		&raw.ID, &raw.TraceNo, &raw.Date, &raw.RegionCode, &raw.BusinessCode,
		&raw.OrderNo, &raw.Amount, &raw.UserID, &raw.ProductName,
		&raw.Quantity, &tradeTimeStr, &raw.Extra, &createdAtStr,
	)
	if err != nil {
		return nil, err
	}

	raw.TradeTime, _ = time.Parse("2006-01-02 15:04:05", tradeTimeStr)
	raw.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)

	return &raw, nil
}
