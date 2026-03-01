// Package _err is for error helper functions
package _err

import "fmt"

func WrapNotNil(resultErr, wrapToErr error) error {
	if resultErr == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", resultErr, wrapToErr)
}
