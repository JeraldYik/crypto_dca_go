package config

import (
	"time"
)

type Time struct {
	now time.Time
}

var t Time

func timeInit(now *time.Time) {
	if now != nil {
		t.now = *now
	} else {
		t.now = time.Now()
	}
}

func GetTime() *Time {
	return &t
}

func (t Time) NowTimestamp(now *time.Time) int64 {
	timeInit(now)
	return t.now.UnixNano()
}

func (t Time) GetTodayDate() time.Time {
	return time.Date(t.now.Year(), t.now.Month(), t.now.Day(), 0, 0, 0, 0, time.UTC)
}

func (t Time) GetNowDateString() string {
	return t.now.Format("02/01/2006")
}
