package entity

type Points struct {
	ID           uint64 `json:"id,string" form:"id,string" csv:"id" gorm:"primaryKey"`
	CheckingPeer string `json:"checkingpeer" form:"checkingpeer" csv:"checkingpeer" gorm:"column:checkingpeer"`
	CheckedPeer  string `json:"checkedpeer" form:"checkedpeer" csv:"checkedpeer" gorm:"column:checkedpeer"`
	PointsAmount int    `json:"pointsamount,string" form:"pointsamount,string" csv:"pointsamount" gorm:"column:pointsamount;default:null"`
}

func (Points) TableName() string {
	return "transferredpoints"
}
