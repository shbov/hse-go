package Structs

import (
	"fmt"
	"hw-1/Errors"
)

// unique id for book = index

type SliceStorage struct {
	books []Book
}

func NewSliceStorage() *SliceStorage {
	return &SliceStorage{
		books: make([]Book, 0),
	}
}

func (sliceStorage *SliceStorage) GetBookByUid(uid string) (*Book, error) {
	for _, book := range sliceStorage.books {
		if book.Uid == uid {
			return &book, nil
		}
	}

	return nil, fmt.Errorf("[library:error] book with uid='%s': %w", uid, Errors.BookNotFound)
}

func (sliceStorage *SliceStorage) AddBook(book Book) error {
	_, err := sliceStorage.GetBookByUid(book.Uid)
	if err == nil {
		return fmt.Errorf("[library:error] book with uid='%s': %w", book.Uid, Errors.BookExists)
	}

	sliceStorage.books = append(sliceStorage.books, book)
	return nil
}
