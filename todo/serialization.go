package todo

import (
	"encoding/json"
	"time"
)

type MyTime time.Time

func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).Format(`"02.01.2006 15:04:05"`)
	return []byte(formatted), nil
}

func (t *MyTime) UnmarshalJSON(data []byte) error {
	const layout = `"02.01.2006 15:04:05"`
	parsed, err := time.ParseInLocation(layout, string(data), time.Local)
	if err != nil {
		return err
	}
	*t = MyTime(parsed)
	return nil
}

type DurationString time.Duration

func (d DurationString) MarshalJSON() ([]byte, error) {
	s := time.Duration(d).String() // "1h2m3s"
	return json.Marshal(s)
}

func (d *DurationString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = DurationString(duration)
	return nil
}
