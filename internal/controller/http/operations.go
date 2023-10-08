package http

import (
	"net/http"

	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/pkg/tools"
	"github.com/gin-gonic/gin"
)

const (
	operationsEndpoint = "/operations"
)

type OperationUseCase interface {
	GetAllOperations() []entity.Operation
}
type Operations struct {
	uc OperationUseCase
}

func NewOperations(uc OperationUseCase) *Operations {
	return &Operations{uc: uc}
}

func (os *Operations) GetAll(ctx *gin.Context) {
	var e entity.Operation
	operations := os.uc.GetAllOperations()
	ginMap := &gin.H{
		"endpoint":      operationsEndpoint,
		"table_title":   "Operations",
		"table_data":    operations,
		"table_headers": tools.GetFieldNames(e),
		"table_type":    "function",
	}
	ctx.HTML(http.StatusOK, "operations", ginMap)
}
