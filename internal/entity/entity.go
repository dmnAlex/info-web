package entity

type Entity interface {
	Peer | Friends | Recommendations | Task | XP | Points | Check | P2P | Verter | TimeTracking
}

type Tabler interface {
	TableName() string
}
