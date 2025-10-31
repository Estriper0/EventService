package service

import "errors"

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrRepositoryError = errors.New("error in the repository")
)
