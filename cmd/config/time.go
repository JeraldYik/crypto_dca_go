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

// Keeps changing, to satisfy uniqueness in gemini http requests
func (t Time) NowTimestamp() int64 {
	return time.Now().UnixNano()
}

func (t Time) GetTodayDate() time.Time {
	return time.Date(t.now.Year(), t.now.Month(), t.now.Day(), 0, 0, 0, 0, time.UTC)
}

func (t Time) GetNowDateString() string {
	return t.now.Format("02/01/2006")
}

func (t Time) Now() time.Time {
	return t.now
}
