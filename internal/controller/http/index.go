package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (i *Index) GetIndexPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index", "")
}
