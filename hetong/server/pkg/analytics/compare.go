package analytics

import "math"

func CalculateYoY(current float64, lastYear float64) (float64, float64) {
	if lastYear == 0 {
		return 0, 0
	}
	diff := current - lastYear
	percent := (diff / math.Abs(lastYear)) * 100
	return diff, roundFloat(percent, 2)
}

func CalculateMoM(current float64, lastMonth float64) (float64, float64) {
	if lastMonth == 0 {
		return 0, 0
	}
	diff := current - lastMonth
	percent := (diff / math.Abs(lastMonth)) * 100
	return diff, roundFloat(percent, 2)
}

func DetectAnomaly(values []float64, value float64) (bool, string, string) {
	if len(values) < 3 {
		return false, "", ""
	}

	mean := calculateMean(values)
	stdDev := calculateStdDev(values, mean)

	threshold := 2.0 * stdDev
	diffFromMean := math.Abs(value - mean)

	if diffFromMean > threshold {
		severity := "low"
		reason := "偏离平均值超过2个标准差"

		if diffFromMean > 3.0*stdDev {
			severity = "high"
			reason = "严重偏离平均值超过3个标准差"
		} else if diffFromMean > 2.5*stdDev {
			severity = "medium"
			reason = "中度偏离平均值超过2.5个标准差"
		}

		if value > mean {
			reason += "，数据异常偏高"
		} else {
			reason += "，数据异常偏低"
		}

		return true, severity, reason
	}

	return false, "", ""
}

func calculateMean(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateStdDev(values []float64, mean float64) float64 {
	sumSq := 0.0
	for _, v := range values {
		diff := v - mean
		sumSq += diff * diff
	}
	variance := sumSq / float64(len(values))
	return math.Sqrt(variance)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func CalculateAggregateMetrics(values []float64) map[string]float64 {
	if len(values) == 0 {
		return map[string]float64{
			"sum":   0,
			"avg":   0,
			"max":   0,
			"min":   0,
			"count": 0,
		}
	}

	sum := 0.0
	max := values[0]
	min := values[0]

	for _, v := range values {
		sum += v
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}

	return map[string]float64{
		"sum":   roundFloat(sum, 2),
		"avg":   roundFloat(sum/float64(len(values)), 2),
		"max":   max,
		"min":   min,
		"count": float64(len(values)),
	}
}
