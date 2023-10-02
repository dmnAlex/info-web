package usecase

import "github.com/dsnikitin/info-web/internal/entity"

type OperationsRepo interface {
	Get() []entity.Operation
}

type Operations struct {
	repo OperationsRepo
}

func NewOperations(r OperationsRepo) *Operations {
	return &Operations{repo: r}
}

func (os *Operations) GetAllOperations() []entity.Operation {
	return os.repo.Get()
}
