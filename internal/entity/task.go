package entity

type Task struct {
	Title      string `json:"title" form:"title" csv:"title" gorm:"primaryKey;default:null"`
	ParentTask string `json:"parenttask" form:"parenttask" csv:"parenttask" gorm:"column:parenttask"`
	MaxXP      uint64 `json:"maxxp,string" form:"maxxp,string" csv:"maxxp" gorm:"column:maxxp;default:null"`
}
