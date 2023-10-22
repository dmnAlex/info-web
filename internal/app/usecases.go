package app

import (
	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/usecase"
)

type UseCases struct {
	operations      *usecase.Operations
	checks          *usecase.DataSubsection[entity.Check]
	friends         *usecase.DataSubsection[entity.Friends]
	p2p             *usecase.DataSubsection[entity.P2P]
	peers           *usecase.DataSubsection[entity.Peer]
	recommendations *usecase.DataSubsection[entity.Recommendations]
	tasks           *usecase.DataSubsection[entity.Task]
	time_tracking   *usecase.DataSubsection[entity.TimeTracking]
	points          *usecase.DataSubsection[entity.Points]
	verter          *usecase.DataSubsection[entity.Verter]
	xp              *usecase.DataSubsection[entity.XP]
}

func initUseCases(rs *Repositories) *UseCases {
	return &UseCases{
		operations:      usecase.NewOperations(rs.operations),
		checks:          usecase.NewDataSubsection[entity.Check](rs.checks),
		friends:         usecase.NewDataSubsection[entity.Friends](rs.friends),
		p2p:             usecase.NewDataSubsection[entity.P2P](rs.p2p),
		peers:           usecase.NewDataSubsection[entity.Peer](rs.peers),
		recommendations: usecase.NewDataSubsection[entity.Recommendations](rs.recommendations),
		tasks:           usecase.NewDataSubsection[entity.Task](rs.tasks),
		time_tracking:   usecase.NewDataSubsection[entity.TimeTracking](rs.time_tracking),
		points:          usecase.NewDataSubsection[entity.Points](rs.points),
		verter:          usecase.NewDataSubsection[entity.Verter](rs.verter),
		xp:              usecase.NewDataSubsection[entity.XP](rs.xp),
	}
}
