package entity

type Points struct {
	ID           uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	CheckingPeer string `json:"checking_peer" form:"checking_peer" csv:"checking_peer" gorm:"column:checkingpeer"`
	CheckedPeer  string `json:"checked_peer" form:"checked_peer" csv:"checked_peer" gorm:"column:checkedpeer"`
	PointsAmount int    `json:"points_amount" form:"points_amount" csv:"points_amount" gorm:"column:pointsamount"`
}

func (Points) TableName() string {
	return "transferredpoints"
}
