package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/pkg/tools"
	"github.com/dsnikitin/info-web/internal/template"
	"github.com/gin-gonic/gin"
)

type OperationUseCase interface {
	GetAllOperations() []entity.Operation
	CallOperation(query string, arguments []interface{}) (entity.TableData, error)
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
		"endpoint":      OperationsEndpoint,
		"table_title":   "Operations",
		"table_data":    operations,
		"table_headers": tools.GetFieldNames(e),
		"table_type":    "function",
	}
	ctx.HTML(http.StatusOK, template.Operations, ginMap)
}

func (os *Operations) RawRequest(ctx *gin.Context) {
	var e entity.RequestData
	if err := ctx.Bind(&e); err != nil {
		// todo вернуть страницу bad request
		log.Println(err)
		return
	}

	placeholders := make([]string, len(e.Arguments))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("SELECT * FROM %s(%s)", e.FunctionName, strings.Join(placeholders, ", "))

	table, err := os.uc.CallOperation(query, e.Arguments)
	if err != nil {
		// todo вернуть страницу bad request
		log.Println(err)
		return
	}
	fmt.Println(table) // TODO

	ginMap := &gin.H{
		"table_title":   "Result",
		"table_data":    table.Rows,
		"table_headers": table.Headers,
	}
	ctx.HTML(http.StatusOK, template.Result, ginMap)
}
