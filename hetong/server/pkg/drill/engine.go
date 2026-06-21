package drill

import (
	"dashboard/internal/model"
)

type DrillEngine struct {
	dimensionOrder []model.DimensionLevel
}

func NewDrillEngine() *DrillEngine {
	return &DrillEngine{
		dimensionOrder: []model.DimensionLevel{
			model.DimTimeLevel,
			model.DimRegionLevel,
			model.DimBusinessLevel,
		},
	}
}

func (e *DrillEngine) GetNextDimension(currentPath []model.DrillLevel) (model.DimensionLevel, bool) {
	usedDims := make(map[model.DimensionLevel]bool)
	for _, level := range currentPath {
		usedDims[level.Dimension] = true
	}

	for _, dim := range e.dimensionOrder {
		if !usedDims[dim] {
			return dim, true
		}
	}

	return "", false
}

func (e *DrillEngine) DrillDown(
	currentPath []model.DrillLevel,
	dimension model.DimensionLevel,
	value string,
	label string,
) []model.DrillLevel {
	newPath := make([]model.DrillLevel, len(currentPath))
	copy(newPath, currentPath)

	for i, level := range newPath {
		if level.Dimension == dimension {
			newPath[i].Value = value
			newPath[i].Label = label
			return newPath
		}
	}

	newPath = append(newPath, model.DrillLevel{
		Dimension: dimension,
		Value:     value,
		Label:     label,
	})

	return newPath
}

func (e *DrillEngine) RollUp(currentPath []model.DrillLevel, targetIndex int) []model.DrillLevel {
	if targetIndex < 0 || targetIndex >= len(currentPath) {
		return currentPath
	}
	return currentPath[:targetIndex+1]
}

func (e *DrillEngine) RollUpOne(currentPath []model.DrillLevel) []model.DrillLevel {
	if len(currentPath) <= 1 {
		return currentPath
	}
	return currentPath[:len(currentPath)-1]
}

func (e *DrillEngine) CanDrillDown(currentPath []model.DrillLevel) bool {
	_, hasNext := e.GetNextDimension(currentPath)
	return hasNext
}

func (e *DrillEngine) GetDrillLevel(currentPath []model.DrillLevel) int {
	return len(currentPath)
}

func (e *DrillEngine) GetAggLevel(currentPath []model.DrillLevel) int {
	return len(currentPath) + 1
}

func (e *DrillEngine) BuildAggregatePath(
	drillPath []model.DrillLevel,
) []map[string]interface{} {
	path := make([]map[string]interface{}, 0)

	for i, level := range drillPath {
		path = append(path, map[string]interface{}{
			"level":     i + 1,
			"dimension": string(level.Dimension),
			"value":     level.Value,
			"label":     level.Label,
			"filter": map[string]interface{}{
				string(level.Dimension): level.Value,
			},
		})
	}

	return path
}

func (e *DrillEngine) GetPathSummary(drillPath []model.DrillLevel) string {
	if len(drillPath) == 0 {
		return "全局视图"
	}

	summary := ""
	for i, level := range drillPath {
		if i > 0 {
			summary += " → "
		}
		summary += level.Label
	}
	return summary
}

func (e *DrillEngine) ValidateDrillAction(
	currentPath []model.DrillLevel,
	action string,
	targetDimension model.DimensionLevel,
) (bool, string) {
	switch action {
	case "down":
		usedDims := make(map[model.DimensionLevel]bool)
		for _, level := range currentPath {
			usedDims[level.Dimension] = true
		}
		if usedDims[targetDimension] {
			return false, "该维度已在钻取路径中"
		}
		return true, ""
	case "up":
		if len(currentPath) <= 1 {
			return false, "已在最顶层，无法上卷"
		}
		return true, ""
	default:
		return false, "不支持的钻取操作"
	}
}
