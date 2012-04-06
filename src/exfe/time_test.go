package exfe

import (
	"testing"
	"time"
	"fmt"
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
	input string
	shouldSuccess bool
	expect string
}

var zoneTestData = []ZoneTest{
	{"+08:00 CST", true, "+08:00"},
	{"+08:00",     true, "+08:00"},
	{"-08:00 PST", true, "-08:00"},
	{"-08:00",     true, "-08:00"},

	{"+08:00 ",    false, ""},
	{"+8:00 CST",  false, ""},
	{"-8:00 PST",  false, ""},
	{"+8 CST",     false, ""},
	{"+8",         false, ""},
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

func ThisYear() string {
	now := time.Now()
	return now.Format("Mon, Jan 2")
}

func ThisYearDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func LastYear() string {
	now := time.Now()
	last := now.AddDate(-1, 0, 0)
	return last.Format("Mon, Jan 2 2006")
}

func LastYearDate() string {
	now := time.Now()
	last := now.AddDate(-1, 0, 0)
	return last.Format("2006-01-02")
}

type CrossTimeTest struct {
	time CrossTime
	targetZone string
	shouldSuccess bool
	expect string
}

var crossTimeTestData = []CrossTimeTest{
	// if OutputOrigin, then output origin directly
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "+08:00 CST", true, fmt.Sprintf("2:08PM on %s", ThisYear())},
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00 abc", OutputOrigin},            "+08:00 CST", true, fmt.Sprintf("2012-04-04 14:8:00 abc")},

	// Time_word (at) Time Date_word (on) Date
	{CrossTime{EFTime{"This Week", "", "", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                           "+08:00 CST", true, fmt.Sprintf("This Week")},
	{CrossTime{EFTime{"", ThisYearDate(), "", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                        "+08:00 CST", true, fmt.Sprintf("%s", ThisYear())},
	{CrossTime{EFTime{"", "", "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                              "+08:00 CST", true, fmt.Sprintf("Dinner")},
	{CrossTime{EFTime{"", "", "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                            "+08:00 CST", true, fmt.Sprintf("2:08PM")},
	{CrossTime{EFTime{"This Week", ThisYearDate(), "", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},               "+08:00 CST", true, fmt.Sprintf("This Week on %s", ThisYear())},
	{CrossTime{EFTime{"This Week", "", "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                     "+08:00 CST", true, fmt.Sprintf("Dinner This Week")},
	{CrossTime{EFTime{"This Week", "", "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                   "+08:00 CST", true, fmt.Sprintf("2:08PM This Week")},
	{CrossTime{EFTime{"", ThisYearDate(), "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                  "+08:00 CST", true, fmt.Sprintf("Dinner on %s", ThisYear())},
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "+08:00 CST", true, fmt.Sprintf("2:08PM on %s", ThisYear())},
	{CrossTime{EFTime{"", "", "Dinner", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                      "+08:00 CST", true, fmt.Sprintf("Dinner at 2:08PM")},
	{CrossTime{EFTime{"This Week", ThisYearDate(), "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},         "+08:00 CST", true, fmt.Sprintf("Dinner This Week on %s", ThisYear())},
	{CrossTime{EFTime{"This Week", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},       "+08:00 CST", true, fmt.Sprintf("2:08PM This Week on %s", ThisYear())},
	{CrossTime{EFTime{"This Week", "", "Dinner", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},             "+08:00 CST", true, fmt.Sprintf("Dinner at 2:08PM This Week")},
	{CrossTime{EFTime{"This Week", ThisYearDate(), "Dinner", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat}, "+08:00 CST", true, fmt.Sprintf("Dinner at 2:08PM This Week on %s", ThisYear())},

	// different target zone format
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "+08:00", true, fmt.Sprintf("2:08PM on %s", ThisYear())},
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "", true, fmt.Sprintf("2:08PM on %s", ThisYear())},
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "+08:00 PST", true, fmt.Sprintf("2:08PM on %s", ThisYear())},

	// if OutputOrigin, use CrossTime zone
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "+09:00 PST", true, fmt.Sprintf("3:08PM +09:00 PST on %s", ThisYear())},
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00 abc", OutputOrigin},            "+09:00 PST", true, fmt.Sprintf("2012-04-04 14:8:00 abc +08:00 CST")},

	// Time_word (at) Time Zone Date_word (on) Date
	// Only show Zone with Time_word or Time
	{CrossTime{EFTime{"This Week", "", "", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                           "+09:00 PST", true, fmt.Sprintf("This Week")},
	{CrossTime{EFTime{"", ThisYearDate(), "", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                        "+09:00 PST", true, fmt.Sprintf("%s", ThisYear())},
	{CrossTime{EFTime{"", "", "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                              "+09:00 PST", true, fmt.Sprintf("Dinner +08:00 CST")},
	{CrossTime{EFTime{"", "", "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                            "+09:00 PST", true, fmt.Sprintf("3:08PM +09:00 PST")},
	{CrossTime{EFTime{"This Week", ThisYearDate(), "", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},               "+09:00 PST", true, fmt.Sprintf("This Week on %s", ThisYear())},
	{CrossTime{EFTime{"This Week", "", "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                     "+09:00 PST", true, fmt.Sprintf("Dinner +08:00 CST This Week")},
	{CrossTime{EFTime{"This Week", "", "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                   "+09:00 PST", true, fmt.Sprintf("3:08PM +09:00 PST This Week")},
	{CrossTime{EFTime{"", ThisYearDate(), "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                  "+09:00 PST", true, fmt.Sprintf("Dinner +08:00 CST on %s", ThisYear())},
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "+09:00 PST", true, fmt.Sprintf("3:08PM +09:00 PST on %s", ThisYear())},
	{CrossTime{EFTime{"", "", "Dinner", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                      "+09:00 PST", true, fmt.Sprintf("Dinner at 3:08PM +09:00 PST")},
	{CrossTime{EFTime{"This Week", ThisYearDate(), "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},         "+09:00 PST", true, fmt.Sprintf("Dinner +08:00 CST This Week on %s", ThisYear())},
	{CrossTime{EFTime{"This Week", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},       "+09:00 PST", true, fmt.Sprintf("3:08PM +09:00 PST This Week on %s", ThisYear())},
	{CrossTime{EFTime{"This Week", "", "Dinner", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},             "+09:00 PST", true, fmt.Sprintf("Dinner at 3:08PM +09:00 PST This Week")},
	{CrossTime{EFTime{"This Week", ThisYearDate(), "Dinner", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat}, "+09:00 PST", true, fmt.Sprintf("Dinner at 3:08PM +09:00 PST This Week on %s", ThisYear())},

	// different target zone format
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "+09:00", true, fmt.Sprintf("3:08PM +09:00 on %s", ThisYear())},
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "", true, fmt.Sprintf("2:08PM on %s", ThisYear())},
	{CrossTime{EFTime{"", ThisYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "+09:00 PST", true, fmt.Sprintf("3:08PM +09:00 PST on %s", ThisYear())},

	// different year
	// Time_word (at) Time Date_word (on) Date
	{CrossTime{EFTime{"", LastYearDate(), "", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                        "+08:00 CST", true, fmt.Sprintf("%s", LastYear())},
	{CrossTime{EFTime{"This Week", LastYearDate(), "", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},               "+08:00 CST", true, fmt.Sprintf("This Week on %s", LastYear())},
	{CrossTime{EFTime{"", LastYearDate(), "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                  "+08:00 CST", true, fmt.Sprintf("Dinner on %s", LastYear())},
	{CrossTime{EFTime{"", LastYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},                "+08:00 CST", true, fmt.Sprintf("2:08PM on %s", LastYear())},
	{CrossTime{EFTime{"This Week", LastYearDate(), "Dinner", "", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},         "+08:00 CST", true, fmt.Sprintf("Dinner This Week on %s", LastYear())},
	{CrossTime{EFTime{"This Week", LastYearDate(), "", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat},       "+08:00 CST", true, fmt.Sprintf("2:08PM This Week on %s", LastYear())},
	{CrossTime{EFTime{"This Week", LastYearDate(), "Dinner", "14:08:00", "+08:00 CST"}, "2012-04-04 14:8:00", OutputFormat}, "+08:00 CST", true, fmt.Sprintf("Dinner at 2:08PM This Week on %s", LastYear())},
}

func TestCrossTimeInZoneString(t *testing.T) {
	for _, data := range crossTimeTestData {
		got, err := data.time.StringInZone(data.targetZone)
		if data.shouldSuccess {
			if err != nil {
				t.Fatalf("Test %v should success, but got error: %s", data, err)
			}
			if got != data.expect {
				t.Errorf("Test %v expect: %s, but got: %s", data, data.expect, got)
			}
		} else {
			if err == nil {
				t.Fatalf("Test %v should failed, but got: %s", got)
			}
		}
	}
}
