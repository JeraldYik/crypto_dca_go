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
			},
			want: time.Date(2024, time.November, 3, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Time{
				now: tt.fields.now,
			}
			got := tr.GetTodayDate()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTime_GetNowDateString(t *testing.T) {
	type fields struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				now: TestNow,
			},
			want: TestNowDateStr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Time{
				now: tt.fields.now,
			}
			got := tr.GetNowDateString()
			assert.Equal(t, tt.want, got)
		})
	}
}
