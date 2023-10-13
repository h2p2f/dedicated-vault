package servererrors

import "errors"

var (
	UserAlreadyExists = errors.New("user already exists")
	RecordNotFound    = errors.New("record not found")
)
