package util

// Slugify converts a string to a kebab-case slug suitable for branch names and paths.
func Slugify(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result = append(result, c-'A'+'a')
		} else if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '_' {
			result = append(result, c)
		} else if c == ' ' || c == '/' || c == ':' {
			if len(result) > 0 && result[len(result)-1] != '-' {
				result = append(result, '-')
			}
		}
	}
	// Trim trailing dash
	for len(result) > 0 && result[len(result)-1] == '-' {
		result = result[:len(result)-1]
	}
	if len(result) == 0 {
		return "unnamed"
	}
	return string(result)
}
