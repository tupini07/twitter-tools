package data_utils

import "math/rand"

func ShuffleArrayInplace(arr ...interface{}) {
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}
