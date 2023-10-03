package entity

type Points struct {
	ID           uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	CheckingPeer string `json:"checkingpeer" form:"checkingpeer" csv:"checkingpeer" gorm:"column:checkingpeer"`
	CheckedPeer  string `json:"checkedpeer" form:"checkedpeer" csv:"checkedpeer" gorm:"column:checkedpeer"`
	PointsAmount int    `json:"pointsamount" form:"pointsamount" csv:"pointsamount" gorm:"column:pointsamount"`
}

func (Points) TableName() string {
	return "transferredpoints"
}

func (Points) GetRussianFieldNames() []string {
	return []string{"ID", "Проверяющий", "Проверяемый", "Количество поинтов"}
}
