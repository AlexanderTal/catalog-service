package entity

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrCategoryHasProducts = errors.New("category has linked products")
)
