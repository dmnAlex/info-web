package app

import (
	"github.com/dsnikitin/info-web/internal/template"
	"github.com/gin-gonic/gin"

	"github.com/dsnikitin/info-web/internal/controller/http"
)

func initRoutes(r *gin.Engine, h *Handlers) {
	r.NoRoute(func(ctx *gin.Context) { ctx.HTML(404, template.PageNotFound, "") })

	r.GET(http.IndexEndpoint, h.index.GetIndexPage)
	r.GET(http.DataEndpoint, h.data.GetFrontPage)

	r.GET(http.OperationsEndpoint, h.operations.GetAll)
	r.POST(http.OperationsEndpoint, h.operations.RawRequest)

	r.GET(http.ChecksEndpoint, h.checks.GetAllEntities)
	r.GET(http.ExportChecksEndpoint, h.checks.ExportData)
	r.POST(http.ChecksEndpoint, h.checks.CreateEntity)
	r.PUT(http.ChecksEndpoint, h.checks.UpdateEntity)
	r.DELETE(http.ChecksEndpoint, h.checks.DeleteEntity)

	r.GET(http.FriendsEndpoint, h.friends.GetAllEntities)
	r.GET(http.ExportFriendsEndpoint, h.friends.ExportData)
	r.POST(http.FriendsEndpoint, h.friends.CreateEntity)
	r.PUT(http.FriendsEndpoint, h.friends.UpdateEntity)
	r.DELETE(http.FriendsEndpoint, h.friends.DeleteEntity)

	r.GET(http.P2PEndpoint, h.p2p.GetAllEntities)
	r.GET(http.ExportP2PEndpoint, h.p2p.ExportData)
	r.POST(http.P2PEndpoint, h.p2p.CreateEntity)
	r.PUT(http.P2PEndpoint, h.p2p.UpdateEntity)
	r.DELETE(http.P2PEndpoint, h.p2p.DeleteEntity)

	r.GET(http.PeersEndpoint, h.peers.GetAllEntities)
	r.GET(http.ExportPeersEndpoint, h.peers.ExportData)
	r.POST(http.PeersEndpoint, h.peers.CreateEntity)
	r.PUT(http.PeersEndpoint, h.peers.UpdateEntity)
	r.DELETE(http.PeersEndpoint, h.peers.DeleteEntity)

	r.GET(http.RecommendationsEndpoint, h.recommendations.GetAllEntities)
	r.GET(http.ExportRecommendationsEndpoint, h.recommendations.ExportData)
	r.POST(http.RecommendationsEndpoint, h.recommendations.CreateEntity)
	r.PUT(http.RecommendationsEndpoint, h.recommendations.UpdateEntity)
	r.DELETE(http.RecommendationsEndpoint, h.recommendations.DeleteEntity)

	r.GET(http.TasksEndpoint, h.tasks.GetAllEntities)
	r.GET(http.ExportTasksEndpoint, h.tasks.ExportData)
	r.POST(http.TasksEndpoint, h.tasks.CreateEntity)
	r.PUT(http.TasksEndpoint, h.tasks.UpdateEntity)
	r.DELETE(http.TasksEndpoint, h.tasks.DeleteEntity)

	r.GET(http.TimeTrackingEndpoint, h.time_tracking.GetAllEntities)
	r.GET(http.ExportTimeTrackingEndpoint, h.time_tracking.ExportData)
	r.POST(http.TimeTrackingEndpoint, h.time_tracking.CreateEntity)
	r.PUT(http.TimeTrackingEndpoint, h.time_tracking.UpdateEntity)
	r.DELETE(http.TimeTrackingEndpoint, h.time_tracking.DeleteEntity)

	r.GET(http.PointsEndpoint, h.points.GetAllEntities)
	r.GET(http.ExportPointsEndpoint, h.points.ExportData)
	r.POST(http.PointsEndpoint, h.points.CreateEntity)
	r.PUT(http.PointsEndpoint, h.points.UpdateEntity)
	r.DELETE(http.PointsEndpoint, h.points.DeleteEntity)

	r.GET(http.VerterEndpoint, h.verter.GetAllEntities)
	r.GET(http.ExportVerterEndpoint, h.verter.ExportData)
	r.POST(http.VerterEndpoint, h.verter.CreateEntity)
	r.PUT(http.VerterEndpoint, h.verter.UpdateEntity)
	r.DELETE(http.VerterEndpoint, h.verter.DeleteEntity)

	r.GET(http.XPEndpont, h.xp.GetAllEntities)
	r.GET(http.ExportXPEndpoint, h.xp.ExportData)
	r.POST(http.XPEndpont, h.xp.CreateEntity)
	r.PUT(http.XPEndpont, h.xp.UpdateEntity)
	r.DELETE(http.XPEndpont, h.xp.DeleteEntity)
}
