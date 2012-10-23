package model

import (
	"fmt"
	"time"
)

const (
	HourInSeconds   int = 1 /*hour*/ * 60 /*minutes*/ * 60 /*seconds*/
	MinuteInSeconds int = 1 /*minute*/ * 60                /*seconds*/
)

func LoadLocation(zone string) (*time.Location, error) {
	var hour, minute int
	_, err := fmt.Sscanf(zone, "+%02d%02d", &hour, &minute)
	if err != nil {
		_, err = fmt.Sscanf(zone, "-%02d%02d", &hour, &minute)
		hour = -hour
	}
	if err != nil {
		_, err = fmt.Sscanf(zone, "+%02d:%02d", &hour, &minute)
	}
	if err != nil {
		_, err = fmt.Sscanf(zone, "-%02d:%02d", &hour, &minute)
		hour = -hour
	}
	if err != nil {
		return nil, fmt.Errorf("Zone format invalid")
	}

	offset := hour*HourInSeconds + minute*MinuteInSeconds
	return time.FixedZone(fmt.Sprintf("%+03d:%02d", hour, minute), offset), nil
}
