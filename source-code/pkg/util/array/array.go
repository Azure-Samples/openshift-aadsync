package array

// ContainsString determines whether the specified array contains the specified value
func ContainsString(arr []string, value string) bool {

	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}
