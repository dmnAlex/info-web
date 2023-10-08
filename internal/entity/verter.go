package entity

type Verter struct {
	ID      uint64      `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	CheckID uint64      `json:"check" form:"check" csv:"check" gorm:"column:Check"`
	State   CheckStatus `json:"state" form:"state" csv:"state" gorm:"type:checkstatus('start', 'success', 'failure')"`
	Time    string      `json:"time" form:"time" csv:"time" gorm:"column:Time"`
}

func (Verter) TableName() string {
	return "verter"
}
