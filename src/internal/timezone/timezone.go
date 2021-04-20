package timezone

import (
	"time"
)

const (
	BasicFormat = "2006-01-02 3:04:05PM"
	PST         = "America/Los_Angeles"
)

func GetTime(timezone string) (time.Time, string) {
	now := time.Now()
	loc, _ := time.LoadLocation(PST)
	currentTime := now.In(loc)
	return currentTime, currentTime.Format(BasicFormat)
}

func TimeParse(currentTime string) (time.Time, error) {
	t, err := time.Parse(BasicFormat, currentTime)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}
