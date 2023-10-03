package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	dataEndPoint = "/data"
	dataTemplate = "data"

	// subsections
	peersEndpoint           = dataEndPoint + "/peers"
	peersTemplate           = "peers"
	friendsEndpoint         = dataEndPoint + "/friends"
	friendsTemplate         = "friends"
	recommendationsEndpoint = dataEndPoint + "/recommendations"
	recommendationsTemplate = "recommendations"
	tasksEndpoint           = dataEndPoint + "/tasks"
	tasksTemplate           = "tasks"
	xpEndpont               = dataEndPoint + "/xp"
	xpTemplate              = "xp"
	pointsEndpoint          = dataEndPoint + "/points"
	pointsTemplate          = "points"
	checksEndpoint          = dataEndPoint + "/checks"
	checksTemplate          = "checks"
	p2pEndpoint             = dataEndPoint + "/p2p"
	p2pTemplate             = "p2p"
	verterEndpoint          = dataEndPoint + "/verter"
	verterTemplate          = "verter"
	timeTrackingEndpoint    = dataEndPoint + "/time_tracking"
	timeTrackingTemplate    = "time_tracking"

	// common
	tableTemplate     = "data_table"
	errDialogTemplate = "error_dialog"
)

type Data struct {
}

func NewData() *Data {
	return &Data{}
}

func (d *Data) GetFrontPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, dataTemplate, "")
}
