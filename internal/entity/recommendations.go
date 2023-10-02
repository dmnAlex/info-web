package entity

type Recommendations struct {
	ID              uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Peer            string `json:"peer" form:"peer" csv:"peer"`
	RecommendedPeer string `json:"recommended_peer" form:"recommended_peer" csv:"recommended_peer" gorm:"column:recommendedpeer"`
}
