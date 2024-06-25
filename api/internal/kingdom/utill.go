package kingdom

func sliceContainsAny[V comparable](slice []V, target ...V) bool {
	for _, v := range slice {
		for _, t := range target {
			if v == t {
				return true
			}
		}
	}
	return false
}
