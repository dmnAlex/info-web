package entity

type Recommendations struct {
	ID              uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Peer            string `json:"peer" form:"peer" csv:"peer"`
	RecommendedPeer string `json:"recommendedpeer" form:"recommendedpeer" csv:"recommendedpeer" gorm:"column:recommendedpeer"`
}
