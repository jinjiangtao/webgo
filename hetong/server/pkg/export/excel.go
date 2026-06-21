package export

import (
	"dashboard/internal/model"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelExporter struct{}

func NewExcelExporter() *ExcelExporter {
	return &ExcelExporter{}
}

func (e *ExcelExporter) ExportTraceReport(
	traceData *model.TraceResponse,
	drillPath []model.DrillLevel,
	filters model.Filters,
	metrics []string,
) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	f.SetSheetName("Sheet1", "数据溯源报告")

	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#3B82F6"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "all", Color: "#CCCCCC", Style: 1},
		},
	})

	styleTitle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16, Color: "#1E3A8A"},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})

	f.SetCellValue("数据溯源报告", "A1", "企业数据可视化平台 - 数据溯源报告")
	f.MergeCell("数据溯源报告", "A1", "H1")
	f.SetCellStyle("数据溯源报告", "A1", "H1", styleTitle)

	f.SetCellValue("数据溯源报告", "A2", "生成时间")
	f.SetCellValue("数据溯源报告", "B2", time.Now().Format("2006-01-02 15:04:05"))

	f.SetCellValue("数据溯源报告", "A3", "溯源ID")
	f.SetCellValue("数据溯源报告", "B3", traceData.Data.TraceID)

	f.SetCellValue("数据溯源报告", "A4", "钻取路径")
	pathStr := ""
	for i, level := range drillPath {
		if i > 0 {
			pathStr += " → "
		}
		pathStr += fmt.Sprintf("%s: %s", level.Dimension, level.Label)
	}
	f.SetCellValue("数据溯源报告", "B4", pathStr)
	f.MergeCell("数据溯源报告", "B4", "H4")

	headerRow := 6
	headers := []string{"层级", "维度", "值", "标签"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c%d", 'A'+i, headerRow)
		f.SetCellValue("数据溯源报告", cell, header)
	}
	f.SetCellStyle("数据溯源报告", fmt.Sprintf("A%d", headerRow), fmt.Sprintf("%c%d", 'A'+len(headers)-1, headerRow), styleHeader)

	for i, path := range traceData.Data.AggregatePath {
		row := headerRow + 1 + i
		f.SetCellValue("数据溯源报告", fmt.Sprintf("A%d", row), path["level"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("B%d", row), path["dimension"])
		filter := path["filter"].(map[string]interface{})
		for _, v := range filter {
			f.SetCellValue("数据溯源报告", fmt.Sprintf("C%d", row), v)
		}
		f.SetCellValue("数据溯源报告", fmt.Sprintf("D%d", row), path["label"])
	}

	rawStartRow := headerRow + len(traceData.Data.AggregatePath) + 3
	f.SetCellValue("数据溯源报告", fmt.Sprintf("A%d", rawStartRow-1), "原始明细数据")
	f.SetCellStyle("数据溯源报告", fmt.Sprintf("A%d", rawStartRow-1), fmt.Sprintf("H%d", rawStartRow-1), styleTitle)

	rawHeaders := []string{"溯源编号", "订单号", "日期", "区域", "业务类型", "产品名称", "数量", "金额", "交易时间", "用户ID"}
	for i, header := range rawHeaders {
		cell := fmt.Sprintf("%c%d", 'A'+i, rawStartRow)
		f.SetCellValue("数据溯源报告", cell, header)
	}
	f.SetCellStyle("数据溯源报告", fmt.Sprintf("A%d", rawStartRow), fmt.Sprintf("%c%d", 'A'+len(rawHeaders)-1, rawStartRow), styleHeader)

	for i, raw := range traceData.Data.RawData {
		row := rawStartRow + 1 + i
		f.SetCellValue("数据溯源报告", fmt.Sprintf("A%d", row), raw["TraceNo"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("B%d", row), raw["OrderNo"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("C%d", row), raw["Date"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("D%d", row), raw["RegionCode"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("E%d", row), raw["BusinessCode"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("F%d", row), raw["ProductName"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("G%d", row), raw["Quantity"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("H%d", row), raw["Amount"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("I%d", row), raw["TradeTime"])
		f.SetCellValue("数据溯源报告", fmt.Sprintf("J%d", row), raw["UserID"])
	}

	for i := range rawHeaders {
		col := fmt.Sprintf("%c", 'A'+i)
		f.SetColWidth("数据溯源报告", col, col, 18)
	}

	filename := fmt.Sprintf("./data/trace_report_%s.xlsx", time.Now().Format("20060102150405"))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}

	return filename, nil
}

func (e *ExcelExporter) ExportAggregateReport(
	data []model.AggRecord,
	drillPath []model.DrillLevel,
	filters model.Filters,
	metrics []string,
) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	f.SetSheetName("Sheet1", "聚合数据报表")

	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 12, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#06B6D4"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "all", Color: "#CCCCCC", Style: 1},
		},
	})

	stylePositive, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#10B981"},
	})

	styleNegative, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#EF4444"},
	})

	headers := []string{"时间", "区域", "业务类型"}
	for _, m := range metrics {
		headers = append(headers, m)
		headers = append(headers, m+"同比(%)")
		headers = append(headers, m+"环比(%)")
	}
	headers = append(headers, "是否异常", "异常级别")

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("聚合数据报表", cell, header)
	}
	f.SetCellStyle("聚合数据报表", "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), styleHeader)

	for i, record := range data {
		row := i + 2
		dims := record.Dimensions
		f.SetCellValue("聚合数据报表", fmt.Sprintf("A%d", row), dims[model.DimTimeLevel])
		f.SetCellValue("聚合数据报表", fmt.Sprintf("B%d", row), dims[model.DimRegionLevel])
		f.SetCellValue("聚合数据报表", fmt.Sprintf("C%d", row), dims[model.DimBusinessLevel])

		col := 3
		for _, m := range metrics {
			metricVal := record.Metrics[m]
			f.SetCellValue("聚合数据报表", fmt.Sprintf("%c%d", 'A'+col, row), metricVal)
			col++

			if record.Comparison != nil {
				yoyCell := fmt.Sprintf("%c%d", 'A'+col, row)
				f.SetCellValue("聚合数据报表", yoyCell, record.Comparison.YoYPercent)
				if record.Comparison.YoYPercent >= 0 {
					f.SetCellStyle("聚合数据报表", yoyCell, yoyCell, stylePositive)
				} else {
					f.SetCellStyle("聚合数据报表", yoyCell, yoyCell, styleNegative)
				}
				col++

				momCell := fmt.Sprintf("%c%d", 'A'+col, row)
				f.SetCellValue("聚合数据报表", momCell, record.Comparison.MoMPercent)
				if record.Comparison.MoMPercent >= 0 {
					f.SetCellStyle("聚合数据报表", momCell, momCell, stylePositive)
				} else {
					f.SetCellStyle("聚合数据报表", momCell, momCell, styleNegative)
				}
				col++
			} else {
				col += 2
			}
		}

		if record.Anomaly != nil {
			f.SetCellValue("聚合数据报表", fmt.Sprintf("%c%d", 'A'+col, row), record.Anomaly.IsAnomaly)
			col++
			f.SetCellValue("聚合数据报表", fmt.Sprintf("%c%d", 'A'+col, row), record.Anomaly.Severity)
		}
	}

	for i := range headers {
		col := fmt.Sprintf("%c", 'A'+i)
		f.SetColWidth("聚合数据报表", col, col, 15)
	}

	filename := fmt.Sprintf("./data/aggregate_report_%s.xlsx", time.Now().Format("20060102150405"))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}

	return filename, nil
}
