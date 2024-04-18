package directmq

import "math/rand"

// Fisher-Yates shuffle algorithm
func randomOrder[T interface{}](items []T) []T {
	reordered := make([]T, len(items))
	copy(reordered, items)

	for i := len(reordered) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		reordered[i], reordered[j] = reordered[j], reordered[i]
	}

	return reordered
}

func contains[T comparable](items []T, item T) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}

	return false
}
