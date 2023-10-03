package entity

type Task struct {
	Title      string `json:"title" form:"title" csv:"title" gorm:"primaryKey"`
	ParentTask string `json:"parenttask" form:"parenttask" csv:"parenttask" gorm:"column:parenttask"`
	MaxXP      uint64 `json:"maxxp" form:"maxxp" csv:"maxxp" gorm:"column:maxxp"`
}

func (Task) GetRussianFieldNames() []string {
	return []string{"Название", "Родительское задание", "Максимальный опыт"}
}
