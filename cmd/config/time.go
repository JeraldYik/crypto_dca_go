package config

import (
	"time"
)

type Time struct {
	now time.Time
}

var t Time

func timeInit(config *Config, now *time.Time) {
	if now != nil {
		t.now = *now
	} else {
		t.now = time.Now()
	}
}

func GetTime() *Time {
	return &t
}

func (t Time) NowTimestamp() int64 {
	return t.now.Unix()
}

func (t Time) GetTodayDate() time.Time {
	return time.Date(t.now.Year(), t.now.Month(), t.now.Day(), 0, 0, 0, 0, time.UTC)
}

func (t Time) GetNowDateString() string {
	return t.now.Format("02/01/2006")
}
