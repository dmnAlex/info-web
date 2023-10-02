package entity

type P2P struct {
	ID           uint64      `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Check        uint64      `json:"check_id" form:"check_id" csv:"check_id"`
	CheckingPeer string      `json:"checking_peer" form:"checking_peer" csv:"checking_peer" gorm:"column:checkingpeer"`
	State        CheckStatus `json:"state" form:"state" csv:"state" gorm:"type:checkstatus('start', 'success', 'failure')"`
	Time         string      `json:"time" form:"time" csv:"time"`
}

func (P2P) TableName() string {
	return "p2p"
}
