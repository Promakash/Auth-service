package time

import "time"

type UnixTime = int64

func DaysToUnix(days time.Duration) UnixTime {
	return UnixTime(days * 86400)
}
