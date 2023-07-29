package common

import (
	"time"
)

const (
	layout = "2006-01-02 15:04:05"
)

func CompareTime(t1, t2 string) (string, error) {
	time1, err1 := time.Parse(layout, t1)
	if err1 != nil {
		return "", err1
	}
	time2, err2 := time.Parse(layout, t2)
	if err2 != nil {
		return "", err2
	}

	if time1.Equal(time2) {
		return "=", nil
	} else if time1.Before(time2) {
		return "<", nil
	} else if time1.After(time2) {
		return ">", nil
	}

	return "", nil
}
