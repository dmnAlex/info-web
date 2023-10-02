package entity

type Peer struct {
	ID       uint64 `json:"id" form:"id" csv:"id" gorm:"primaryKey"`
	Nickname string `json:"nickname" form:"nickname" csv:"nickname"`
	Birthday string `json:"birthday" form:"birthday" csv:"birthday"`
}
