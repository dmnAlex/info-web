package entity

import (
	"gorm.io/datatypes"
)

type Verter struct {
	ID      uint64         `json:"id,string" form:"id,string" csv:"id" gorm:"primaryKey"`
	CheckID uint64         `json:"checkid,string" form:"checkid,string" csv:"check" gorm:"column:Check"`
	State   CheckStatus    `json:"state" form:"state" csv:"state" gorm:"type:checkstatus('start', 'success', 'failure');default:null"`
	Time    datatypes.Time `json:"time" form:"time" csv:"time" gorm:"column:Time;default:null"`
}

func (Verter) TableName() string {
	return "verter"
}
