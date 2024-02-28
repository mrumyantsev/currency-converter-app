package utils

import "fmt"

func DecorateError(desc string, err error) error {
	return fmt.Errorf("%s: %w", desc, err)
}
