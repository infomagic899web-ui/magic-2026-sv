package timein

import "time"

func Years(n int) time.Duration {
	return time.Duration(n) * 365 * 24 * time.Hour
}
