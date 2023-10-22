package entity

type Friends struct {
	ID    uint64 `json:"id,string" form:"id,string" csv:"id" gorm:"primaryKey"`
	Peer1 string `json:"peer1" form:"peer1" csv:"peer1"`
	Peer2 string `json:"peer2" form:"peer2" csv:"peer2"`
}
