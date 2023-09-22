package Funcs

import "github.com/gosimple/slug"

func GetSlug(name string) string {
	return slug.Make(name)
}
