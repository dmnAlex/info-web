package entity

type Task struct {
	ID         uint64 `json:"id" form:"id" csv:"id"`
	Title      string `json:"title" form:"title" csv:"title" gorm:"primaryKey"`
	ParentTask string `json:"parent_task" form:"parent_task" csv:"parent_task" gorm:"column:parenttask"`
	MaxXP      uint64 `json:"max_xp" form:"max_xp" csv:"max_xp" gorm:"column:maxxp"`
}
