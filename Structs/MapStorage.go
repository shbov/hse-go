package Structs

import (
	"fmt"
	"hw-1/Errors"
)

// unique id for book = slug

type MapStorage struct {
	books map[string]Book
}

func NewMapStorage() *MapStorage {
	return &MapStorage{books: make(map[string]Book)}
}

func (mapStorage *MapStorage) GetBookByUid(uid string) (*Book, error) {
	book, ok := mapStorage.books[uid]
	if ok {
		return &book, nil
	}

	return nil, fmt.Errorf("[library:error] book with uid='%s': %w", uid, Errors.BookNotFound)
}

func (mapStorage *MapStorage) AddBook(book Book) error {
	_, ok := mapStorage.books[book.Uid]
	if ok {
		return fmt.Errorf("[library:error] book with uid='%s': %w", book.Uid, Errors.BookExists)
	}

	mapStorage.books[book.Uid] = book
	return nil
}
