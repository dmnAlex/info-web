package http

const (
	IndexEndpoint           = "/"
	OperationsEndpoint      = "/operations"
	DataEndpoint            = "/data"
	PeersEndpoint           = DataEndpoint + "/peers"
	FriendsEndpoint         = DataEndpoint + "/friends"
	RecommendationsEndpoint = DataEndpoint + "/recommendations"
	TasksEndpoint           = DataEndpoint + "/tasks"
	XPEndpont               = DataEndpoint + "/xp"
	PointsEndpoint          = DataEndpoint + "/points"
	ChecksEndpoint          = DataEndpoint + "/checks"
	P2PEndpoint             = DataEndpoint + "/p2p"
	VerterEndpoint          = DataEndpoint + "/verter"
	TimeTrackingEndpoint    = DataEndpoint + "/time_tracking"

	ExportPeersEndpoint           = PeersEndpoint + "/export"
	ExportFriendsEndpoint         = FriendsEndpoint + "/export"
	ExportRecommendationsEndpoint = RecommendationsEndpoint + "/export"
	ExportTasksEndpoint           = TasksEndpoint + "/export"
	ExportXPEndpoint              = XPEndpont + "/export"
	ExportPointsEndpoint          = PointsEndpoint + "/export"
	ExportChecksEndpoint          = ChecksEndpoint + "/export"
	ExportP2PEndpoint             = P2PEndpoint + "/export"
	ExportVerterEndpoint          = VerterEndpoint + "/export"
	ExportTimeTrackingEndpoint    = TimeTrackingEndpoint + "/export"
)
