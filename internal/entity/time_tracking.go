package entity

type TimeTracking struct {
	ID    uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Peer  string `json:"peer" form:"peer" csv:"peer"`
	Date  string `json:"date" form:"date" csv:"date" gorm:"column:Date"`
	Time  string `json:"time" form:"time" csv:"time" gorm:"column:Time;type:time"`
	State int    `json:"state" form:"state" csv:"state"`
}

func (TimeTracking) TableName() string {
	return "timetracking"
}
