package kismet

func sliceContainsString(slice []string, s string) bool {
	for _, c := range slice {
		if c == s {
			return true
		}
	}
	return false
}
