package middleware

func matchRoute(path string, rules map[string]string) string {
	longest := ""
	for prefix := range rules {
		if len(prefix) > len(longest) && len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			longest = prefix
		}
	}
	return longest
}
