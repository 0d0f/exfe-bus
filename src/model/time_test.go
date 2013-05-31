package model

import (
	"fmt"
	"testing"
	"time"
)

func shouldEqual(got, expect interface{}) (isFail bool, info string) {
	isFail = !(got == expect)
	if !isFail {
		return
	}
	info = fmt.Sprintf("Got: %s, Expect: %s", got, expect)
	return
}

type ZoneTest struct {
	input         string
	shouldSuccess bool
	expect        string
}

var zoneTestData = []ZoneTest{
	{"+08:00 CST", true, "+08:00"},
	{"+08:00", true, "+08:00"},
	{"-08:00 PST", true, "-08:00"},
	{"-08:00", true, "-08:00"},

	{"+0800 CST", true, "+08:00"},
	{"+0800", true, "+08:00"},
	{"-0800 PST", true, "-08:00"},
	{"-0800", true, "-08:00"},

	{"+8 CST", false, ""},
	{"+8", false, ""},
}

func TestZoneToLocation(t *testing.T) {
	for i, data := range zoneTestData {
		got, err := LoadLocation(data.input)
		if data.shouldSuccess {
			if err != nil {
				t.Fatalf("Test %v should success, but got error: %s", data, err)
			}
			if got.String() != data.expect {
				t.Errorf("Test %v expect: %s, but got: %s", data, data.expect, got.String())
			}
		} else {
			if err == nil {
				t.Fatalf("Test %v should failed, but got: %s", i, got.String())
			}
		}
	}
}

type CrossTimeTest struct {
	time       CrossTime
	targetZone string
	expect     string
	title      string
	desc       string
}

