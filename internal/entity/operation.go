package entity

import "github.com/lib/pq"

type Operation struct {
	ID         uint64         `json:"id" form:"id" csv:"id"`
	Name       string         `json:"name" form:"name" csv:"name"`
	Lngname    string         `json:"lngname" form:"lngname" csv:"lngname"`
	Kind       string         `json:"kind" form:"kind" csv:"kind"`
	Argnumber  uint64         `json:"argnumber" form:"argnumber" csv:"argnumber"`
	Returntype string         `json:"returntype" form:"returntype" csv:"returntype"`
	Inargs     pq.StringArray `json:"inargs" form:"inargs" csv:"inargs" gorm:"type:text[]"`
	Allargs    pq.StringArray `json:"allargs" form:"allargs" csv:"allargs" gorm:"type:text[]"`
	Argmodes   pq.StringArray `json:"argmodes" form:"argmodes" csv:"argmodes" gorm:"type:text[]"`
	Argnames   pq.StringArray `json:"argnames" form:"argnames" csv:"argnames" gorm:"type:text[]"`
}
