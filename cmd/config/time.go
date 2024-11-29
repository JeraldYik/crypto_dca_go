package config

import (
	"log"
	"time"
)

type Time struct {
	now time.Time
	loc *time.Location
}

var t Time

func timeInit(config *Config, now *time.Time) {
	location, err := time.LoadLocation(config.Location)
	if err != nil {
		log.Panicf("[config.timeInit] Unable to load location, err: %+v", err)
	}
	t.loc = location
	if now != nil {
		t.now = *now
	} else {
		t.now = time.Now().In(t.loc)
	}
}

func GetTime() *Time {
	return &t
}

func (t Time) NowTimestamp() int64 {
	return t.now.In(t.loc).Unix()
}

func (t Time) GetTodayDate() time.Time {
	return time.Date(t.now.Year(), t.now.Month(), t.now.Day(), 0, 0, 0, 0, t.loc)
}

func (t Time) GetNowDateString() string {
	return t.now.Format("02/01/2006")
}
