package usecase

import "github.com/dsnikitin/info-web/internal/entity"

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
