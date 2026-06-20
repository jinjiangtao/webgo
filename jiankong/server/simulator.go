package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func detectStatus(deviceID string, cpu, temp float64) string {
	var thresholdCPU, thresholdTemp float64
	db.QueryRow("SELECT threshold_cpu, threshold_temp FROM devices WHERE id=?", deviceID).Scan(&thresholdCPU, &thresholdTemp)
	if cpu >= thresholdCPU || temp >= thresholdTemp {
		return "abnormal"
	}
	return "online"
}

func handleStatusChange(deviceID string, newStatus string, cpu, temp float64) {
	var oldStatus, deviceName string
	db.QueryRow("SELECT status, name FROM devices WHERE id=?", deviceID).Scan(&oldStatus, &deviceName)
	if oldStatus == newStatus {
		return
	}
	var idInt int64
	fmt.Sscanf(deviceID, "%d", &idInt)

	switch newStatus {
	case "abnormal":
		var thresholdCPU, thresholdTemp float64
		db.QueryRow("SELECT threshold_cpu, threshold_temp FROM devices WHERE id=?", deviceID).Scan(&thresholdCPU, &thresholdTemp)
		msg := fmt.Sprintf("设备 %s 运行异常，CPU=%.1f%%，温度=%.1f°C", deviceName, cpu, temp)
		level := "warning"
		if cpu >= thresholdCPU && temp >= thresholdTemp {
			level = "critical"
		}
		db.Exec("INSERT INTO alarms (device_id, level, message, status, created_at) VALUES (?, ?, ?, 'active', ?)",
			idInt, level, msg, now())
		db.Exec("INSERT INTO logs (device_id, level, message, created_at) VALUES (?, 'error', ?, ?)",
			idInt, msg, now())
	case "online":
		db.Exec("INSERT INTO logs (device_id, level, message, created_at) VALUES (?, 'info', ?, ?)",
			idInt, fmt.Sprintf("设备 %s 已恢复正常运行", deviceName), now())
		db.Exec("UPDATE alarms SET status='resolved', acknowledged_at=? WHERE device_id=? AND status='active'",
			now(), idInt)
	case "offline":
		db.Exec("INSERT INTO logs (device_id, level, message, created_at) VALUES (?, 'warning', ?, ?)",
			idInt, fmt.Sprintf("设备 %s 心跳超时，状态置为离线", deviceName), now())
	}
}

func startSimulator() {
	go func() {
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()
		checkOfflineTicker := time.NewTicker(10 * time.Second)
		defer checkOfflineTicker.Stop()
		for {
			select {
			case <-ticker.C:
				simulateAllDevices()
			case <-checkOfflineTicker.C:
				markStaleDevicesOffline()
			}
		}
	}()
	log.Println("device data simulator started (interval: 4s)")
}

func simulateAllDevices() {
	rows, err := db.Query("SELECT id, status, cpu_usage, memory_usage, temperature, threshold_cpu, threshold_temp FROM devices")
	if err != nil {
		return
	}
	type devInfo struct {
		ID            int64
		Status        string
		CPU           float64
		Mem           float64
		Temp          float64
		ThresholdCPU  float64
		ThresholdTemp float64
	}
	devs := []devInfo{}
	for rows.Next() {
		var d devInfo
		rows.Scan(&d.ID, &d.Status, &d.CPU, &d.Mem, &d.Temp, &d.ThresholdCPU, &d.ThresholdTemp)
		devs = append(devs, d)
	}
	rows.Close()

	for _, d := range devs {
		if d.Status == "offline" {
			if rand.Float64() < 0.08 {
				db.Exec("UPDATE devices SET status='online', cpu_usage=20, memory_usage=30, temperature=40, last_heartbeat=?, updated_at=? WHERE id=?",
					now(), now(), d.ID)
				handleStatusChange(fmt.Sprintf("%d", d.ID), "online", 20, 40)
			}
			continue
		}

		cpuDelta := (rand.Float64() - 0.5) * 25
		newCPU := d.CPU + cpuDelta
		if newCPU < 5 {
			newCPU = 5
		}
		if newCPU > 99 {
			newCPU = 99
		}
		newMem := d.Mem + (rand.Float64()-0.5)*10
		if newMem < 10 {
			newMem = 10
		}
		if newMem > 98 {
			newMem = 98
		}
		newTemp := d.Temp + (rand.Float64()-0.4)*6
		if newTemp < 30 {
			newTemp = 30
		}
		if newTemp > 95 {
			newTemp = 95
		}

		if d.Status == "abnormal" && rand.Float64() < 0.25 {
			newCPU = d.ThresholdCPU - rand.Float64()*10
			newTemp = d.ThresholdTemp - rand.Float64()*10
		}

		newStatus := "online"
		if newCPU >= d.ThresholdCPU || newTemp >= d.ThresholdTemp {
			newStatus = "abnormal"
		}

		db.Exec("UPDATE devices SET cpu_usage=?, memory_usage=?, temperature=?, status=?, last_heartbeat=?, updated_at=? WHERE id=?",
			newCPU, newMem, newTemp, newStatus, now(), now(), d.ID)
		handleStatusChange(fmt.Sprintf("%d", d.ID), newStatus, newCPU, newTemp)
	}
}

func markStaleDevicesOffline() {
	cutoff := time.Now().Add(-15 * time.Second).Format("2006-01-02 15:04:05")
	rows, err := db.Query("SELECT id, name FROM devices WHERE status != 'offline' AND last_heartbeat < ?", cutoff)
	if err != nil {
		return
	}
	type info struct {
		ID   int64
		Name string
	}
	stale := []info{}
	for rows.Next() {
		var i info
		rows.Scan(&i.ID, &i.Name)
		stale = append(stale, i)
	}
	rows.Close()
	for _, s := range stale {
		db.Exec("UPDATE devices SET status='offline', updated_at=? WHERE id=?", now(), s.ID)
		handleStatusChange(fmt.Sprintf("%d", s.ID), "offline", 0, 0)
	}
}
