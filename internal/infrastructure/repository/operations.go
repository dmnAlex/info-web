package repository

import (
	"gorm.io/gorm"
	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/queries"
)

type Operations struct {
	db *gorm.DB
}

func NewOperations(db *gorm.DB) *Operations {
	return &Operations{db: db}
}

func (os *Operations) Get() []entity.Operation {
	var operations []entity.Operation
	os.db.Raw(queries.GetFunctionData).Scan(&operations)
	return operations
}
