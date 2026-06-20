package main

import "time"

type DeviceType struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type Device struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	DeviceCode    string  `json:"device_code"`
	TypeID        int64   `json:"type_id"`
	TypeName      string  `json:"type_name"`
	Location      string  `json:"location"`
	Status        string  `json:"status"`
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryUsage   float64 `json:"memory_usage"`
	Temperature   float64 `json:"temperature"`
	ThresholdCPU  float64 `json:"threshold_cpu"`
	ThresholdTemp float64 `json:"threshold_temp"`
	LastHeartbeat string  `json:"last_heartbeat"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type Alarm struct {
	ID             int64  `json:"id"`
	DeviceID       int64  `json:"device_id"`
	DeviceName     string `json:"device_name"`
	Level          string `json:"level"`
	Message        string `json:"message"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	AcknowledgedAt string `json:"acknowledged_at"`
}

type Log struct {
	ID         int64  `json:"id"`
	DeviceID   int64  `json:"device_id"`
	DeviceName string `json:"device_name"`
	Level      string `json:"level"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
}

type Stats struct {
	Total    int `json:"total"`
	Online   int `json:"online"`
	Offline  int `json:"offline"`
	Abnormal int `json:"abnormal"`
	Alarms   int `json:"alarms"`
}

func now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
