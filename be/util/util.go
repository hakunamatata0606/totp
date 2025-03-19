package util

import "time"

func RoundTimeUTC(t time.Time, interval time.Duration) time.Time {
	return t.UTC().Truncate(interval)
}
