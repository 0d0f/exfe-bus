package model

import (
	"fmt"
	"time"
)

var nowFunc = func() time.Time {
	return time.Now()
}

type EFTime struct {
	DateWord string `json:"date_word,omitempty"`
	Date     string `json:"date,omitempty"`
	TimeWord string `json:"time_word,omitempty"`
	Time     string `json:"time,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

func (t EFTime) UTCTime(layout string) (string, error) {
	if t.Time == "" && t.Date == "" {
		return "", nil
	}

	var time_ time.Time
	var err error
	if t.Time == "" {
		time_, err = time.Parse("2006-1-02", t.Date)
	} else if t.Date != "" {
		str := fmt.Sprintf("%s %s", t.Date, t.Time)
		time_, err = time.Parse("2006-1-02 15:04:05", str)
	}
	if err != nil {
		return "", err
	}
	return time_.Format(layout), nil
}

func (t EFTime) StringInZone(targetZone string, targetLoc *time.Location) (string, error) {
	t_, err := t.timeInZone(targetLoc)
	if err != nil {
		return "", err
	}

	ret := ""

	if t.TimeWord != "" {
		ret += t.TimeWord
	}

	if t.Time != "" {
		if ret != "" {
			ret += " at "
		}
		ret += t_.Format("3:04PM")
	}

	if loc, _ := LoadLocation(t.Timezone); loc.String() != targetLoc.String() && ret != "" {
		ret += " "
		if t.Time != "" {
			ret += targetZone
		} else {
			ret += t.Timezone
		}
	}

	if t.DateWord != "" {
		if ret != "" {
			ret += " "
		}
		ret += t.DateWord
	}

	if t.Date != "" {
		if ret != "" {
			ret += " on "
		}
		now := nowFunc()
		if now.Year() == t_.Year() {
			ret += t_.Format("Mon, Jan 2")
		} else {
			ret += t_.Format("Mon, Jan 2 2006")
		}
	}

	return ret, nil
}

func (t EFTime) timeInZone(targetLoc *time.Location) (time.Time, error) {
	var t_ time.Time
	var err error

	switch {
	case t.Time != "" && t.Date != "":
		t_, err = time.Parse("2006-1-02 15:04:05", fmt.Sprintf("%s %s", t.Date, t.Time[:8]))
	case t.Time != "" && t.Date == "":
		t_, err = time.Parse("15:04:05", fmt.Sprintf("%s", t.Time[:8]))
	case t.Time == "" && t.Date != "":
		t_, err = time.Parse("2006-1-02", fmt.Sprintf("%s", t.Date))
	}

	if err != nil {
		return t_, fmt.Errorf("Parse time error: %s", err)
	}

	loc, _ := LoadLocation(t.Timezone)
	t_ = t_.In(loc)

	if loc.String() != targetLoc.String() && t.Time != "" {
		t_ = t_.In(targetLoc)
	}

	return t_, nil
}

type OutputFormat uint

const (
	TimeFormat OutputFormat = 0
	TimeOrigin              = 1
)

type CrossTime struct {
	BeginAt      EFTime       `json:"begin_at,omitempty"`
	Origin       string       `json:"origin,omitempty"`
	OutputFormat OutputFormat `json:"output_format,omitempty"`
}

func (t CrossTime) StringInZone(targetZone string) (string, error) {
	if t.Origin == "" {
		return "", nil
	}

	loc, err := LoadLocation(t.BeginAt.Timezone)
	if err != nil {
		return "", fmt.Errorf("timezone(%s) invalid: %s", t.BeginAt.Timezone, err)
	}
	targetLoc := loc
	if targetZone != "" {
		targetLoc, err = LoadLocation(targetZone)
		if err != nil {
			return "", fmt.Errorf("targetZone(%s) invalid: %s", targetZone, err)
		}
	}

	switch t.OutputFormat {
	case TimeFormat:
		ret, err := t.BeginAt.StringInZone(targetZone, targetLoc)
		return ret, err
	}

	if loc.String() == targetLoc.String() {
		return t.Origin, nil
	}
	return fmt.Sprintf("%s %s", t.Origin, t.BeginAt.Timezone), nil
}
