package util

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
    return min + rand.Int63n(max - min + 1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
    var sb strings.Builder
    k := len(alphabet)

    for i := 0; i < n; i++ {
        c := alphabet[rand.Intn(k)]
        sb.WriteByte(c)
    }

    return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
    return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() float64 {
    return roundFloat64(rand.Float64() * float64(RandomInt(0, 10000)), 2)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
    currencies := []string{CNY, USD, EUR}
    n := len(currencies)
    return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email
func RandomEmail() string {
    return fmt.Sprintf("%s@email.com", RandomString(6))
}

// roundFloat64 rounds a float64 number to a decimal place
func roundFloat64(number float64, decimalPlace int) float64 {
    if decimalPlace < 0 || decimalPlace > 20 {
        log.Fatal("decimalPlace is out of range (0-20)")
    }
    // Calculate 10 to the power of decimalPlace
    temp := math.Pow10(decimalPlace)
    // Multiply the float number with 10**decimalPlace and round it
    return math.Round(number * temp) / temp
}