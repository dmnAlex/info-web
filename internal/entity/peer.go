package entity

type Peer struct {
	Nickname string `json:"nickname" form:"nickname" csv:"nickname" gorm:"primaryKey"`
	Birthday string `json:"birthday" form:"birthday" csv:"birthday"`
}
