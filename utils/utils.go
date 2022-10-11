package utils

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func BoolPtr(b bool) *bool { return &b }
func IntPtr(i int) *int { return &i }
func StringPtr(s string) *string { return &s }

func UseStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
