package entity

import "errors"

var (
	ErrPeerNotFound      = errors.New("peer doesn't exists")
	ErrPeerAlreadyExists = errors.New("peer with such nickname already exists")
)

func GetErrorDescription(err error) string {
	switch err {
	case ErrPeerAlreadyExists:
		return "Студент с таким ником уже существует"
	case ErrPeerNotFound:
		return "Студента с таким ником не существует"
	default:
		return "Произошла неизвестная ошибка"
	}
}
