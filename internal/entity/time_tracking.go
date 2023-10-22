package entity

import "gorm.io/datatypes"

type TimeTracking struct {
	ID    uint64         `json:"id,string" form:"id,string" csv:"id" gorm:"primaryKey"`
	Peer  string         `json:"peer" form:"peer" csv:"peer"`
	Date  string         `json:"date" form:"date" csv:"date" gorm:"column:Date;default:null"`
	Time  datatypes.Time `json:"time" form:"time" csv:"time" gorm:"column:Time;default:null"`
	State int            `json:"state,string" form:"state,string" csv:"state"`
}

func (TimeTracking) TableName() string {
	return "timetracking"
}
