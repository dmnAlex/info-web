package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	dataEndPoint            = "/data"
	peersEndpoint           = dataEndPoint + "/peers"
	friendsEndpoint         = dataEndPoint + "/friends"
	recommendationsEndpoint = dataEndPoint + "/recommendations"
	tasksEndpoint           = dataEndPoint + "/tasks"
	xpEndpont               = dataEndPoint + "/xp"
	pointsEndpoint          = dataEndPoint + "/points"
	checksEndpoint          = dataEndPoint + "/checks"
	p2pEndpoint             = dataEndPoint + "/p2p"
	verterEndpoint          = dataEndPoint + "/verter"
	timeTrackingEndpoint    = dataEndPoint + "/time_tracking"

	dataTemplate       = "data"
	subsectionTemplate = "data_subsection"
	tableTemplate      = "data_table"
	errDialogTemplate  = "error_dialog"
)

type Data struct {
}

func NewData() *Data {
	return &Data{}
}

func (d *Data) GetFrontPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, dataTemplate, "")
}
