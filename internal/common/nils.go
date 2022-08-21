package common

func OrDefault[T any](value *T, defaultValue T) T {
	if value != nil {
		return *value
	}
	return defaultValue
}
