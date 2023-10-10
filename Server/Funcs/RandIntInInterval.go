package Funcs

import "math/rand"

func RandIntInInterval(min int, max int) int {
	return rand.Intn(max-min) + min
}
