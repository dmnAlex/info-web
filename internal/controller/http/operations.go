package http

import (
	"net/http"

	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/gin-gonic/gin"
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
	operations := os.uc.GetAllOperations()
	ginMap := &gin.H{
		"endpoint":      peersEndpoint,
		"object_fields": []string{"id", "name", "lngname", "kind", "argnumber", "returntype", "inargs", "allargs", "argmodes", "argnames"},
		"table_title":   "Операции",
		"table_data":    operations,
		"table_headers": []string{"OID", "Название", "Язык", "Тип", "Агрументы", "Тип возврата", "Типы сигнатуры", "Типы всех аргументов", "Режим аргументов", "Название агрументов"},
		"table_type":    "function",
	}
	ctx.HTML(http.StatusOK, "operations", ginMap)
}
