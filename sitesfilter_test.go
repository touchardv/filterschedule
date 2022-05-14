package filterschedule

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TesLoadFromFile(t *testing.T) {
	f, err := LoadFromFile("filterschedule.yaml")
	if err != nil {
		t.Fatal("Expected error to be nil but got: ", err)
	}
	if len(f) != 2 {
		t.Fatal("Expected 2 filters but got: ", len(f))
	}
}

func TestPermanentForAllFilter(t *testing.T) {
	sf := SitesFilter{
		Sites: []string{"facebook", "instagram"},
		For:   Targets{All: true},
		When:  Schedule{Always: true},
	}
	timestamp := time.Now()
	ip := "192.168.1.2"

	assertIsMatching(t, sf, "www.facebook.com", ip, timestamp)
	assertIsMatching(t, sf, "www.instagram.com", ip, timestamp)
	assertIsNotMatching(t, sf, "www.wikipedia.com", ip, timestamp)
}

func TestPermanentNotForAllFilter(t *testing.T) {
	sf := SitesFilter{
		Sites: []string{"facebook", "instagram"},
		For:   Targets{All: false, Hosts: []string{"192.168.1.66"}},
		When:  Schedule{Always: true},
	}
	timestamp := time.Now()

	assertIsNotMatching(t, sf, "www.facebook.com", "192.168.1.2", timestamp)
	assertIsNotMatching(t, sf, "www.instagram.com", "192.168.1.2", timestamp)
	assertIsNotMatching(t, sf, "www.wikipedia.com", "192.168.1.2", timestamp)

	assertIsMatching(t, sf, "www.facebook.com", "192.168.1.66", timestamp)
	assertIsMatching(t, sf, "www.instagram.com", "192.168.1.66", timestamp)
	assertIsNotMatching(t, sf, "www.wikipedia.com", "192.168.1.66", timestamp)
}

func TestScheduledDayForAllFilter(t *testing.T) {
	sf := SitesFilter{
		Sites: []string{"facebook", "instagram"},
		For:   Targets{All: true},
		When: Schedule{Always: false, Schedule: []ScheduleEntry{{Days: []DayOfTheWeek{
			DayOfTheWeek(time.Sunday),
		}}}},
	}
	ip := "192.168.1.2"

	sunday := time.Date(2022, 5, 8, 10, 30, 0, 0, time.UTC)
	assertIsMatching(t, sf, "www.facebook.com", ip, sunday)
	assertIsMatching(t, sf, "www.instagram.com", ip, sunday)
	assertIsNotMatching(t, sf, "www.wikipedia.com", ip, sunday)

	monday := time.Date(2022, 5, 9, 10, 30, 0, 0, time.UTC)
	assertIsNotMatching(t, sf, "www.facebook.com", ip, monday)
	assertIsNotMatching(t, sf, "www.instagram.com", ip, monday)
	assertIsNotMatching(t, sf, "www.wikipedia.com", ip, monday)
}

func TestScheduledHourForAllFilter(t *testing.T) {
	sf := SitesFilter{
		Sites: []string{"facebook", "instagram"},
		For:   Targets{All: true},
		When: Schedule{Always: false, Schedule: []ScheduleEntry{{Hours: []TimeInteval{
			{From: "00:00", To: "20:00"},
			{From: "21:00", To: "23:59"},
		}}}},
	}

	ip := "192.168.1.2"

	before8pm := time.Date(2022, 5, 9, 10, 30, 0, 0, time.UTC)
	assertIsMatching(t, sf, "www.facebook.com", ip, before8pm)
	assertIsMatching(t, sf, "www.instagram.com", ip, before8pm)
	assertIsNotMatching(t, sf, "www.wikipedia.com", ip, before8pm)

	eightpm := time.Date(2022, 5, 9, 20, 0, 0, 0, time.UTC)
	assertIsNotMatching(t, sf, "www.facebook.com", ip, eightpm)
	assertIsNotMatching(t, sf, "www.instagram.com", ip, eightpm)
	assertIsNotMatching(t, sf, "www.wikipedia.com", ip, eightpm)

	eigth30pm := time.Date(2022, 5, 9, 20, 30, 0, 0, time.UTC)
	assertIsNotMatching(t, sf, "www.facebook.com", ip, eigth30pm)
	assertIsNotMatching(t, sf, "www.instagram.com", ip, eigth30pm)
	assertIsNotMatching(t, sf, "www.wikipedia.com", ip, eigth30pm)

	ninepm := time.Date(2022, 5, 9, 21, 0, 0, 0, time.UTC)
	assertIsNotMatching(t, sf, "www.facebook.com", ip, ninepm)
	assertIsNotMatching(t, sf, "www.instagram.com", ip, ninepm)
	assertIsNotMatching(t, sf, "www.wikipedia.com", ip, ninepm)

	after9pm := time.Date(2022, 5, 9, 21, 0, 1, 0, time.UTC)
	assertIsMatching(t, sf, "www.facebook.com", ip, after9pm)
	assertIsMatching(t, sf, "www.instagram.com", ip, after9pm)
	assertIsNotMatching(t, sf, "www.wikipedia.com", ip, after9pm)
}

func assertIsMatching(t *testing.T, sf SitesFilter, name string, ip string, timestamp time.Time) {
	if !sf.IsMatching(name, ip, timestamp) {
		t.Log(runtime.Caller(1))
		t.Fatalf("%s should match with name=%s ip=%s timestamp=%s", fmt.Sprint(sf), name, ip, timestamp)
	}
}

func assertIsNotMatching(t *testing.T, sf SitesFilter, name string, ip string, timestamp time.Time) {
	if sf.IsMatching(name, ip, timestamp) {
		t.Log(runtime.Caller(1))
		t.Fatalf("%s should not match with name=%s ip=%s timestamp=%s", fmt.Sprint(sf), name, ip, timestamp)
	}
}
