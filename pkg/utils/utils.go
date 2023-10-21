package utils

func ArrayContains(array []string, element string) bool {
	for _, el := range array {
		if el == element {
			return true
		}
	}
	return false
}
