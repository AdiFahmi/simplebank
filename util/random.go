package util

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lowerAlphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt returns a random integer between min and max
func RandomInteger(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

// RandomString returns a random string of length n
func RandomString(n int, isLower bool) string {
	var letterRunes = []rune(alphabet)
	if isLower {
		letterRunes = []rune(lowerAlphabet)
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6, false)
}

func RandomMoney() int64 {
	return RandomInteger(1, 100)
}

func RandomCurrency() string {
	mapCurr := supportedCurrencyMap
	currencies := make([]string, len(mapCurr))
	i := 0
	for k := range mapCurr {
		currencies[i] = k
		i++
	}
	rCur := currencies[rand.Intn(len(currencies))]
	return rCur
}

func RandomEmail() string {
	return RandomString(6, true) + "@" + RandomString(6, true) + ".com"
}
