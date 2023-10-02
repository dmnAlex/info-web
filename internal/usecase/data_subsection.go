package usecase

import "github.com/dsnikitin/info-web/internal/entity"

type DataSubsectionRepository[E entity.Entity] interface {
	Create(e *E) error
	ReadAll() []E
	Update(e *E)
	Delete(id string)
}

type DataSubsection[E entity.Entity] struct {
	repo DataSubsectionRepository[E]
}

func NewDataSubsection[E entity.Entity](r DataSubsectionRepository[E]) *DataSubsection[E] {
	return &DataSubsection[E]{repo: r}
}

func (uc *DataSubsection[E]) GetAll() []E {
	return uc.repo.ReadAll()
}

func (uc *DataSubsection[E]) Create(e *E) error {
	uc.repo.Create(e)
	// todo обработка ошибки
	return nil
}

func (uc *DataSubsection[E]) Update(e *E) error {
	uc.repo.Update(e)
	// todo обработка ошибки
	return nil
}

func (uc *DataSubsection[E]) Delete(id string) error {
	uc.repo.Delete(id)
	// todo обработка ошибки
	return nil
}
