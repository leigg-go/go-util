package _util

import (
	"github.com/stretchr/testify/assert"
	"reflect"
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
		got := GetEachDayZeroClockFromStToEt(tt.start, tt.end)
		assert.Equal(t, tt.want, got, "[%s]: got:%+v not equal want:%+v", tt.name, got, tt.want)
	}
}

const (
	StdDateTimeLayout = "2006-01-02 15:04:05"
)

func TestGetMonthList(t *testing.T) {
	t1, _ := time.Parse(StdDateTimeLayout, "2020-11-01 00:00:00")
	t2, _ := time.Parse(StdDateTimeLayout, "2020-11-01 01:00:00")
	args1 := []time.Time{t1, t2}
	want1 := []string{"2020-11-01 00:00:00"}

	t3, _ := time.Parse(StdDateTimeLayout, "2020-11-30 00:00:00")
	t4, _ := time.Parse(StdDateTimeLayout, "2020-12-01 00:00:00")
	args2 := []time.Time{t3, t4}
	want2 := []string{"2020-11-01 00:00:00", "2020-12-01 00:00:00"}

	t5, _ := time.Parse(StdDateTimeLayout, "2020-11-01 00:00:00")
	t6, _ := time.Parse(StdDateTimeLayout, "2021-01-02 00:00:00")
	args3 := []time.Time{t5, t6}
	want3 := []string{"2020-11-01 00:00:00", "2020-12-01 00:00:00", "2021-01-01 00:00:00"}

	test := []struct {
		name string
		args []time.Time
		want []string
	}{
		{
			name: "2020-11-01 ~ 2020-11-01",
			args: args1,
			want: want1,
		},
		{
			name: "2020-11-30 ~ 2020-12-01",
			args: args2,
			want: want2,
		},
		{
			name: "2020-11-01 ~ 2021-01-02",
			args: args3,
			want: want3,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ret := GetMonthList(tt.args[0], tt.args[1])
			var comp []string
			for _, t := range ret {
				comp = append(comp, t.Format(StdDateTimeLayout))
			}
			if !reflect.DeepEqual(comp, tt.want) {
				t.Errorf("got:%+v not want:%+v", comp, tt.want)
			}
		})
	}
}
