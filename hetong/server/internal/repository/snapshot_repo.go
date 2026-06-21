package repository

import (
	"dashboard/internal/model"
	"encoding/json"
	"time"
)

type SnapshotRepository struct{}

func NewSnapshotRepository() *SnapshotRepository {
	return &SnapshotRepository{}
}

func (r *SnapshotRepository) SaveSnapshot(req model.SnapshotRequest) (uint, error) {
	stateJSON, _ := json.Marshal(req.State)
	drillPathJSON, _ := json.Marshal(req.DrillPath)
	filtersJSON, _ := json.Marshal(req.Filters)

	sql := `
		INSERT INTO snapshots (name, state, drill_path, filters, created_by)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := DB.Exec(sql, req.Name, string(stateJSON), string(drillPathJSON), string(filtersJSON), req.CreatedBy)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (r *SnapshotRepository) GetSnapshotList() ([]model.SnapshotItem, error) {
	sql := `
		SELECT id, name, created_at, created_by
		FROM snapshots
		ORDER BY created_at DESC
	`

	rows, err := DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.SnapshotItem, 0)
	for rows.Next() {
		var s model.SnapshotItem
		var createdAtStr string
		err := rows.Scan(&s.ID, &s.Name, &createdAtStr, &s.CreatedBy)
		if err != nil {
			return nil, err
		}
		s.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		result = append(result, s)
	}

	return result, nil
}

func (r *SnapshotRepository) GetSnapshotDetail(id uint) (*model.SnapshotDetailResponse, error) {
	sql := `
		SELECT id, name, state, drill_path, filters, created_at, created_by
		FROM snapshots
		WHERE id = ?
	`

	var snapshot struct {
		ID        uint
		Name      string
		State     string
		DrillPath string
		Filters   string
		CreatedAt string
		CreatedBy string
	}

	err := DB.QueryRow(sql, id).Scan(
		&snapshot.ID, &snapshot.Name, &snapshot.State, &snapshot.DrillPath,
		&snapshot.Filters, &snapshot.CreatedAt, &snapshot.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	var state map[string]interface{}
	var drillPath []model.DrillLevel
	var filters model.Filters

	json.Unmarshal([]byte(snapshot.State), &state)
	json.Unmarshal([]byte(snapshot.DrillPath), &drillPath)
	json.Unmarshal([]byte(snapshot.Filters), &filters)

	createdAt, _ := time.Parse("2006-01-02 15:04:05", snapshot.CreatedAt)

	resp := &model.SnapshotDetailResponse{}
	resp.Code = 0
	resp.Message = "success"
	resp.Data.ID = snapshot.ID
	resp.Data.Name = snapshot.Name
	resp.Data.State = state
	resp.Data.DrillPath = drillPath
	resp.Data.Filters = filters
	resp.Data.CreatedAt = createdAt
	resp.Data.CreatedBy = snapshot.CreatedBy

	return resp, nil
}

func (r *SnapshotRepository) DeleteSnapshot(id uint) error {
	sql := `DELETE FROM snapshots WHERE id = ?`
	_, err := DB.Exec(sql, id)
	return err
}
