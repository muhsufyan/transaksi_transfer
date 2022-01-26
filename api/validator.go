package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/muhsufyan/transaksi_transfer/util"
)

/*
saat melakukan transaksi currency we set 3 buah, bagaimana jika ada 100 ? tdk mungkin dg oneof solusinya dg custom validator sprti ini
*/

var validCurrency validator.Func = func(FieldLevel validator.FieldLevel) bool {
	// get nilai dr field, is reflection, convert to string. jika ok maka currency valid
	if currency, ok := FieldLevel.Field().Interface().(string); ok {
		// cek currency is supported
		return util.IsSupportedCurrency(currency)
	}
	// currency not string
	return false
}
