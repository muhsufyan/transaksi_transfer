package util

/*
implementasi logik to cek if a currency is supported or not in this file
*/
// semua currency yg didukung/dpt digunakan
const (
	USD    = "USD"
	EUR    = "EUR"
	CAD    = "CAD"
	RP     = "RP"
	DINAR  = "DINAR"
	DIRHAM = "DIRHAM"
)

// true if currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, RP, DINAR, DIRHAM:
		return true
	}
	return false
}
