package service

import "errors"

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrRepositoryError = errors.New("error in the repository")
	ErrRegistered      = errors.New("the user is already registered")
	ErrNotRegistered   = errors.New("the user is not registered")
	ErrMaxRegistered   = errors.New("the maximum number of users has been registered")
)
