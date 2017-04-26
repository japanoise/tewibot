package main

import (
	"math/rand"
)

func randoms(fr []string) string {
	return fr[rand.Intn(len(fr))]
}
