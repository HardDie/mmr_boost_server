package utils

import "time"

func Allocate[T any](val T) *T {
	return &val
}

func ToInt32[T ~int32](val *T) *int32 {
	if val == nil {
		return nil
	}
	res := int32(*val)
	return &res
}

func Compare[T comparable](a, b *T) bool {
	switch {
	case a == nil && b == nil:
	case a != nil && b != nil:
		if *a != *b {
			return false
		}
	default:
		return false
	}
	return true
}
func CompareTime(a, b *time.Time) bool {
	switch {
	case a == nil && b == nil:
	case a != nil && b != nil:
		if !a.Equal(*b) {
			return false
		}
	default:
		return false
	}
	return true
}
