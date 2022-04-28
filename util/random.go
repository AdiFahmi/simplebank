package util

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt returns a random integer between min and max
func RandomInteger(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

// RandomString returns a random string of length n
func RandomString(n int) string {
	var letterRunes = []rune(alphabet)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInteger(1, 100)
}

func RandomCurrency() string {
	mapCurr := supportedCurrency
	currencies := make([]string, len(mapCurr))
	for k := range mapCurr {
		currencies = append(currencies, k)
	}
	return currencies[rand.Intn(len(currencies))]
}
