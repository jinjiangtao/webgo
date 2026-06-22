package database

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"whiteboard/internal/model"
)

var db *gorm.DB

func InitDB() error {
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(dataDir, "whiteboard.db")
	var err error
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	return db.AutoMigrate(
		&model.User{},
		&model.Whiteboard{},
		&model.Operation{},
		&model.Snapshot{},
	)
}

func GetDB() *gorm.DB {
	return db
}

type WhiteboardRepo struct{}

func NewWhiteboardRepo() *WhiteboardRepo {
	return &WhiteboardRepo{}
}

func (r *WhiteboardRepo) Create(w *model.Whiteboard) error {
	return db.Create(w).Error
}

func (r *WhiteboardRepo) GetByID(id uint) (*model.Whiteboard, error) {
	var w model.Whiteboard
	err := db.First(&w, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *WhiteboardRepo) List() ([]model.Whiteboard, error) {
	var list []model.Whiteboard
	err := db.Order("updated_at desc").Find(&list).Error
	return list, err
}

func (r *WhiteboardRepo) Update(w *model.Whiteboard) error {
	return db.Save(w).Error
}

func (r *WhiteboardRepo) Delete(id uint) error {
	return db.Delete(&model.Whiteboard{}, id).Error
}

type OperationRepo struct{}

func NewOperationRepo() *OperationRepo {
	return &OperationRepo{}
}

func (r *OperationRepo) CreateBatch(ops []model.Operation) error {
	if len(ops) == 0 {
		return nil
	}
	return db.Create(&ops).Error
}

func (r *OperationRepo) ListByWhiteboard(whiteboardID uint) ([]model.Operation, error) {
	var list []model.Operation
	err := db.Where("whiteboard_id = ?", whiteboardID).Order("timestamp asc").Find(&list).Error
	return list, err
}

func (r *OperationRepo) DeleteByWhiteboard(whiteboardID uint) error {
	return db.Where("whiteboard_id = ?", whiteboardID).Delete(&model.Operation{}).Error
}

type SnapshotRepo struct{}

func NewSnapshotRepo() *SnapshotRepo {
	return &SnapshotRepo{}
}

func (r *SnapshotRepo) Create(s *model.Snapshot) error {
	return db.Create(s).Error
}

func (r *SnapshotRepo) ListByWhiteboard(whiteboardID uint) ([]model.Snapshot, error) {
	var list []model.Snapshot
	err := db.Where("whiteboard_id = ?", whiteboardID).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *SnapshotRepo) GetLastByWhiteboard(whiteboardID uint) (*model.Snapshot, error) {
	var s model.Snapshot
	err := db.Where("whiteboard_id = ?", whiteboardID).Order("created_at desc").First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Create(u *model.User) error {
	return db.Create(u).Error
}

func (r *UserRepo) GetOrCreate(u *model.User) error {
	var existing model.User
	err := db.Where("username = ?", u.Username).First(&existing).Error
	if err == nil {
		*u = existing
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return db.Create(u).Error
}
