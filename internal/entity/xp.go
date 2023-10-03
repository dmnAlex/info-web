package entity

type XP struct {
	ID       uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Check    uint64 `json:"check" form:"check" csv:"check" gorm:"column:Check"`
	XPAmount uint64 `json:"xpamount" form:"xpamount" csv:"xpamount" gorm:"column:xpamount"`
}

func (XP) TableName() string {
	return "xp"
}

func (XP) GetRussianFieldNames() []string {
	return []string{"ID", "ID проверки", "Количество опыта"}
}
