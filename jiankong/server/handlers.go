package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func newHandler() *Handler { return &Handler{} }

func jsonResponse(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(code, gin.H{"code": code, "message": msg, "data": data})
}

func (h *Handler) GetStats(c *gin.Context) {
	var s Stats
	db.QueryRow("SELECT COUNT(*) FROM devices").Scan(&s.Total)
	db.QueryRow("SELECT COUNT(*) FROM devices WHERE status='online'").Scan(&s.Online)
	db.QueryRow("SELECT COUNT(*) FROM devices WHERE status='offline'").Scan(&s.Offline)
	db.QueryRow("SELECT COUNT(*) FROM devices WHERE status='abnormal'").Scan(&s.Abnormal)
	db.QueryRow("SELECT COUNT(*) FROM alarms WHERE status='active'").Scan(&s.Alarms)
	jsonResponse(c, http.StatusOK, s, "ok")
}

func (h *Handler) GetDevices(c *gin.Context) {
	typeID := c.Query("device_type")
	status := c.Query("status")
	keyword := strings.TrimSpace(c.Query("keyword"))

	query := `SELECT d.id, d.name, d.device_code, d.type_id, t.name, d.location, d.status,
		d.cpu_usage, d.memory_usage, d.temperature, d.threshold_cpu, d.threshold_temp,
		d.last_heartbeat, d.created_at, d.updated_at
		FROM devices d LEFT JOIN device_types t ON d.type_id = t.id WHERE 1=1`
	args := []interface{}{}
	if typeID != "" {
		query += " AND d.type_id = ?"
		args = append(args, typeID)
	}
	if status != "" {
		query += " AND d.status = ?"
		args = append(args, status)
	}
	if keyword != "" {
		query += " AND (d.name LIKE ? OR d.device_code LIKE ? OR d.location LIKE ?)"
		kw := "%" + keyword + "%"
		args = append(args, kw, kw, kw)
	}
	query += " ORDER BY d.id ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		jsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	defer rows.Close()

	devices := []Device{}
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.Name, &d.DeviceCode, &d.TypeID, &d.TypeName, &d.Location,
			&d.Status, &d.CPUUsage, &d.MemoryUsage, &d.Temperature, &d.ThresholdCPU, &d.ThresholdTemp,
			&d.LastHeartbeat, &d.CreatedAt, &d.UpdatedAt); err != nil {
			continue
		}
		devices = append(devices, d)
	}
	jsonResponse(c, http.StatusOK, devices, "ok")
}

func (h *Handler) GetDevice(c *gin.Context) {
	id := c.Param("id")
	var d Device
	err := db.QueryRow(`SELECT d.id, d.name, d.device_code, d.type_id, t.name, d.location, d.status,
		d.cpu_usage, d.memory_usage, d.temperature, d.threshold_cpu, d.threshold_temp,
		d.last_heartbeat, d.created_at, d.updated_at
		FROM devices d LEFT JOIN device_types t ON d.type_id = t.id WHERE d.id = ?`, id).
		Scan(&d.ID, &d.Name, &d.DeviceCode, &d.TypeID, &d.TypeName, &d.Location,
			&d.Status, &d.CPUUsage, &d.MemoryUsage, &d.Temperature, &d.ThresholdCPU, &d.ThresholdTemp,
			&d.LastHeartbeat, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		jsonResponse(c, http.StatusNotFound, nil, "device not found")
		return
	}
	jsonResponse(c, http.StatusOK, d, "ok")
}

