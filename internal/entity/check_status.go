package entity

type CheckStatus string

const (
	Start   CheckStatus = "Start"
	Success CheckStatus = "Success"
	Failure CheckStatus = "Failure"
)
