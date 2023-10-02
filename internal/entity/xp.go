package entity

type XP struct {
	ID       uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Check    uint64 `json:"check_id" form:"check_id" csv:"check_id" gorm:"column:Check"`
	XPAmount uint64 `json:"xp_amount" form:"xp_amount" csv:"xp_amount" gorm:"column:xpamount"`
}

func (XP) TableName() string {
  return "xp"
}
