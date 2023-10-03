package entity

type P2P struct {
	ID           uint64      `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Check        uint64      `json:"check" form:"check" csv:"check"`
	CheckingPeer string      `json:"checkingpeer" form:"checkingpeer" csv:"checkingpeer" gorm:"column:checkingpeer"`
	State        CheckStatus `json:"state" form:"state" csv:"state" gorm:"type:checkstatus('start', 'success', 'failure')"`
	Time         string      `json:"time" form:"time" csv:"time"`
}

func (P2P) TableName() string {
	return "p2p"
}

func (P2P) GetRussianFieldNames() []string {
	return []string{"ID", "ID проверки", "Ник проверяющего", "Статус проверки", "Время"}
}
