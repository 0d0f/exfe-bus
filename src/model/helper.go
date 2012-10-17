package model

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
