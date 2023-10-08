package repository

import (
	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/pkg/tools"
	"gorm.io/gorm"
)

type DataManager[E entity.Entity] struct {
	db *gorm.DB
}

func NewDataManager[E entity.Entity](db *gorm.DB) *DataManager[E] {
	var e E
	db.AutoMigrate(e)
	return &DataManager[E]{db: db}
}

func (m *DataManager[E]) Create(e *E) error {
	if err := m.db.Create(e).Error; err != nil {
		return err
	}

	return nil
}

func (m *DataManager[E]) ReadAll() ([]E, error) {
	var entities []E
	if err := m.db.Find(&entities).Error; err != nil {
		return entities, err
	}

	return entities, nil
}

func (m *DataManager[E]) Update(e *E) error {
	if err := m.db.Save(e).Error; err != nil {
		return err
	}

	return nil
}

func (m *DataManager[E]) Delete(id string) error {
	var e E
	if err := m.db.Where(tools.GetPrimaryKeyName(e)+" = ?", id).Delete(e).Error; err != nil {
		return err
	}

	return nil
}
