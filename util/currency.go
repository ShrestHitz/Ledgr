package util

// Supported currency codes.
const (
	INR = "INR"
	USD = "USD"
	EUR = "EUR"
)

// IsSupportedCurrency returns true if the currency code is supported.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case INR, USD, EUR:
		return true
	}
	return false
}
