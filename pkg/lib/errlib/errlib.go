package errlib

import "fmt"

const (
	tmplWrap = "%s: %w"
)

func Wrap(err error, msg string) error {
	return fmt.Errorf(tmplWrap, msg, err)
}
