package services

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrMissingTitle = errors.New("title is required")
)
