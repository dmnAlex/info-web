package app

import (
	"github.com/dsnikitin/info-web/internal/controller/http"
	"github.com/dsnikitin/info-web/internal/entity"
)

type Handlers struct {
	index           *http.Index
	data            *http.Data
	operations      *http.Operations
	checks          *http.DataSubsection[entity.Check]
	friends         *http.DataSubsection[entity.Friends]
	p2p             *http.DataSubsection[entity.P2P]
	peers           *http.DataSubsection[entity.Peer]
	recommendations *http.DataSubsection[entity.Recommendations]
	tasks           *http.DataSubsection[entity.Task]
	time_tracking   *http.DataSubsection[entity.TimeTracking]
	points          *http.DataSubsection[entity.Points]
	verter          *http.DataSubsection[entity.Verter]
	xp              *http.DataSubsection[entity.XP]
}

func initHandlers(uc *UseCases) *Handlers {
	return &Handlers{
		index:           http.NewIndex(),
		data:            http.NewData(),
		operations:      http.NewOperations(uc.operations),
		checks:          http.NewDataSubsection[entity.Check](uc.checks),
		friends:         http.NewDataSubsection[entity.Friends](uc.friends),
		p2p:             http.NewDataSubsection[entity.P2P](uc.p2p),
		peers:           http.NewDataSubsection[entity.Peer](uc.peers),
		recommendations: http.NewDataSubsection[entity.Recommendations](uc.recommendations),
		tasks:           http.NewDataSubsection[entity.Task](uc.tasks),
		time_tracking:   http.NewDataSubsection[entity.TimeTracking](uc.time_tracking),
		points:          http.NewDataSubsection[entity.Points](uc.points),
		verter:          http.NewDataSubsection[entity.Verter](uc.verter),
		xp:              http.NewDataSubsection[entity.XP](uc.xp),
	}
}
