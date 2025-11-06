package repositories

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrAlreadyExists  = errors.New("the record exists")
	ErrMaxRegistered  = errors.New("the maximum number of users has been registered")
)
