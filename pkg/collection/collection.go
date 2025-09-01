// Package collection is a set of collection helper functions
package collection

import "iter"

func MapValues[IN, OUT any](s []IN, fn func(IN) OUT) iter.Seq[OUT] {
	return func(yield func(OUT) bool) {
		for _, v := range s {
			if !yield(fn(v)) {
				break
			}
		}
	}
}

func FilterValues[T any](s []T, fn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			if fn(v) {
				if !yield(v) {
					break
				}
			}
		}
	}
}
