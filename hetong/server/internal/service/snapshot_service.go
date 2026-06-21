package service

import (
	"dashboard/internal/model"
	"dashboard/internal/repository"
)

type SnapshotService struct {
	snapshotRepo *repository.SnapshotRepository
}

func NewSnapshotService() *SnapshotService {
	return &SnapshotService{
		snapshotRepo: repository.NewSnapshotRepository(),
	}
}

func (s *SnapshotService) SaveSnapshot(req model.SnapshotRequest) (uint, error) {
	return s.snapshotRepo.SaveSnapshot(req)
}

func (s *SnapshotService) GetSnapshotList() ([]model.SnapshotItem, error) {
	return s.snapshotRepo.GetSnapshotList()
}

func (s *SnapshotService) GetSnapshotDetail(id uint) (*model.SnapshotDetailResponse, error) {
	return s.snapshotRepo.GetSnapshotDetail(id)
}

func (s *SnapshotService) DeleteSnapshot(id uint) error {
	return s.snapshotRepo.DeleteSnapshot(id)
}
