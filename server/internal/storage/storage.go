package storage

import "errors"

var (
	ErrNotFound  = errors.New("thumbnail not found")
	ErrUrlExists = errors.New("thumbnail url already exists")
)
