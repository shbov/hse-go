package Errors

import "errors"

var (
	BookNotFound = errors.New("not found")
	BookExists   = errors.New("already exists")
)