func (h *Handler) CreateDevice(c *gin.Context) {
	var d Device
	if err := c.ShouldBindJSON(&d); err != nil {
		jsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if d.ThresholdCPU == 0 {
		d.ThresholdCPU = 90
	}
	if d.ThresholdTemp == 0 {
		d.ThresholdTemp = 75
	}
	res, err := db.Exec(`INSERT INTO devices (name, device_code, type_id, location, status, cpu_usage, memory_usage, temperature, threshold_cpu, threshold_temp, last_heartbeat, created_at, updated_at)
		VALUES (?, ?, ?, ?, 'offline', 0, 0, 0, ?, ?, ?, ?, ?)`,
		d.Name, d.DeviceCode, d.TypeID, d.Location, d.ThresholdCPU, d.ThresholdTemp, now(), now(), now())
	if err != nil {
		jsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	d.ID = id
	d.Status = "offline"
	d.CreatedAt = now()
	d.UpdatedAt = now()
	db.Exec("INSERT INTO logs (device_id, level, message, created_at) VALUES (?, 'info', ?, ?)",
		id, "设备已录入系统", now())
	jsonResponse(c, http.StatusOK, d, "device created")
}

func (h *Handler) UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	var d Device
	if err := c.ShouldBindJSON(&d); err != nil {
		jsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	_, err := db.Exec(`UPDATE devices SET name=?, device_code=?, type_id=?, location=?, threshold_cpu=?, threshold_temp=?, updated_at=? WHERE id=?`,
		d.Name, d.DeviceCode, d.TypeID, d.Location, d.ThresholdCPU, d.ThresholdTemp, now(), id)
	if err != nil {
		jsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	db.Exec("INSERT INTO logs (device_id, level, message, created_at) VALUES (?, 'info', ?, ?)",
		id, "设备基础信息已更新", now())
	jsonResponse(c, http.StatusOK, gin.H{"id": id}, "device updated")
}

func (h *Handler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	db.Exec("DELETE FROM alarms WHERE device_id=?", id)
	db.Exec("DELETE FROM logs WHERE device_id=?", id)
	db.Exec("DELETE FROM devices WHERE id=?", id)
	jsonResponse(c, http.StatusOK, nil, "device deleted")
}

func (h *Handler) UpdateDeviceData(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		CPUUsage    float64 `json:"cpu_usage"`
		MemoryUsage float64 `json:"memory_usage"`
		Temperature float64 `json:"temperature"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		jsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	newStatus := detectStatus(id, body.CPUUsage, body.Temperature)
	_, err := db.Exec(`UPDATE devices SET cpu_usage=?, memory_usage=?, temperature=?, status=?, last_heartbeat=?, updated_at=? WHERE id=?`,
		body.CPUUsage, body.MemoryUsage, body.Temperature, newStatus, now(), now(), id)
	if err != nil {
		jsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	handleStatusChange(id, newStatus, body.CPUUsage, body.Temperature)
	var d Device
	db.QueryRow(`SELECT d.id, d.name, d.device_code, d.type_id, t.name, d.location, d.status,
		d.cpu_usage, d.memory_usage, d.temperature, d.threshold_cpu, d.threshold_temp,
		d.last_heartbeat, d.created_at, d.updated_at
		FROM devices d LEFT JOIN device_types t ON d.type_id = t.id WHERE d.id = ?`, id).
		Scan(&d.ID, &d.Name, &d.DeviceCode, &d.TypeID, &d.TypeName, &d.Location,
			&d.Status, &d.CPUUsage, &d.MemoryUsage, &d.Temperature, &d.ThresholdCPU, &d.ThresholdTemp,
			&d.LastHeartbeat, &d.CreatedAt, &d.UpdatedAt)
	jsonResponse(c, http.StatusOK, d, "data updated")
}

func (h *Handler) GetDeviceTypes(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, description, created_at FROM device_types ORDER BY id ASC")
	if err != nil {
		jsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	defer rows.Close()
	types := []DeviceType{}
	for rows.Next() {
		var t DeviceType
		rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt)
		types = append(types, t)
	}
	jsonResponse(c, http.StatusOK, types, "ok")
}

func (h *Handler) CreateDeviceType(c *gin.Context) {
	var t DeviceType
	if err := c.ShouldBindJSON(&t); err != nil {
		jsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	res, err := db.Exec("INSERT INTO device_types (name, description, created_at) VALUES (?, ?, ?)",
		t.Name, t.Description, now())
	if err != nil {
		jsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	t.ID, _ = res.LastInsertId()
	t.CreatedAt = now()
	jsonResponse(c, http.StatusOK, t, "type created")
}

func (h *Handler) DeleteDeviceType(c *gin.Context) {
	id := c.Param("id")
	var count int
	db.QueryRow("SELECT COUNT(*) FROM devices WHERE type_id=?", id).Scan(&count)
	if count > 0 {
		jsonResponse(c, http.StatusBadRequest, nil, fmt.Sprintf("该类型下仍有 %d 台设备，无法删除", count))
		return
	}
	db.Exec("DELETE FROM device_types WHERE id=?", id)
	jsonResponse(c, http.StatusOK, nil, "type deleted")
}

func (h *Handler) GetAlarms(c *gin.Context) {
	status := c.Query("status")
	query := `SELECT a.id, a.device_id, d.name, a.level, a.message, a.status, a.created_at, a.acknowledged_at
		FROM alarms a LEFT JOIN devices d ON a.device_id = d.id`
	args := []interface{}{}
	if status != "" {
		query += " WHERE a.status = ?"
		args = append(args, status)
	}
	query += " ORDER BY a.id DESC"
	if status == "" || status == "active" {
		query += " LIMIT 100"
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		jsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	defer rows.Close()
	alarms := []Alarm{}
	for rows.Next() {
		var a Alarm
		rows.Scan(&a.ID, &a.DeviceID, &a.DeviceName, &a.Level, &a.Message, &a.Status, &a.CreatedAt, &a.AcknowledgedAt)
		alarms = append(alarms, a)
	}
	jsonResponse(c, http.StatusOK, alarms, "ok")
}

func (h *Handler) AcknowledgeAlarm(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("UPDATE alarms SET status='acknowledged', acknowledged_at=? WHERE id=?", now(), id)
	if err != nil {
		jsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	var deviceID int64
	db.QueryRow("SELECT device_id FROM alarms WHERE id=?", id).Scan(&deviceID)
	db.Exec("INSERT INTO logs (device_id, level, message, created_at) VALUES (?, 'info', ?, ?)",
		deviceID, "告警已被运维人员确认处理", now())
	jsonResponse(c, http.StatusOK, nil, "alarm acknowledged")
}

func (h *Handler) GetLogs(c *gin.Context) {
	deviceID := c.Query("device_id")
	level := c.Query("level")
	limit := c.DefaultQuery("limit", "200")
	if n, err := strconv.Atoi(limit); err == nil && n > 0 && n <= 1000 {
		limit = strconv.Itoa(n)
	} else {
		limit = "200"
	}
	query := `SELECT l.id, l.device_id, d.name, l.level, l.message, l.created_at
		FROM logs l LEFT JOIN devices d ON l.device_id = d.id WHERE 1=1`
	args := []interface{}{}
	if deviceID != "" {
		query += " AND l.device_id = ?"
		args = append(args, deviceID)
	}
	if level != "" {
		query += " AND l.level = ?"
		args = append(args, level)
	}
	query += " ORDER BY l.id DESC LIMIT " + limit
	rows, err := db.Query(query, args...)
	if err != nil {
		jsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	defer rows.Close()
	logs := []Log{}
	for rows.Next() {
		var l Log
		rows.Scan(&l.ID, &l.DeviceID, &l.DeviceName, &l.Level, &l.Message, &l.CreatedAt)
		logs = append(logs, l)
	}
	jsonResponse(c, http.StatusOK, logs, "ok")
}

func (h *Handler) ClearLogs(c *gin.Context) {
	db.Exec("DELETE FROM logs")
	jsonResponse(c, http.StatusOK, nil, "logs cleared")
}
