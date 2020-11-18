package _util

import "time"

// Get days start time via start-time to end-time args
// e.g. input: st="2020-09-01 xx:xx:xx"  et="2020-09-03 xx:xx:xx"
// return []{"2020-09-01 00:00:00", "2020-09-02 00:00:00", "2020-09-03 00:00:00"}
func GetEachDayZeroClockFromStToEt(st, et time.Time) []time.Time {
	var ret = make([]time.Time, 0)

	if et.Sub(st) < 0 {
		return ret
	}

	sameDay := st.Year() == et.Year() && st.Month() == et.Month() && st.Day() == et.Day()
	if sameDay {
		ret = append(ret, time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location()))
		return ret
	}

	// reset them to zero o'clock
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	for {
		ret = append(ret, time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location()))
		st = st.Add(time.Hour * 24)
		if et.Sub(st) < 0 {
			break
		}
	}
	return ret
}

func NextMonthFirstDay(t time.Time) time.Time {
	t = t.AddDate(0, 1, 0)
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// 返回给定起始结束时间范围内的月份时间，t为每月1日0点0分...
func GetMonthList(startTime, endTime time.Time) []time.Time {
	var ret []time.Time
	for endTime.Sub(startTime) >= 0 {
		ret = append(ret, time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, startTime.Location()))
		startTime = NextMonthFirstDay(startTime)
	}
	return ret
}
