package Structs

type Book struct {
	Name   string
	Author string
	Uid    string
}

func CreateDefaultBook(name string, author string) Book {
	return Book{
		Name:   name,
		Author: author,
		Uid:    "",
	}
}
