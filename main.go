package main

import (
	"fmt"
	"hw-1/Generators"
	"hw-1/Structs"
)

func testLibrary(library *Structs.Library, books []Structs.Book) {
	fmt.Printf("Testing library %T\n", library)

	// adding books to the library
	_ = library.AddBook(books[0])
	_ = library.AddBook(books[1])
	_ = library.AddBook(books[2])

	// try to add book with existing slug
	err := library.AddBook(books[2])
	if err != nil {
		fmt.Println(err)
	}

	// try to find book by name
	book, err := library.GetBookByName("Книга 1")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("The book was found: %s\n", book)
	}

}

func main() {
	// create library
	library := Structs.NewLibrary(Structs.NewMapStorage(), Generators.NewSlugGenerator())

	// create slice of books
	var books []Structs.Book
	books = append(books, Structs.CreateDefaultBook(
		"Книга 1",
		"Автор 1",
	))
	books = append(books, Structs.CreateDefaultBook(
		"Книга 2",
		"Автор 1",
	))
	books = append(books, Structs.CreateDefaultBook(
		"Книга 3",
		"Автор 1",
	))

	// test library with map-based storage system
	testLibrary(library, books)

	fmt.Println()

	// test library with slice-based storage system
	library = Structs.NewLibrary(Structs.NewSliceStorage(), Generators.NewSlugGenerator())
	testLibrary(library, books)
}