var crossTimeTestData = []CrossTimeTest{
	// if OutputOrigin, then output origin directly
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 2:08:00 pm abc", TimeOrigin}, "+08:00 CST", "2012-04-04 2:08:00 pm abc", "2012-04-04 2:08:00 pm abc", ""},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:08:00", TimeFormat}, "+08:00 CST", "2:08PM on Wed, Apr 4", "Apr 4", "2:08PM on Wed"},
	{CrossTime{EFTime{"", "", "", "", "+08:00 CST"}, "", TimeFormat}, "+08:00 CST", "", "", ""},

	// Time_word (at) Time Date_word (on) Date
	{CrossTime{EFTime{"This Week", "", "", "", "+08:00 CST"}, "This week", TimeFormat}, "+08:00 CST", "This Week", "This Week", ""},
	{CrossTime{EFTime{"", "2012-04-04", "", "", "+08:00 CST"}, "2012 4 4", TimeFormat}, "+08:00 CST", "Wed, Apr 4", "Apr 4", "Wed"},
	{CrossTime{EFTime{"", "", "Dinner", "", "+08:00 CST"}, "dinner", TimeFormat}, "+08:00 CST", "Dinner", "Dinner", ""},
	{CrossTime{EFTime{"", "", "", "06:08:00", "+08:00 CST"}, "14:08:00", TimeFormat}, "+08:00 CST", "2:08PM", "2:08PM", ""},
	{CrossTime{EFTime{"This Week", "2012-04-04", "", "", "+08:00 CST"}, "This week 2012 04 04", TimeFormat}, "+08:00 CST", "This Week on Wed, Apr 4", "Apr 4", "This Week on Wed"},
	{CrossTime{EFTime{"This Week", "", "Dinner", "", "+08:00 CST"}, "dinner this week", TimeFormat}, "+08:00 CST", "Dinner This Week", "This Week", "Dinner"},
	{CrossTime{EFTime{"This Week", "", "", "06:08:00", "+08:00 CST"}, "14:08 this week", TimeFormat}, "+08:00 CST", "2:08PM This Week", "This Week", "2:08PM"},
	{CrossTime{EFTime{"", "2012-04-04", "Dinner", "", "+08:00 CST"}, "dinner 2012-04-04", TimeFormat}, "+08:00 CST", "Dinner on Wed, Apr 4", "Apr 4", "Dinner on Wed"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012 04 04 14:08", TimeFormat}, "+08:00 CST", "2:08PM on Wed, Apr 4", "Apr 4", "2:08PM on Wed"},
	{CrossTime{EFTime{"", "", "Dinner", "06:08:00", "+08:00 CST"}, "dinner at 14:08", TimeFormat}, "+08:00 CST", "Dinner at 2:08PM", "Dinner at 2:08PM", ""},
	{CrossTime{EFTime{"This Week", "2012-04-04", "Dinner", "", "+08:00 CST"}, "dinner this week 2012-04-04", TimeFormat}, "+08:00 CST", "Dinner This Week on Wed, Apr 4", "Apr 4", "Dinner This Week on Wed"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "14:08 this week 2012-04-04", TimeFormat}, "+08:00 CST", "2:08PM This Week on Wed, Apr 4", "Apr 4", "2:08PM This Week on Wed"},
	{CrossTime{EFTime{"This Week", "", "Dinner", "06:08:00", "+08:00 CST"}, "dinner 14:08 this week", TimeFormat}, "+08:00 CST", "Dinner at 2:08PM This Week", "This Week", "Dinner at 2:08PM"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "Dinner", "06:08:00", "+08:00 CST"}, "dinner 14:08 this week 2012-4-4", TimeFormat}, "+08:00 CST", "Dinner at 2:08PM This Week on Wed, Apr 4", "Apr 4", "Dinner at 2:08PM This Week on Wed"},

	// different target zone Timeformat
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", TimeFormat}, "+08:00", "2:08PM on Wed, Apr 4", "Apr 4", "2:08PM on Wed"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", TimeFormat}, "", "2:08PM on Wed, Apr 4", "Apr 4", "2:08PM on Wed"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", TimeFormat}, "+08:00 PST", "2:08PM on Wed, Apr 4", "Apr 4", "2:08PM on Wed"},

	// if Origin, use CrossTime zone
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", TimeFormat}, "+09:00 PST", "3:08PM +09:00 PST on Wed, Apr 4", "Apr 4", "3:08PM +09:00 PST on Wed"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00 abc", TimeOrigin}, "+09:00 PST", "2012-04-04 14:8:00 abc +08:00 CST", "2012-04-04 14:8:00 abc +08:00 CST", ""},

	// Time_word (at) Time Zone Date_word (on) Date
	// Only show Zone with Time_word or Time
	{CrossTime{EFTime{"This Week", "", "", "", "+08:00 CST"}, "this week", TimeFormat}, "+09:00 PST", "This Week", "This Week", ""},
	{CrossTime{EFTime{"", "2012-04-04", "", "", "+08:00 CST"}, "2012-04-04", TimeFormat}, "+09:00 PST", "Wed, Apr 4", "Apr 4", "Wed"},
	{CrossTime{EFTime{"", "", "Dinner", "", "+08:00 CST"}, "dinner", TimeFormat}, "+09:00 PST", "Dinner +08:00 CST", "Dinner +08:00 CST", ""},
	{CrossTime{EFTime{"", "", "", "06:08:00", "+08:00 CST"}, "14:08", TimeFormat}, "+09:00 PST", "3:08PM +09:00 PST", "3:08PM +09:00 PST", ""},
	{CrossTime{EFTime{"This Week", "2012-04-04", "", "", "+08:00 CST"}, "this week 2012 4 4", TimeFormat}, "+09:00 PST", "This Week on Wed, Apr 4", "Apr 4", "This Week on Wed"},
	{CrossTime{EFTime{"This Week", "", "Dinner", "", "+08:00 CST"}, "dinner this week", TimeFormat}, "+09:00 PST", "Dinner +08:00 CST This Week", "This Week", "Dinner +08:00 CST"},
	{CrossTime{EFTime{"This Week", "", "", "06:08:00", "+08:00 CST"}, "14:08 this week", TimeFormat}, "+09:00 PST", "3:08PM +09:00 PST This Week", "This Week", "3:08PM +09:00 PST"},
	{CrossTime{EFTime{"", "2012-04-04", "Dinner", "", "+08:00 CST"}, "dinner 2012-04-04", TimeFormat}, "+09:00 PST", "Dinner +08:00 CST on Wed, Apr 4", "Apr 4", "Dinner +08:00 CST on Wed"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:08", TimeFormat}, "+09:00 PST", "3:08PM +09:00 PST on Wed, Apr 4", "Apr 4", "3:08PM +09:00 PST on Wed"},
	{CrossTime{EFTime{"", "", "Dinner", "06:08:00", "+08:00 CST"}, "dinner 14:08", TimeFormat}, "+09:00 PST", "Dinner at 3:08PM +09:00 PST", "Dinner at 3:08PM +09:00 PST", ""},
	{CrossTime{EFTime{"This Week", "2012-04-04", "Dinner", "", "+08:00 CST"}, "dinner this week 2012-04-04", TimeFormat}, "+09:00 PST", "Dinner +08:00 CST This Week on Wed, Apr 4", "Apr 4", "Dinner +08:00 CST This Week on Wed"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "14:08 this week 2012 04 04", TimeFormat}, "+09:00 PST", "3:08PM +09:00 PST This Week on Wed, Apr 4", "Apr 4", "3:08PM +09:00 PST This Week on Wed"},
	{CrossTime{EFTime{"This Week", "", "Dinner", "06:08:00", "+08:00 CST"}, "14:08 dinner this week", TimeFormat}, "+09:00 PST", "Dinner at 3:08PM +09:00 PST This Week", "This Week", "Dinner at 3:08PM +09:00 PST"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "Dinner", "06:08:00", "+08:00 CST"}, "14:08 dinner this week 2012 04 04", TimeFormat}, "+09:00 PST", "Dinner at 3:08PM +09:00 PST This Week on Wed, Apr 4", "Apr 4", "Dinner at 3:08PM +09:00 PST This Week on Wed"},

	// different target zone Timeformat
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", TimeFormat}, "+09:00", "3:08PM +09:00 on Wed, Apr 4", "Apr 4", "3:08PM +09:00 on Wed"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", TimeFormat}, "", "2:08PM on Wed, Apr 4", "Apr 4", "2:08PM on Wed"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", TimeFormat}, "+09:00 PST", "3:08PM +09:00 PST on Wed, Apr 4", "Apr 4", "3:08PM +09:00 PST on Wed"},

	// different year
	// Time_word (at) Time Date_word (on) Date
	{CrossTime{EFTime{"", "2011-04-04", "", "", "+08:00 CST"}, "2012-04-04", TimeFormat}, "+08:00 CST", "Mon, Apr 4 2011", "Apr 4 2011", "Mon"},
	{CrossTime{EFTime{"This Week", "2011-04-04", "", "", "+08:00 CST"}, "this week 2012-04-04", TimeFormat}, "+08:00 CST", "This Week on Mon, Apr 4 2011", "Apr 4 2011", "This Week on Mon"},
	{CrossTime{EFTime{"", "2011-04-04", "Dinner", "", "+08:00 CST"}, "dinner 2012-04-04", TimeFormat}, "+08:00 CST", "Dinner on Mon, Apr 4 2011", "Apr 4 2011", "Dinner on Mon"},
	{CrossTime{EFTime{"", "2011-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:08", TimeFormat}, "+08:00 CST", "2:08PM on Mon, Apr 4 2011", "Apr 4 2011", "2:08PM on Mon"},
	{CrossTime{EFTime{"This Week", "2011-04-04", "Dinner", "", "+08:00 CST"}, "2012-04-04 dinner this week", TimeFormat}, "+08:00 CST", "Dinner This Week on Mon, Apr 4 2011", "Apr 4 2011", "Dinner This Week on Mon"},
	{CrossTime{EFTime{"This Week", "2011-04-04", "", "06:08:00", "+08:00 CST"}, "this week 2012-04-04 14:8:00", TimeFormat}, "+08:00 CST", "2:08PM This Week on Mon, Apr 4 2011", "Apr 4 2011", "2:08PM This Week on Mon"},
	{CrossTime{EFTime{"This Week", "2011-04-04", "Dinner", "06:08:00", "+08:00 CST"}, "14:08 this week 2012 04 04", TimeFormat}, "+08:00 CST", "Dinner at 2:08PM This Week on Mon, Apr 4 2011", "Apr 4 2011", "Dinner at 2:08PM This Week on Mon"},
}

