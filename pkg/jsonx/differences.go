package jsonx

type Differences []*Difference

func (diffs Differences) addPath(path string) Differences {
	for _, diff := range diffs {
		diff.Path = append(diff.Path, path)
	}
	return diffs
}

func (diffs Differences) String() string {
	var result string
	for _, diff := range diffs {
		result += diff.String()
	}
	return result
}
