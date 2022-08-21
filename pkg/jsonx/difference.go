package jsonx

import (
	"fmt"
	"strings"
)

type Difference struct {
	Path        []string
	ExpectedRaw JsonX
	Expected    JsonX
	Actual      JsonX
	Message     string
}

func (d Difference) String() string {
	return fmt.Sprintf("%s: %s.\n  Expected: %s\n  Actual: %s\n", strings.Join(reverse(d.Path), "."), d.Message, d.Expected, d.Actual)
}

func reverse(ds []string) []string {
	l := len(ds)
	if l == 0 {
		return ds
	}

	result := make([]string, l)
	maxI := l - 1
	for i := maxI; i >= 0; i-- {
		result[maxI-i] = ds[i]
	}
	return result
}
