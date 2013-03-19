package ics

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strings"
	"time"
)

type Timezone struct {
	ID       string
	Standard struct {
		Start      time.Time
		OffsetFrom string
	}
}

type Attendee struct {
	Name     string
	Email    string
	PartStat string
}

type Event struct {
	ID          string
	Organizer   Attendee
	Start       time.Time
	DateStart   bool
	End         time.Time
	DateEnd     bool
	Location    string
	Description string
	URL         string
	Summary     string
	Attendees   []Attendee
}

type Calendar struct {
	Event []Event
}

func ParseCalendar(reader io.Reader) (Calendar, error) {
	r := textproto.NewReader(bufio.NewReader(reader))
	ret := Calendar{}
	tzs := make(map[string]Timezone)
	for {
		key, _, value, err := GetIcsLine(r)
		if err != nil {
			return ret, err
		}
		switch key {
		case "BEGIN":
			switch value {
			case "VTIMEZONE":
				tz, err := ParseTimezone(r)
				if err != nil {
					return ret, err
				}
				tzs[tz.ID] = tz
			case "VEVENT":
				event, err := ParseEvent(r, tzs)
				if err != nil {
					return ret, err
				}
				ret.Event = append(ret.Event, event)
			}
		case "END":
			if value == "VCALENDAR" {
				return ret, nil
			}
		}
	}
	return ret, nil
}

func ParseEvent(r *textproto.Reader, tzs map[string]Timezone) (Event, error) {
	ret := Event{}
	for {
		key, params, value, err := GetIcsLine(r)
		if err != nil {
			return ret, err
		}
		switch key {
		case "UID":
			ret.ID = value
		case "DTSTART":
			tz := tzs[params["TZID"]]
			start, err := ParseDateTime(value, tz.Standard.OffsetFrom)
			if err != nil {
				return ret, err
			}
			ret.Start = start
			ret.DateStart = params["VALUE"] == "DATE"
		case "DTEND":
			tz := tzs[params["TZID"]]
			end, err := ParseDateTime(value, tz.Standard.OffsetFrom)
			if err != nil {
				return ret, err
			}
			ret.End = end
			ret.DateEnd = params["VALUE"] == "DATE"
		case "SUMMARY":
			ret.Summary = value
		case "ORGANIZER":
			ret.Organizer = ParseAttendee(params, value)
		case "DESCRIPTION":
			ret.Description = value
		case "LOCATION":
			ret.Location = value
		case "URL":
			ret.URL = value
		case "ATTENDEE":
			ret.Attendees = append(ret.Attendees, ParseAttendee(params, value))
		case "END":
			if value == "VEVENT" {
				return ret, nil
			}
		case "BEGIN":
			err := IgnoreTo(r, value)
			if err != nil {
				return ret, err
			}
		}
	}
	return ret, nil
}

func IgnoreTo(r *textproto.Reader, name string) error {
	for {
		key, _, value, err := GetIcsLine(r)
		if err != nil {
			return err
		}
		if key == "END" && value == name {
			break
		}
	}
	return nil
}

func ParseTimezone(r *textproto.Reader) (Timezone, error) {
	ret := Timezone{}
	for {
		key, _, value, err := GetIcsLine(r)
		if err != nil {
			return ret, err
		}
		switch key {
		case "END":
			if value == "VTIMEZONE" {
				return ret, nil
			}
		case "TZID":
			ret.ID = value
		case "BEGIN":
			switch value {
			case "STANDARD":
				err = ParseTimezoneStandard(r, &ret)
				if err != nil {
					return ret, err
				}
			default:
				err = fmt.Errorf("invalid BEGIN line: %s", value)
				return ret, err
			}
		}
	}
	return ret, nil
}

func ParseTimezoneStandard(r *textproto.Reader, timezone *Timezone) error {
	for {
		key, _, value, err := GetIcsLine(r)
		if err != nil {
			return err
		}
		switch key {
		case "DTSTART":
			time, err := ParseDateTime(value, "")
			if err != nil {
				return err
			}
			timezone.Standard.Start = time
		case "TZOFFSETFROM":
			timezone.Standard.OffsetFrom = value
		case "END":
			if value == "STANDARD" {
				return nil
			}
		}
	}
	return nil
}

func GetIcsLine(r *textproto.Reader) (key string, params map[string]string, value string, err error) {
	replacer := strings.NewReplacer("\\n", "\n")
	var line string
	line, err = r.ReadContinuedLine()
	if err != nil {
		return
	}

	kv := strings.SplitN(line, ":", 2)
	if len(kv) < 2 {
		err = fmt.Errorf("invalid line(%s): no ':' finded", line)
		return
	}
	keys := strings.Split(kv[0], ";")
	key = strings.Trim(keys[0], " \r\n\t")
	value = replacer.Replace(strings.Trim(kv[1], " \r\n\t"))
	for _, v := range keys[1:] {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) < 2 {
			err = fmt.Errorf("invalid line(%s): no '=' in parameters", line)
			return
		}
		if params == nil {
			params = make(map[string]string)
		}
		params[strings.Trim(kv[0], " \r\n\t")] = strings.Trim(kv[1], " \r\n\t")
	}
	return
}

func ParseDateTime(value, tz string) (time.Time, error) {
	if len(value) < len("20060102T150405") {
		return time.Parse("20060102", value)
	}
	if tz == "" && value[len(value)-1] != 'Z' {
		value = value + "Z"
	}
	if value[len(value)-1] == 'Z' {
		return time.Parse("20060102T150405Z", value)
	}
	return time.Parse("20060102T150405 -0700", fmt.Sprintf("%s %s", value, tz))
}

func ParseAttendee(params map[string]string, value string) Attendee {
	kv := strings.SplitN(strings.ToLower(value), ":", 2)
	if len(kv) == 2 {
		if kv[0] == "invalid" {
			value = params["CN"]
		} else {
			value = kv[1]
		}
	}

	return Attendee{
		Name:     params["CN"],
		Email:    value,
		PartStat: params["PARTSTAT"],
	}
}
