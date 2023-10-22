package entity

type RequestData struct {
	FunctionName string        `json:"functionname" csv:"functionname"`
	Arguments    []interface{} `json:"arguments" csv:"arguments"`
}
