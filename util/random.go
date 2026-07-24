package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano()) //nolint:staticcheck
}

// RandomInt generates a random integer between min and max (inclusive).
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random lowercase alphabetical string of length n.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates a random username.
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money in paise (1–100000).
func RandomMoney() int64 {
	return RandomInt(100, 100000)
}

// RandomCurrency returns a random supported currency code.
func RandomCurrency() string {
	currencies := []string{INR, USD, EUR}
	return currencies[rand.Intn(len(currencies))]
}

// RandomEmail generates a random email address.
func RandomEmail() string {
	return fmt.Sprintf("%s@ledgr.dev", RandomString(6))
}

// RandomAccountName generates a random account label.
func RandomAccountName() string {
	names := []string{"Cash", "Savings", "HDFC Checking", "SBI Savings", "Wallet", "Emergency Fund"}
	return names[rand.Intn(len(names))]
}

// RandomCategoryName generates a random spending category name.
func RandomCategoryName() string {
	names := []string{"Food", "Transport", "Entertainment", "Health", "Shopping", "Utilities", "Rent", "Salary"}
	return names[rand.Intn(len(names))]
}

// RandomCategoryType returns either "income" or "expense".
func RandomCategoryType() string {
	types := []string{"income", "expense"}
	return types[rand.Intn(len(types))]
}

// RandomFrequency returns a random recurring payment frequency.
func RandomFrequency() string {
	freqs := []string{"daily", "weekly", "monthly", "yearly"}
	return freqs[rand.Intn(len(freqs))]
}

// RandomFutureDate returns a random date between tomorrow and 1 year from now.
func RandomFutureDate() time.Time {
	daysAhead := rand.Intn(365) + 1
	return time.Now().AddDate(0, 0, daysAhead).Truncate(24 * time.Hour)
}
