package util

const (
    CNY = "CNY"
    USD = "USD"
    EUR = "EUR"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
    switch currency {
    case CNY, USD, EUR:
        return true
    }
    return false
}