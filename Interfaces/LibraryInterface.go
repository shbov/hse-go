package Interfaces

import "hw-1/Structs"

type LibraryInterface interface {
	GetBookByName(name string) (*Structs.Book, error)
	AddBook(book Structs.Book) error
}
