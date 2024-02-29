package e

import "fmt"

func Wrap(desc string, err error) error {
	return fmt.Errorf("%s: %w", desc, err)
}
