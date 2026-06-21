package handlers

import (
	"daka/server/database"
	"daka/server/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DayRecord struct {
	Date           string          `json:"date"`
	Records        []models.Record `json:"records"`
	TaskCount      int             `json:"task_count"`
	CompletedCount int             `json:"completed_count"`
}

func GetCalendarData(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	var records []models.Record
	database.DB.Preload("Task").
		Where("record_date >= ? AND record_date <= ?",
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02")).
		Order("record_date asc").
		Find(&records)

	calendar := make(map[string][]models.Record)
	for _, record := range records {
		calendar[record.RecordDate] = append(calendar[record.RecordDate], record)
	}

	var activeTasksCount int64
	database.DB.Model(&models.Task{}).Where("is_active = ?", true).Count(&activeTasksCount)

	result := make([]DayRecord, 0)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dayRecords := calendar[dateStr]
		if dayRecords == nil {
			dayRecords = make([]models.Record, 0)
		}
		completedCount := 0
		for _, r := range dayRecords {
			if r.Status != "absent" {
				completedCount++
			}
		}
		result = append(result, DayRecord{
			Date:           dateStr,
			Records:        dayRecords,
			TaskCount:      int(activeTasksCount),
			CompletedCount: completedCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"year":         year,
		"month":        month,
		"active_tasks": activeTasksCount,
		"calendar":     result,
	})
}

type MonthlyStats struct {
	Year           int     `json:"year"`
	Month          int     `json:"month"`
	TotalTasks     int64   `json:"total_tasks"`
	TotalCheckIns  int64   `json:"total_check_ins"`
	OnTimeCount    int64   `json:"on_time_count"`
	LateCount      int64   `json:"late_count"`
	MakeupCount    int64   `json:"makeup_count"`
	AbsentCount    int64   `json:"absent_count"`
	TotalDuration  int64   `json:"total_duration"`
	CompletionRate float64 `json:"completion_rate"`
	MaxStreak      int     `json:"max_streak"`
	CurrentStreak  int     `json:"current_streak"`
}

func GetMonthlyStats(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	var totalTasks int64
	database.DB.Model(&models.Task{}).Count(&totalTasks)

	var records []models.Record
	database.DB.Where("record_date >= ? AND record_date <= ?",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02")).Find(&records)

	var onTimeCount, lateCount, makeupCount, absentCount, totalDuration int64
	recordMap := make(map[string]int)

	for _, r := range records {
		switch r.Status {
		case "on_time", "checked_in":
			onTimeCount++
		case "late":
			lateCount++
		case "makeup":
			makeupCount++
		case "absent":
			absentCount++
		}
		totalDuration += int64(r.Duration)
		if r.Status != "absent" {
			recordMap[r.RecordDate]++
		}
	}

	totalCheckIns := int64(len(records)) - absentCount

	expectedDays := endDate.Day()
	var completionRate float64
	if totalTasks > 0 && expectedDays > 0 {
		completionRate = float64(totalCheckIns) / float64(totalTasks) / float64(expectedDays) * 100
		if completionRate > 100 {
			completionRate = 100
		}
	}

	maxStreak, currentStreak := calculateStreaks(recordMap, startDate, endDate)

	stats := MonthlyStats{
		Year:           year,
		Month:          month,
		TotalTasks:     totalTasks,
		TotalCheckIns:  totalCheckIns,
		OnTimeCount:    onTimeCount,
		LateCount:      lateCount,
		MakeupCount:    makeupCount,
		AbsentCount:    absentCount,
		TotalDuration:  totalDuration,
		CompletionRate: completionRate,
		MaxStreak:      maxStreak,
		CurrentStreak:  currentStreak,
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func calculateStreaks(recordMap map[string]int, startDate, endDate time.Time) (int, int) {
	maxStreak := 0
	currentStreak := 0
	today := time.Now().Format("2006-01-02")

	streak := 0
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		if recordMap[dateStr] > 0 {
			streak++
			if streak > maxStreak {
				maxStreak = streak
			}
		} else {
			streak = 0
		}
	}

	streak = 0
	for d := endDate; !d.Before(startDate); d = d.AddDate(0, 0, -1) {
		dateStr := d.Format("2006-01-02")
		if dateStr <= today && recordMap[dateStr] > 0 {
			streak++
		} else if dateStr <= today {
			break
		}
	}
	currentStreak = streak

	return maxStreak, currentStreak
}

func GetTaskStats(c *gin.Context) {
	taskID := c.Param("task_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	var records []models.Record
	database.DB.Where("task_id = ? AND record_date >= ? AND record_date <= ?",
		taskID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Find(&records)

	var onTime, late, makeup, absent, totalDuration int
	for _, r := range records {
		switch r.Status {
		case "on_time", "checked_in":
			onTime++
		case "late":
			late++
		case "makeup":
			makeup++
		case "absent":
			absent++
		}
		totalDuration += r.Duration
	}

	totalDays := endDate.Day()
	totalCompleted := onTime + late + makeup
	var rate float64
	if totalDays > 0 {
		rate = float64(totalCompleted) / float64(totalDays) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"task_id":         taskID,
			"total_days":      totalDays,
			"on_time":         onTime,
			"late":            late,
			"makeup":          makeup,
			"absent":          absent,
			"total_duration":  totalDuration,
			"completion_rate": rate,
		},
	})
}
