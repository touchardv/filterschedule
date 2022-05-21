package filterschedule

import (
	"encoding/json"
	"errors"
	"time"
)

type DayOfTheWeek time.Weekday

var daysOfWeek = map[string]DayOfTheWeek{}

func init() {
	for d := time.Sunday; d <= time.Saturday; d++ {
		name := d.String()
		daysOfWeek[name] = DayOfTheWeek(d)
	}
}

func (i *DayOfTheWeek) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if d, ok := daysOfWeek[s]; ok {
		*i = d
		return nil
	}
	return errors.New("invalid week day name: " + s)
}

type Schedule struct {
	Always   bool
	Schedule []ScheduleEntry
}

type ScheduleEntry struct {
	Days  []DayOfTheWeek
	Hours []TimeInteval
}

type TimeInteval struct {
	From string
	To   string
}

func (ti *TimeInteval) includes(t time.Time) bool {
	date := t.Format("2006-01-02")
	from, _ := time.ParseInLocation("2006-01-02T15:04", date+"T"+ti.From, t.Location())
	to, _ := time.ParseInLocation("2006-01-02T15:04", date+"T"+ti.To, t.Location())

	return from.Before(t) && t.Before(to)
}

func isApplicableAtDays(days []DayOfTheWeek, t time.Time) bool {
	if len(days) == 0 {
		return true
	} else {
		for _, d := range days {
			if DayOfTheWeek(t.Weekday()) == d {
				log.Debug("Matched: ", days)
				return true
			}
		}
	}
	return false
}

func isApplicableAtHours(timeIntervals []TimeInteval, t time.Time) bool {
	if len(timeIntervals) == 0 {
		return true
	} else {
		for _, ti := range timeIntervals {
			if ti.includes(t) {
				log.Debug("Matched: ", timeIntervals)
				return true
			}
		}
	}
	return false
}
