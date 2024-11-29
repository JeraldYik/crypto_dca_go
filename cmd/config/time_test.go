package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime_GetTodayDate(t *testing.T) {
	TestInit(nil, &TestNow)
	type fields struct {
		now time.Time
		loc *time.Location
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "ok",
			fields: fields{
				now: TestNow,
				loc: time.UTC,
			},
			want: time.Date(2024, time.November, 3, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Time{
				now: tt.fields.now,
				loc: tt.fields.loc,
			}
			got := tr.GetTodayDate()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTime_GetNowDateString(t *testing.T) {
	type fields struct {
		now time.Time
		loc *time.Location
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "UTC",
			fields: fields{
				now: TestNow,
				loc: time.UTC,
			},
			want: TestNowDateStr,
		},
		{
			name: "PST",
			fields: fields{
				now: time.Date(2024, time.November, 3, 0, 0, 0, 0, time.FixedZone("PST", -8*60*60)),
				loc: time.FixedZone("PST", -8*60*60),
			},
			want: TestNowDateStr,
		},
		{
			name: "GMT+8",
			fields: fields{
				now: time.Date(2024, time.November, 3, 0, 0, 0, 0, time.FixedZone("GMT+8", 8*60*60)),
				loc: time.FixedZone("GMT+8", 8*60*60),
			},
			want: TestNowDateStr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Time{
				now: tt.fields.now,
				loc: tt.fields.loc,
			}
			got := tr.GetNowDateString()
			assert.Equal(t, tt.want, got)
		})
	}
}
