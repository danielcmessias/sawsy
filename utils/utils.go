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

func BoolPtr(b bool) *bool          { return &b }
func IntPtr(i int) *int             { return &i }
func StringPtr(s string) *string    { return &s }
func Float64Ptr(f float64) *float64 { return &f }

func UseStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Compose2[A0 any, T1 any, T2 any](fn0 func(A0) T1, fn1 func(T1) T2) func(A0) T2 {
	return func(data A0) T2 {
		return fn1(fn0(data))
	}
}

func BytesToMB(f float64) float64 {
	return f / (1024 * 1024)
}
