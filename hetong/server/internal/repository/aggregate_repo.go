package repository

import (
	"dashboard/internal/model"
	"fmt"
	"strings"
)

type AggregateRepository struct{}

func NewAggregateRepository() *AggregateRepository {
	return &AggregateRepository{}
}

func (r *AggregateRepository) QueryAggregatedData(
	aggLevel int,
	filters model.Filters,
	drillPath []model.DrillLevel,
	dimensions []model.DimensionLevel,
	metrics []string,
) ([]map[string]interface{}, error) {
	var args []interface{}

	selectClause := r.buildSelectClause(dimensions, metrics)
	whereClause, args := r.buildWhereClause(aggLevel, filters, drillPath, args)
	groupByClause := r.buildGroupByClause(dimensions)

	sql := fmt.Sprintf(`
		SELECT %s
		FROM agg_data a
		JOIN dim_time t ON a.date = t.date
		JOIN dim_region reg ON a.region_code = reg.region_code
		JOIN dim_business biz ON a.business_code = biz.business_code
		%s
		%s
	`, selectClause, whereClause, groupByClause)

	rows, err := DB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

func (r *AggregateRepository) buildSelectClause(dimensions []model.DimensionLevel, metrics []string) string {
	selectParts := make([]string, 0)

	for _, dim := range dimensions {
		switch dim {
		case model.DimTimeLevel:
			selectParts = append(selectParts, "t.year as time_year")
		case model.DimRegionLevel:
			selectParts = append(selectParts, "reg.region_code as region_code", "reg.region_name as region_name")
		case model.DimBusinessLevel:
			selectParts = append(selectParts, "biz.business_code as business_code", "biz.business_name as business_name")
		}
	}

	for _, metric := range metrics {
		switch metric {
		case "sales":
			selectParts = append(selectParts, "SUM(a.sales) as sales")
		case "orders":
			selectParts = append(selectParts, "SUM(a.orders) as orders")
		case "users":
			selectParts = append(selectParts, "SUM(a.users) as users")
		case "amount":
			selectParts = append(selectParts, "SUM(a.amount) as amount")
		}
	}

	return strings.Join(selectParts, ", ")
}

func (r *AggregateRepository) buildGroupByClause(dimensions []model.DimensionLevel) string {
	groupParts := make([]string, 0)
	for _, dim := range dimensions {
		switch dim {
		case model.DimTimeLevel:
			groupParts = append(groupParts, "t.year")
		case model.DimRegionLevel:
			groupParts = append(groupParts, "reg.region_code", "reg.region_name")
		case model.DimBusinessLevel:
			groupParts = append(groupParts, "biz.business_code", "biz.business_name")
		}
	}
	if len(groupParts) > 0 {
		return "GROUP BY " + strings.Join(groupParts, ", ")
	}
	return ""
}

func (r *AggregateRepository) buildWhereClause(
	aggLevel int,
	filters model.Filters,
	drillPath []model.DrillLevel,
	args []interface{},
) (string, []interface{}) {
	var conditions []string

	conditions = append(conditions, "a.agg_level = ?")
	args = append(args, aggLevel)

	if len(filters.TimeRange) == 2 {
		conditions = append(conditions, "a.date BETWEEN ? AND ?")
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

func (r *AggregateRepository) GetComparisonData(
	baseDate string,
	regionCode string,
	businessCode string,
	metrics []string,
) (map[string]float64, map[string]float64, error) {
	yoyData := make(map[string]float64)
	momData := make(map[string]float64)

	for _, metric := range metrics {
		var yoyValue, momValue float64

		yoySQL := fmt.Sprintf(`
			SELECT COALESCE(SUM(a.%s), 0) as val
			FROM agg_data a
			JOIN dim_time t ON a.date = t.date
			WHERE t.year = (SELECT year-1 FROM dim_time WHERE date = ?)
			AND a.region_code = ?
			AND a.business_code = ?
		`, metric)
		err := DB.QueryRow(yoySQL, baseDate, regionCode, businessCode).Scan(&yoyValue)
		if err != nil {
			return nil, nil, err
		}
		yoyData[metric] = yoyValue

		momSQL := fmt.Sprintf(`
			SELECT COALESCE(SUM(a.%s), 0) as val
			FROM agg_data a
			JOIN dim_time t ON a.date = t.date
			WHERE t.date < ?
			AND t.year = (SELECT year FROM dim_time WHERE date = ?)
			AND a.region_code = ?
			AND a.business_code = ?
			ORDER BY t.date DESC
			LIMIT 1
		`, metric)
		err = DB.QueryRow(momSQL, baseDate, baseDate, regionCode, businessCode).Scan(&momValue)
		if err != nil {
			return nil, nil, err
		}
		momData[metric] = momValue
	}

	return yoyData, momData, nil
}
