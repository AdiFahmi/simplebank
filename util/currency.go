package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	IDR = "IDR"
)

var supportedCurrencyMap = map[string]bool{
	USD: true,
	IDR: true,
	EUR: true,
}

func IsSupportedCurrency(currency string) bool {
	_, ok := supportedCurrencyMap[currency]
	return ok
}
