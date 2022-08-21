package model

type AnyValue = any

type AnyStr struct {
	value string
}

func (a AnyStr) ToString() string {
	return a.value
}
