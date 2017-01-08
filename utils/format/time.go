package format

import (
	"fmt"
	"time"
)

func GetUnixTimeFromDBTimestamp(timestamp string) string {
	var formattedTime string

	t, err := time.Parse(time.RFC3339, timestamp)
	if err == nil {
		formattedTime = t.Format(time.UnixDate)
	}

	return formattedTime
}

func GetPrettyTimeFromDBTimestamp(timestamp string) string {
	var formattedTime string

	t, err := time.Parse(time.RFC3339, timestamp)
	if err == nil {
		suffix := "AM"
		if t.Hour() > 12 {
			suffix = "PM"
		}
		formattedTime = fmt.Sprintf("%02d-%02d-%04d %d:%02d:%02d %s",
			t.Month(), t.Day(), t.Year(), t.Hour()%12, t.Minute(), t.Second(), suffix)
	}

	return formattedTime
}
