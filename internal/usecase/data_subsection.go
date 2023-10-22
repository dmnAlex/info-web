package usecase

import (
	"fmt"

	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/pkg/tools"
)

type DataSubsectionRepository[E entity.Entity] interface {
	Create(e *E) error
	ReadAll() ([]E, error)
	Update(e *E) error
	Delete(id string) error
}

type DataSubsection[E entity.Entity] struct {
	repo DataSubsectionRepository[E]
}

func NewDataSubsection[E entity.Entity](r DataSubsectionRepository[E]) *DataSubsection[E] {
	return &DataSubsection[E]{repo: r}
}

func (uc *DataSubsection[E]) GetAll() ([]E, error) {
	return uc.repo.ReadAll()
}

func (uc *DataSubsection[E]) Create(e *E) error {
	return uc.repo.Create(e)
}

func (uc *DataSubsection[E]) Update(e *E) error {
	return uc.repo.Update(e)
}

func (uc *DataSubsection[E]) Delete(id string) error {
	return uc.repo.Delete(id)
}

func (uc *DataSubsection[E]) Export() ([][]string, error) {
	entities, err := uc.repo.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	var e E
	data := make([][]string, 0, len(entities)+1)
	data = append(data, tools.GetFieldNames(e)) // Header
	for _, entity := range entities {
		values := tools.GetFieldValues(entity)
		valuesStr := make([]string, 0, len(values))
		for i := range values {
			valuesStr = append(valuesStr, fmt.Sprintf("%v", values[i]))
		}
		data = append(data, valuesStr) // Items
	}

	return data, nil
}
