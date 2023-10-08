package entity

type Friends struct {
	ID     uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Peer   string `json:"peer1" form:"peer1" csv:"peer1" gorm:"column:peer1"`
	Friend string `json:"peer2" form:"peer2" csv:"peer2" gorm:"column:peer2"`
}
