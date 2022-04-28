package api

import (
	"github.com/adifahmi/simplebank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldlevel validator.FieldLevel) bool {
	if currency, ok := fieldlevel.Field().Interface().(string); ok {
		// check currency
		util.IsSupportedCurrency(currency)
	}
	return false
}
