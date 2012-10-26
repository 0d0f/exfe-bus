package exfe_model

import (
	"fmt"
	"regexp"
	"time"
)

const (
	HourInSeconds   int = 1 /*hour*/ * 60 /*minutes*/ * 60 /*seconds*/
	MinuteInSeconds int = 1 /*minute*/ * 60                /*seconds*/
)

func LoadLocation(zone string) (*time.Location, error) {
	isOK, err := regexp.MatchString("^[+-]\\d\\d:\\d\\d( [A-Z]{3})?$", zone)
	if err != nil {
		return nil, err
	}
	if !isOK {
		return nil, fmt.Errorf("Zone format not fit /^[+-]\\d\\d:\\d\\d( [A-Z]{3})?$/")
	}

	var hour, minute int
	fmt.Sscanf(zone, "%d:%d", &hour, &minute)
	offset := hour*HourInSeconds + minute*MinuteInSeconds
	return time.FixedZone(fmt.Sprintf("%+03d:%02d", hour, minute), offset), nil
}

type EFTime struct {
	Date_word string
	Date      string
	Time_word string
	Time      string
	Timezone  string
}

func (t EFTime) differentZone(targetZone string) bool {
	return targetZone != "" && t.Timezone[0:6] != targetZone[0:6]
}

func (t EFTime) timeInZone(targetZone string) (time.Time, error) {
	var t_ time.Time
	var err error

	switch {
	case t.Time != "" && t.Date != "":
		t_, err = time.Parse("2006-1-2 15:4:5", fmt.Sprintf("%s %s", t.Date, t.Time[:8]))
	case t.Time != "" && t.Date == "":
		t_, err = time.Parse("15:4:5", fmt.Sprintf("%s", t.Time[:8]))
	case t.Time == "" && t.Date != "":
		t_, err = time.Parse("2006-1-2", fmt.Sprintf("%s", t.Date))
	}

	if err != nil {
		return t_, fmt.Errorf("Parse time error: %s", err)
	}

	loc, err := LoadLocation(t.Timezone)
	if err != nil {
		return t_, fmt.Errorf("Parse timezone(%s) error: %s", t.Timezone, err)
	}
	t_ = t_.In(loc)

	if t.differentZone(targetZone) && t.Time != "" {
		targetLocation, err := LoadLocation(targetZone)
		if err != nil {
			return t_, fmt.Errorf("Parse target zone(%s) error: %s", targetZone, err)
		}
		t_ = t_.In(targetLocation)
	}

	return t_, nil
}

func (t EFTime) StringInZone(targetZone string) (string, error) {
	t_, err := t.timeInZone(targetZone)
	if err != nil {
		return "", err
	}

	ret := ""

	if t.Time_word != "" {
		ret += t.Time_word
	}

	if t.Time != "" {
		if ret != "" {
			ret += " at "
		}
		ret += t_.Format("3:04PM")
	}

	if t.differentZone(targetZone) && ret != "" {
		ret += " "
		if t.Time != "" {
			ret += targetZone
		} else {
			ret += t.Timezone
		}
	}

	if t.Date_word != "" {
		if ret != "" {
			ret += " "
		}
		ret += t.Date_word
	}

	if t.Date != "" {
		if ret != "" {
			ret += " on "
		}
		now := time.Now()
		if now.Year() == t_.Year() {
			ret += t_.Format("Mon, Jan 2")
		} else {
			ret += t_.Format("Mon, Jan 2 2006")
		}
	}

	return ret, nil
}

func (t EFTime) UTCTime(layout string) (string, error) {
	if t.Time == "" && t.Date == "" {
		return "", nil
	}

	var time_ time.Time
	var err error
	if t.Time == "" {
		time_, err = time.Parse("2006-1-2", t.Date)
	} else if t.Date != "" {
		str := fmt.Sprintf("%s %s", t.Date, t.Time)
		time_, err = time.Parse("2006-1-2 15:4:5", str)
	}
	if err != nil {
		return "", err
	}
	return time_.Format(layout), nil
}

type OutputFormat uint

const (
	Format OutputFormat = 0
	Origin              = 1
)

type CrossTime struct {
	Begin_at     EFTime
	Origin       string
	OutputFormat OutputFormat
}

func (t CrossTime) StringInZone(targetZone string) (string, error) {
	switch t.OutputFormat {
	case Format:
		ret, err := t.Begin_at.StringInZone(targetZone)
		return ret, err
	}

	if targetZone[0:6] == t.Begin_at.Timezone[0:6] {
		return t.Origin, nil
	}
	return fmt.Sprintf("%s %s", t.Origin, t.Begin_at.Timezone), nil
}
