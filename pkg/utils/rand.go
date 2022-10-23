package utils

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandomInt(min int, max int) int {
	if max-min == 0 {
		return min
	}
	return rand.Intn(max-min) + min
}

func RandomFloat(min int, max int) float64 {
	return rand.Float64() + float64(RandomInt(min, max))
}
