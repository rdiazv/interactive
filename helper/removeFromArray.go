package helper

func RemoveFromArray(array []interface{}, value interface{}) []interface{} {
	for i, item := range array {
		if item == value {
			return append(array[:i], array[i+1:]...)
		}
	}

	return array
}
