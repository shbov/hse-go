package Generators

import "github.com/gosimple/slug"

type UidGenerator interface {
	Generate(name string) string
}

type SlugGenerator struct {
}

func NewSlugGenerator() *SlugGenerator {
	return &SlugGenerator{}
}

func (SlugGenerator *SlugGenerator) Generate(name string) string {
	return slug.Make(name)
}
