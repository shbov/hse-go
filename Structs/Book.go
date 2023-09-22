package Structs

import (
	"hw-1/Funcs"
)

type Book struct {
	Name   string
	Author string
	uid    string
}

func CreateBook(name string, author string) Book {
	return Book{
		Name:   name,
		Author: author,
		uid:    Funcs.GetSlug(name),
	}
}

func (Book *Book) getSlug() string {
	return Funcs.GetSlug(Book.Name)
}
