package Structs

import (
	"fmt"
)

// unique id for book = index

type SliceLibrary struct {
	books []Book
}

func (SliceLibrary *SliceLibrary) GetBookByName(name string) (*Book, error) {
	for _, book := range SliceLibrary.books {
		if book.Name == name {
			return &book, nil
		}
	}

	return nil, fmt.Errorf("[library:error] book with name='%s' wasn't found", name)
}

func CreateSliceLibrary() *SliceLibrary {
	return &SliceLibrary{
		books: make([]Book, 0),
	}
}

func (SliceLibrary *SliceLibrary) AddBook(book Book) error {
	_, err := SliceLibrary.GetBookByName(book.Name)
	if err == nil {
		return fmt.Errorf("[library:error] book with name='%s' already exists", book.Name)
	}

	SliceLibrary.books = append(SliceLibrary.books, book)
	return nil
}
