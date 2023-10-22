package app

import (
	"github.com/dsnikitin/info-web/internal/entity"
	"github.com/dsnikitin/info-web/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	operations      *repository.Operations
	checks          *repository.DataManager[entity.Check]
	friends         *repository.DataManager[entity.Friends]
	p2p             *repository.DataManager[entity.P2P]
	peers           *repository.DataManager[entity.Peer]
	recommendations *repository.DataManager[entity.Recommendations]
	tasks           *repository.DataManager[entity.Task]
	time_tracking   *repository.DataManager[entity.TimeTracking]
	points          *repository.DataManager[entity.Points]
	verter          *repository.DataManager[entity.Verter]
	xp              *repository.DataManager[entity.XP]
}

func initRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		operations:      repository.NewOperations(db),
		checks:          repository.NewDataManager[entity.Check](db),
		friends:         repository.NewDataManager[entity.Friends](db),
		p2p:             repository.NewDataManager[entity.P2P](db),
		peers:           repository.NewDataManager[entity.Peer](db),
		recommendations: repository.NewDataManager[entity.Recommendations](db),
		tasks:           repository.NewDataManager[entity.Task](db),
		time_tracking:   repository.NewDataManager[entity.TimeTracking](db),
		points:          repository.NewDataManager[entity.Points](db),
		verter:          repository.NewDataManager[entity.Verter](db),
		xp:              repository.NewDataManager[entity.XP](db),
	}
}
