package util

var supportedCurrency = map[string]bool{
	"USD": true,
	"IDR": true,
	"EUR": true,
}

func IsSupportedCurrency(currency string) bool {
	_, ok := supportedCurrency[currency]
	return ok
}
