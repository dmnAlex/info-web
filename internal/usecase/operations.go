package usecase

import "github.com/dsnikitin/info-web/internal/entity"

type OperationsRepo interface {
	Get() []entity.Operation
	Call(query string, arguments []interface{}) (entity.TableData, error)
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

func (os *Operations) CallOperation(query string, arguments []interface{}) (entity.TableData, error) {
	return os.repo.Call(query, arguments)
}
