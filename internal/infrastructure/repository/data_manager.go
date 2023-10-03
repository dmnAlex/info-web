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

func (r *DataManager[E]) Create(e *E) error {
	if err := r.db.Create(e).Error; err != nil {
		return err
	}

	return nil
}

func (r *DataManager[E]) ReadAll() []E {
	var es []E
	r.db.Find(&es)
	return es
}

func (r *DataManager[E]) Update(e *E) {
	r.db.Save(e)
}

func (r *DataManager[E]) Delete(id string) {
	var e E
	r.db.Where(tools.GetPrimaryKeyName(e)+" = ?", id).Delete(e)
}
