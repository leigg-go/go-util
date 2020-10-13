package util

import (
	"github.com/leigg-go/go-util/_util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GetEachDayZeroClockFromStToEt(t *testing.T) {
	oneDayZero := time.Date(2020, 10, 1, 0, 0, 0, 0, time.Local)
	d1031 := time.Date(2020, 10, 31, 0, 0, 0, 0, time.Local)
	d1101 := time.Date(2020, 11, 1, 0, 0, 0, 0, time.Local)
	d1102 := time.Date(2020, 11, 2, 0, 0, 0, 0, time.Local)
	test := []struct {
		name       string
		start, end time.Time
		want       []time.Time
	}{
		{
			name:  "Wrong time args test: end-time < start-time",
			start: oneDayZero.Add(time.Duration(1)),
			end:   oneDayZero,
			want:  []time.Time{},
		},
		{
			name:  "Same day test",
			start: oneDayZero.Add(time.Duration(1)),
			end:   oneDayZero.Add(time.Duration(2)),
			want:  []time.Time{oneDayZero},
		},
		{
			name:  "Exactly same day test",
			start: oneDayZero,
			end:   oneDayZero,
			want:  []time.Time{oneDayZero},
		},
		{
			name:  "Cross month day test",
			start: d1031,
			end:   d1102,
			want:  []time.Time{d1031, d1101, d1102},
		},
	}

	for _, tt := range test {
		got := _util.GetEachDayZeroClockFromStToEt(tt.start, tt.end)
		assert.Equal(t, tt.want, got, "[%s]: got:%+v not equal want:%+v", tt.name, got, tt.want)
	}
}
