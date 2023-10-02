package repository

import (
	"github.com/dsnikitin/info-web/internal/entity"
	"gorm.io/gorm"
)

type DataManager[E entity.Entity] struct {
	db *gorm.DB
}

func NewDataManager[E entity.Entity](db *gorm.DB) *DataManager[E] {
	var entity E
	db.AutoMigrate(entity)
	return &DataManager[E]{db: db}
}

func (r *DataManager[E]) Create(e *E) error {
	if err := r.db.Create(e).Error; err != nil {
		return err
	}

	return nil
}

func (r *DataManager[E]) ReadAll() []E {
	var entities []E
	r.db.Find(&entities)
	return entities
}

func (r *DataManager[E]) Update(e *E) {
	r.db.Save(e)
}

func (r *DataManager[E]) Delete(id string) {
	var entity E
	r.db.Delete(entity, id)
}
