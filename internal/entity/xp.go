package entity

type XP struct {
	ID       uint64 `json:"id,string" form:"id,string" csv:"id" gorm:"primaryKey"`
	CheckID  uint64 `json:"checkid,string" form:"checkid,string" csv:"check" gorm:"column:Check;default:null"`
	XPAmount uint64 `json:"xpamount,string" form:"xpamount,string" csv:"xpamount" gorm:"column:xpamount;default:null"`
}

func (XP) TableName() string {
	return "xp"
}
