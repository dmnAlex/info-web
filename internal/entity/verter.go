package entity

type Verter struct {
	ID    uint64      `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Check uint64      `json:"checking_id" form:"checking_id" csv:"checking_id" gorm:"column:Date"`
	State CheckStatus `json:"state" form:"state" csv:"state" gorm:"type:checkstatus('start', 'success', 'failure')"`
	Time  string      `json:"time" form:"time" csv:"time" gorm:"column:Time"`
}

func (Verter) TableName() string {
	return "verter"
}
