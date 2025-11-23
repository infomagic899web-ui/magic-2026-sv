package utils

import "time"

var LocationAsiaManila, _ = time.LoadLocation("Asia/Manila")

// Check if voting is open (Monday - Friday UTC+8)
func IsVotingOpen() bool {
	now := time.Now().In(LocationAsiaManila)
	weekday := now.Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
}

// Check if today is Monday morning (for resetting upcoming votes)
func IsMondayMorning() bool {
	now := time.Now().In(LocationAsiaManila)
	return now.Weekday() == time.Monday && now.Hour() < 12 // before noon
}
