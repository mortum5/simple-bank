package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyz"
	ownerLength = 6
	minMoney    = 0
	maxMoney    = 1000
)

var random *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	random = rand.New(source)
}

func RandomInt(min, max int64) int64 {
	return min + random.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[random.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(ownerLength)
}

func RandomMoney() int64 {
	return RandomInt(minMoney, maxMoney)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)

	return currencies[random.Intn(n)]
}
