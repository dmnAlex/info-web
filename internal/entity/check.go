package entity

type Check struct {
	ID   uint64 `json:"id,string" form:"id,string" csv:"id" gorm:"primaryKey"`
	Peer string `json:"peer" form:"peer" csv:"peer"`
	Task string `json:"task" form:"task" csv:"task"`
	Date string `json:"date" form:"date" csv:"date" gorm:"column:Date;type:time;default:null"`
}
