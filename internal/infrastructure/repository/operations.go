package repository

import (
	"fmt"

	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/queries"
	"gorm.io/gorm"
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

func (os *Operations) Call(query string, arguments []interface{}) (entity.TableData, error) {
	// var result map[string]interface{}
	// fmt.Println("calling:", query, arguments) // TODO
	// os.db.Raw(query, arguments...).Scan(&result)
	var table entity.TableData

	rows, err := os.db.Raw(query, arguments...).Rows()
	if err != nil {
		return table, err
	}
	defer rows.Close()

	table.Headers, err = rows.Columns()
	if err != nil {
		return table, err
	}

	size := len(table.Headers)

	for rows.Next() {
		values := make([]interface{}, size)
		valuePtrs := make([]interface{}, size)

		for i := range table.Headers {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		valueStrs := make([]string, size)

		for i, item := range values {
			valueStrs[i] = fmt.Sprint(item)
		}

		table.Rows = append(table.Rows, valueStrs)
	}

	return table, err
}
