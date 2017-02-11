package helper

func InArray(array []interface{}, value interface{}) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}
