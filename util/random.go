package util

import (
	"math/rand"
	"strings"
	"time"
)

// func ini akan otomatis called ketika package util pertama kali used
func init() {
	// disini kita akan set nilai seed dg random generator caranya panggil rand.Seed()
	// normalnya Seed() will return current time time.Now() kita convert jd unix nano UnixNano()
	// jd ketika run akan generate nilai yg berbeda"
	rand.Seed(time.Now().UnixNano())
}

// this func will generate random int
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// deklarasi alpabetnya dulu untuk generate string, total karakter 26
const alphabet = "abcdefghijklmnopqrstuvwxyz"

// this func will generate random string dari n karakter
func RandomString(n int) string {
	// deklarasi objek string builder baru
	var sb strings.Builder
	k := len(alphabet)

	// generate karakter random
	for i := 0; i < n; i++ {
		// generate random position dr 1 sampai k-1
		c := alphabet[rand.Intn(k)]
		// string builder
		sb.WriteByte(c)
	}
	return sb.String()
}

// generate random owner
func RandomOwner() string {
	// return 6 huruf string
	return RandomString(6)
}

// generate random amount of money
func RandomMoney() int64 {
	// isinya random int dr 0 - 1000
	return RandomInt(0, 1000)
}

// generate random currency
func RandomCurrency() string {
	// list berisi 3 buah currencies, EUR, USD, CAD
	currencies := []string{"EUR", "USD", "CAD"}
	// hitung length dr list currency
	n := len(currencies)
	// kembalikan random
	return currencies[rand.Intn(n)]
}
