package Structs

import (
	"hw-1/Generators"
)

type Library struct {
	storage   Storage
	generator Generators.UidGenerator
}

func NewLibrary(storage Storage, generator Generators.UidGenerator) *Library {
	return &Library{storage: storage, generator: generator}
}

func (library *Library) GetBookByName(name string) (*Book, error) {
	return library.storage.GetBookByUid(library.generator.Generate(name))
}

func (library *Library) AddBook(book Book) error {
	book.Uid = library.generator.Generate(book.Name)

	return library.storage.AddBook(book)
}
