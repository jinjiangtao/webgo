package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB(dbPath string) {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	db.SetMaxOpenConns(1)

	schema := `
	CREATE TABLE IF NOT EXISTS device_types (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		created_at TEXT DEFAULT ''
	);
	CREATE TABLE IF NOT EXISTS devices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		device_code TEXT NOT NULL UNIQUE,
		type_id INTEGER NOT NULL,
		location TEXT DEFAULT '',
		status TEXT DEFAULT 'offline',
		cpu_usage REAL DEFAULT 0,
		memory_usage REAL DEFAULT 0,
		temperature REAL DEFAULT 0,
		threshold_cpu REAL DEFAULT 90,
		threshold_temp REAL DEFAULT 75,
		last_heartbeat TEXT DEFAULT '',
		created_at TEXT DEFAULT '',
		updated_at TEXT DEFAULT ''
	);
	CREATE TABLE IF NOT EXISTS alarms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id INTEGER NOT NULL,
		level TEXT DEFAULT 'warning',
		message TEXT DEFAULT '',
		status TEXT DEFAULT 'active',
		created_at TEXT DEFAULT '',
		acknowledged_at TEXT DEFAULT ''
	);
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id INTEGER NOT NULL,
		level TEXT DEFAULT 'info',
		message TEXT DEFAULT '',
		created_at TEXT DEFAULT ''
	);
	`
	if _, err = db.Exec(schema); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	seedData()
}

func seedData() {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM device_types").Scan(&count)
	if count > 0 {
		return
	}

	types := []DeviceType{
		{Name: "服务器", Description: "机架式计算服务器"},
		{Name: "交换机", Description: "网络交换设备"},
		{Name: "路由器", Description: "网络路由设备"},
		{Name: "存储设备", Description: "数据存储阵列"},
		{Name: "温控设备", Description: "机房温控系统"},
	}
	for _, t := range types {
		db.Exec("INSERT INTO device_types (name, description, created_at) VALUES (?, ?, ?)", t.Name, t.Description, now())
	}

	devices := []struct {
		Name, Code, Location string
		TypeID               int
		ThresholdCPU         float64
		ThresholdTemp        float64
	}{
		{"核心数据库服务器", "SRV-001", "机房A-01", 1, 88, 72},
		{"应用服务器-02", "SRV-002", "机房A-02", 1, 90, 75},
		{"核心交换机-01", "SW-001", "机房A-网络柜", 2, 85, 70},
		{"边界路由器-01", "RT-001", "机房B-网络柜", 3, 85, 68},
		{"存储阵列-01", "STG-001", "机房B-01", 4, 88, 70},
		{"精密空调-01", "AC-001", "机房A", 5, 80, 65},
		{"精密空调-02", "AC-002", "机房B", 5, 80, 65},
		{"备份服务器-01", "SRV-003", "机房B-02", 1, 90, 75},
		{"汇聚交换机-02", "SW-002", "机房B-网络柜", 2, 85, 70},
		{"存储阵列-02", "STG-002", "机房B-03", 4, 88, 70},
		{"应用服务器-03", "SRV-004", "机房A-03", 1, 90, 75},
		{"边界路由器-02", "RT-002", "机房C-网络柜", 3, 85, 68},
	}

	for _, d := range devices {
		status := "online"
		cpu := 20 + rand.Float64()*40
		mem := 30 + rand.Float64()*40
		temp := 35 + rand.Float64()*20
		if rand.Float64() < 0.15 {
			status = "abnormal"
			cpu = d.ThresholdCPU + rand.Float64()*8
			temp = d.ThresholdTemp + rand.Float64()*8
		} else if rand.Float64() < 0.1 {
			status = "offline"
		}
		db.Exec(`INSERT INTO devices (name, device_code, type_id, location, status, cpu_usage, memory_usage, temperature, threshold_cpu, threshold_temp, last_heartbeat, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			d.Name, d.Code, d.TypeID, d.Location, status, cpu, mem, temp, d.ThresholdCPU, d.ThresholdTemp, now(), now(), now())

		var deviceID int64
		db.QueryRow("SELECT id FROM devices WHERE device_code = ?", d.Code).Scan(&deviceID)
		if status == "abnormal" {
			db.Exec("INSERT INTO alarms (device_id, level, message, status, created_at) VALUES (?, 'critical', ?, 'active', ?)",
				deviceID, fmt.Sprintf("设备 %s 检测到异常运行状态", d.Name), now())
			db.Exec("INSERT INTO logs (device_id, level, message, created_at) VALUES (?, 'error', ?, ?)",
				deviceID, "设备触发异常告警", now())
		} else {
			db.Exec("INSERT INTO logs (device_id, level, message, created_at) VALUES (?, 'info', ?, ?)",
				deviceID, "设备初始化完成，状态正常", now())
		}
	}
	log.Println("database seeded with sample data")
}
