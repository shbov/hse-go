package Structs

type Storage interface {
	GetBookByUid(uid string) (*Book, error)
	AddBook(book Book) error
}
