package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}
	var result map[string]interface{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}
	*j = result
	return nil
}

type EditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ImageUUID  string    `gorm:"size:36;index;not null" json:"image_uuid"`
	ActionType string    `gorm:"size:50;not null" json:"action_type"`
	ActionData JSONMap   `gorm:"type:json" json:"action_data"`
	CreatedAt  time.Time `json:"created_at"`
}
