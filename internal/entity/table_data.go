package entity

type TableData struct {
	Headers []string   `json:"headers" csv:"headers"`
	Rows    [][]string `json:"data" csv:"data"`
}
