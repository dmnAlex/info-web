package entity

type Friends struct {
	ID    uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Peer1 string `json:"peer_1" form:"peer_1" csv:"peer_1"`
	Peer2 string `json:"peer_2" form:"peer_2" csv:"peer_2"`
}
