package testbench

func contains[T comparable](items []T, item T) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}

	return false
}
