package http

import (
	"net/http"

	"github.com/dsnikitin/info-web/internal/template"
	"github.com/gin-gonic/gin"
)

type Data struct {
}

func NewData() *Data {
	return &Data{}
}

func (d *Data) GetFrontPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, template.Data, "")
}
