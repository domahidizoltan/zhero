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

func MapBy[KEY comparable, ITEM any](items []ITEM, fn func(ITEM) (KEY, ITEM)) iter.Seq2[KEY, ITEM] {
	return func(yield func(KEY, ITEM) bool) {
		for _, item := range items {
			if !yield(fn(item)) {
				break
			}
		}
	}
}

func Unique[T comparable](s []T) iter.Seq[T] {
	seen := make(map[T]struct{})
	return func(yield func(T) bool) {
		for _, v := range s {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				if !yield(v) {
					break
				}
			}
		}
	}
}
