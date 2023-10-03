package entity

type Friends struct {
	ID    uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Peer1 string `json:"peer1" form:"peer1" csv:"peer1"`
	Peer2 string `json:"peer2" form:"peer2" csv:"peer2"`
}

func (Friends) GetRussianFieldNames() []string {
	return []string{"ID", "Ник студента", "Ник друга"}
}
