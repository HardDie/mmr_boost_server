package utils

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
