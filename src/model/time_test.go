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

	{"+08:00 ", false, ""},
	{"+8:00 CST", false, ""},
	{"-8:00 PST", false, ""},
	{"+8 CST", false, ""},
	{"+8", false, ""},
}

func TestZoneToLocation(t *testing.T) {
	for _, data := range zoneTestData {
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
				t.Fatalf("Test %v should failed, but got: %s", got.String())
			}
		}
	}
}

type CrossTimeTest struct {
	time       CrossTime
	targetZone string
	expect     string
}

var crossTimeTestData = []CrossTimeTest{
	// if OutputOrigin, then output origin directly
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 2:08:00 pm abc", Origin}, "+08:00 CST", "2012-04-04 2:08:00 pm abc"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:08:00", Format}, "+08:00 CST", "2:08PM on Wed, Apr 4"},
	{CrossTime{EFTime{"", "", "", "", "+08:00 CST"}, "", Format}, "+08:00 CST", ""},

	// Time_word (at) Time Date_word (on) Date
	{CrossTime{EFTime{"This Week", "", "", "", "+08:00 CST"}, "This week", Format}, "+08:00 CST", "This Week"},
	{CrossTime{EFTime{"", "2012-04-04", "", "", "+08:00 CST"}, "2012 4 4", Format}, "+08:00 CST", "Wed, Apr 4"},
	{CrossTime{EFTime{"", "", "Dinner", "", "+08:00 CST"}, "dinner", Format}, "+08:00 CST", "Dinner"},
	{CrossTime{EFTime{"", "", "", "06:08:00", "+08:00 CST"}, "14:08:00", Format}, "+08:00 CST", "2:08PM"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "", "", "+08:00 CST"}, "This week 2012 04 04", Format}, "+08:00 CST", "This Week on Wed, Apr 4"},
	{CrossTime{EFTime{"This Week", "", "Dinner", "", "+08:00 CST"}, "dinner this week", Format}, "+08:00 CST", "Dinner This Week"},
	{CrossTime{EFTime{"This Week", "", "", "06:08:00", "+08:00 CST"}, "14:08 this week", Format}, "+08:00 CST", "2:08PM This Week"},
	{CrossTime{EFTime{"", "2012-04-04", "Dinner", "", "+08:00 CST"}, "dinner 2012-04-04", Format}, "+08:00 CST", "Dinner on Wed, Apr 4"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012 04 04 14:08", Format}, "+08:00 CST", "2:08PM on Wed, Apr 4"},
	{CrossTime{EFTime{"", "", "Dinner", "06:08:00", "+08:00 CST"}, "dinner at 14:08", Format}, "+08:00 CST", "Dinner at 2:08PM"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "Dinner", "", "+08:00 CST"}, "dinner this week 2012-04-04", Format}, "+08:00 CST", "Dinner This Week on Wed, Apr 4"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "14:08 this week 2012-04-04", Format}, "+08:00 CST", "2:08PM This Week on Wed, Apr 4"},
	{CrossTime{EFTime{"This Week", "", "Dinner", "06:08:00", "+08:00 CST"}, "dinner 14:08 this week", Format}, "+08:00 CST", "Dinner at 2:08PM This Week"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "Dinner", "06:08:00", "+08:00 CST"}, "dinner 14:08 this week 2012-4-4", Format}, "+08:00 CST", "Dinner at 2:08PM This Week on Wed, Apr 4"},

	// different target zone format
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", Format}, "+08:00", "2:08PM on Wed, Apr 4"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", Format}, "", "2:08PM on Wed, Apr 4"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", Format}, "+08:00 PST", "2:08PM on Wed, Apr 4"},

	// if Origin, use CrossTime zone
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", Format}, "+09:00 PST", "3:08PM +09:00 PST on Wed, Apr 4"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00 abc", Origin}, "+09:00 PST", "2012-04-04 14:8:00 abc +08:00 CST"},

	// Time_word (at) Time Zone Date_word (on) Date
	// Only show Zone with Time_word or Time
	{CrossTime{EFTime{"This Week", "", "", "", "+08:00 CST"}, "this week", Format}, "+09:00 PST", "This Week"},
	{CrossTime{EFTime{"", "2012-04-04", "", "", "+08:00 CST"}, "2012-04-04", Format}, "+09:00 PST", "Wed, Apr 4"},
	{CrossTime{EFTime{"", "", "Dinner", "", "+08:00 CST"}, "dinner", Format}, "+09:00 PST", "Dinner +08:00 CST"},
	{CrossTime{EFTime{"", "", "", "06:08:00", "+08:00 CST"}, "14:08", Format}, "+09:00 PST", "3:08PM +09:00 PST"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "", "", "+08:00 CST"}, "this week 2012 4 4", Format}, "+09:00 PST", "This Week on Wed, Apr 4"},
	{CrossTime{EFTime{"This Week", "", "Dinner", "", "+08:00 CST"}, "dinner this week", Format}, "+09:00 PST", "Dinner +08:00 CST This Week"},
	{CrossTime{EFTime{"This Week", "", "", "06:08:00", "+08:00 CST"}, "14:08 this week", Format}, "+09:00 PST", "3:08PM +09:00 PST This Week"},
	{CrossTime{EFTime{"", "2012-04-04", "Dinner", "", "+08:00 CST"}, "dinner 2012-04-04", Format}, "+09:00 PST", "Dinner +08:00 CST on Wed, Apr 4"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:08", Format}, "+09:00 PST", "3:08PM +09:00 PST on Wed, Apr 4"},
	{CrossTime{EFTime{"", "", "Dinner", "06:08:00", "+08:00 CST"}, "dinner 14:08", Format}, "+09:00 PST", "Dinner at 3:08PM +09:00 PST"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "Dinner", "", "+08:00 CST"}, "dinner this week 2012-04-04", Format}, "+09:00 PST", "Dinner +08:00 CST This Week on Wed, Apr 4"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "14:08 this week 2012 04 04", Format}, "+09:00 PST", "3:08PM +09:00 PST This Week on Wed, Apr 4"},
	{CrossTime{EFTime{"This Week", "", "Dinner", "06:08:00", "+08:00 CST"}, "14:08 dinner this week", Format}, "+09:00 PST", "Dinner at 3:08PM +09:00 PST This Week"},
	{CrossTime{EFTime{"This Week", "2012-04-04", "Dinner", "06:08:00", "+08:00 CST"}, "14:08 dinner this week 2012 04 04", Format}, "+09:00 PST", "Dinner at 3:08PM +09:00 PST This Week on Wed, Apr 4"},

	// different target zone format
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", Format}, "+09:00", "3:08PM +09:00 on Wed, Apr 4"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", Format}, "", "2:08PM on Wed, Apr 4"},
	{CrossTime{EFTime{"", "2012-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", Format}, "+09:00 PST", "3:08PM +09:00 PST on Wed, Apr 4"},

	// different year
	// Time_word (at) Time Date_word (on) Date
	{CrossTime{EFTime{"", "2011-04-04", "", "", "+08:00 CST"}, "2012-04-04", Format}, "+08:00 CST", "Mon, Apr 4 2011"},
	{CrossTime{EFTime{"This Week", "2011-04-04", "", "", "+08:00 CST"}, "this week 2012-04-04", Format}, "+08:00 CST", "This Week on Mon, Apr 4 2011"},
	{CrossTime{EFTime{"", "2011-04-04", "Dinner", "", "+08:00 CST"}, "dinner 2012-04-04", Format}, "+08:00 CST", "Dinner on Mon, Apr 4 2011"},
	{CrossTime{EFTime{"", "2011-04-04", "", "06:08:00", "+08:00 CST"}, "2012-04-04 14:08", Format}, "+08:00 CST", "2:08PM on Mon, Apr 4 2011"},
	{CrossTime{EFTime{"This Week", "2011-04-04", "Dinner", "", "+08:00 CST"}, "2012-04-04 dinner this week", Format}, "+08:00 CST", "Dinner This Week on Mon, Apr 4 2011"},
	{CrossTime{EFTime{"This Week", "2011-04-04", "", "06:08:00", "+08:00 CST"}, "this week 2012-04-04 14:8:00", Format}, "+08:00 CST", "2:08PM This Week on Mon, Apr 4 2011"},
	{CrossTime{EFTime{"This Week", "2011-04-04", "Dinner", "06:08:00", "+08:00 CST"}, "14:08 this week 2012 04 04", Format}, "+08:00 CST", "Dinner at 2:08PM This Week on Mon, Apr 4 2011"},
}

func TestCrossTimeInZoneString(t *testing.T) {
	nowFunc = func() time.Time {
		return time.Date(2012, 4, 4, 0, 0, 0, 0, time.FixedZone("utc", 0))
	}
	for _, data := range crossTimeTestData {
		got, err := data.time.StringInZone(data.targetZone)
		if err != nil {
			t.Fatalf("Test %v should success, but got error: %s", data, err)
		}
		if got != data.expect {
			t.Errorf("Test %+v expect: %s, but got: %s", data, data.expect, got)
		}
	}
}
