package utils

func Stringify(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func Intify(v any) int {
	if f, ok := v.(float64); ok {
		return int(f)
	}
	return 0
}
