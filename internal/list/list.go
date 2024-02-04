package list

// ExistInSlice checks whether a comparable element exists in a slice of the same type.
func ExistInSlice[T comparable](item T, list []T) bool {
	if len(list) == 0 {
		return false
	}
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
