package api

import (
	"github.com/dxtym/bankrupt/utils"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check if valid currency
		return utils.CurrencySupported(currency)
	}

	return false
}
