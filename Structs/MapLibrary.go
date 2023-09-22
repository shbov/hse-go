package Structs

import (
	"fmt"
	"hw-1/Funcs"
)

// unique id for book = slug

type MapLibrary struct {
	books map[string]Book
}

func (MapLibrary *MapLibrary) GetBookByName(name string) (*Book, error) {
	uid := Funcs.GetSlug(name)
	book, ok := MapLibrary.books[uid]

	if ok {
		return &book, nil
	}

	return nil, fmt.Errorf("[library:error] book with name='%s' wasn't found", name)
}

func CreateMapLibrary() *MapLibrary {
	return &MapLibrary{
		books: make(map[string]Book),
	}
}

func (MapLibrary *MapLibrary) AddBook(book Book) error {
	_, ok := MapLibrary.books[book.uid]

	if ok {
		return fmt.Errorf("[library:error] book with uid='%s' already exists", book.uid)
	}

	MapLibrary.books[book.uid] = book
	return nil
}
