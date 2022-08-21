package lists

func Map[T any, U any](ts []T, m func(t T) U) []U {
	us := make([]U, len(ts))

	for idx, t := range ts {
		us[idx] = m(t)
	}

	return us
}
