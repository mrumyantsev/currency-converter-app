package utils

import "errors"

const (
	ERROR_WORD = ". error: "
)

func DecorateError(desc string, err error) error {
	return errors.New(desc + ERROR_WORD + err.Error())
}
