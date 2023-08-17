package api

import (
	"github.com/AdamDomagalsky/goes/2023/bank/util"

	"github.com/go-playground/validator/v10"
)

var validatorCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
