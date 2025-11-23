package contextstore

func MapToSlice(m map[string]FileDetails) []FileDetails {
	out := make([]FileDetails, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}