func TestCrossTimeInZoneString(t *testing.T) {
	nowFunc = func() time.Time {
		return time.Date(2012, 4, 4, 0, 0, 0, 0, time.FixedZone("utc", 0))
	}
	for i, data := range crossTimeTestData {
		got, err := data.time.StringInZone(data.targetZone)
		if err != nil {
			t.Fatalf("Test %+v(%d) should success, but got error: %s", data, i, err)
		}
		if got != data.expect {
			t.Errorf("Test %+v(%d) expect: %s, but got: %s", data, i, data.expect, got)
		}
	}
}

func TestCrossTitleString(t *testing.T) {
	nowFunc = func() time.Time {
		return time.Date(2012, 4, 4, 0, 0, 0, 0, time.FixedZone("utc", 0))
	}
	for i, data := range crossTimeTestData {
		got, err := data.time.Title(data.targetZone)
		if err != nil {
			t.Fatalf("Test %d %+v should success, but got error: %s", i, data, err)
		}
		if got != data.title {
			t.Errorf("Test %d %+v expect: %s, but got: %s", i, data, data.title, got)
		}
	}
}

func TestCrossDescriptioinString(t *testing.T) {
	nowFunc = func() time.Time {
		return time.Date(2012, 4, 4, 0, 0, 0, 0, time.FixedZone("utc", 0))
	}
	for i, data := range crossTimeTestData {
		got, err := data.time.Description(data.targetZone)
		if err != nil {
			t.Fatalf("Test %d %+v should success, but got error: %s", i, data, err)
		}
		if got != data.desc {
			title, _ := data.time.Title(data.targetZone)
			t.Errorf("Test %d %+v expect: %s, but got: %s w/ %s", i, data, data.desc, got, title)
		}
	}
}
