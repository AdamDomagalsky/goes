package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	PLN = "PLN"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, PLN:
		return true
	}
	return false
}
