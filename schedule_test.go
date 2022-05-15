package filterschedule

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDayOfTheWeekJSON(t *testing.T) {
	v := `["Monday", "Sunday"]`
	var d []DayOfTheWeek
	err := json.Unmarshal([]byte(v), &d)
	if err != nil {
		t.Fatal("JSON deserialization should return no error:", err)
	}
	if len(d) != 2 {
		t.Fatal("JSON deserialization should return 2 entries")
	}
	if d[0] != DayOfTheWeek(time.Monday) {
		t.Fatal("JSON deserialization should return time.Monday value")
	}
}

func TestTimeInterval(t *testing.T) {
	interval := TimeInteval{From: "12:00", To: "14:00"}

	timestamp := time.Date(2022, 5, 9, 12, 0, 0, 0, time.UTC)
	if interval.includes(timestamp) {
		t.Fatal(interval, " should not include ", timestamp)
	}

	timestamp = time.Date(2022, 5, 9, 12, 0, 1, 0, time.UTC)
	if !interval.includes(timestamp) {
		t.Fatal(interval, " should include ", timestamp)
	}

	timestamp = time.Date(2022, 5, 9, 13, 59, 0, 0, time.UTC)
	if !interval.includes(timestamp) {
		t.Fatal(interval, " should include ", timestamp)
	}

	timestamp = time.Date(2022, 5, 9, 14, 0, 0, 0, time.UTC)
	if interval.includes(timestamp) {
		t.Fatal(interval, " should not include ", timestamp)
	}

	location, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		panic(err)
	}
	timestamp = time.Date(2022, 5, 9, 13, 59, 0, 0, location)
	if !interval.includes(timestamp) {
		t.Fatal(interval, " should include ", timestamp)
	}
}
