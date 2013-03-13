package ics

import (
	"bytes"
	"encoding/json"
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestParseGoogle(t *testing.T) {
	google := `BEGIN:VCALENDAR
PRODID:-//Google Inc//Google Calendar 70.9054//EN
VERSION:2.0
CALSCALE:GREGORIAN
METHOD:REQUEST
BEGIN:VEVENT
DTSTART:20130312T033000Z
DTEND:20130312T043000Z
DTSTAMP:20130311T130930Z
ORGANIZER;CN=Googol Lee:mailto:googollee@gmail.com
UID:qk5r6u61fc1he6oqorfud3hdkg@google.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=
 TRUE;CN=+8613488802891@exfe.com;X-NUM-GUESTS=0:mailto:+8613488802891@exfe.c
 om
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=
 TRUE;CN=googollee@hotmail.com;X-NUM-GUESTS=0:mailto:googollee@hotmail.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=
 TRUE;CN=googollee/twitter@exfe.com;X-NUM-GUESTS=0:mailto:googollee/twitter@
 exfe.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED;RSVP=TRUE
 ;CN=Googol Lee;X-NUM-GUESTS=0:mailto:googollee@gmail.com
CREATED:20130311T130518Z
DESCRIPTION:说明内容\n在 http://www.google.com/calendar/event?action=VIEW&eid=cW
 s1cjZ1NjFmYzFoZTZvcW9yZnVkM2hka2cgZ29vZ29sbGVlQGhvdG1haWwuY29t&tok=MTkjZ29v
 Z29sbGVlQGdtYWlsLmNvbWNhMmQ4MjA3NjNlMDhlMGIyNDJlY2IzODVlYzQzYzFlYjA5MTFmM2I
 &ctz=Asia/Shanghai&hl=zh_CN 查看您的活动。
LAST-MODIFIED:20130311T130930Z
LOCATION:地点
SEQUENCE:0
STATUS:CONFIRMED
SUMMARY:测试
TRANSP:OPAQUE
END:VEVENT
END:VCALENDAR`
	expect := "{\"Event\":[{\"ID\":\"qk5r6u61fc1he6oqorfud3hdkg@google.com\",\"Organizer\":{\"Name\":\"Googol Lee\",\"Email\":\"googollee@gmail.com\",\"PartStat\":\"\"},\"Start\":\"2013-03-12T03:30:00Z\",\"DateStart\":false,\"End\":\"2013-03-12T04:30:00Z\",\"DateEnd\":false,\"Location\":\"地点\",\"Description\":\"说明内容\\n在 http://www.google.com/calendar/event?action=VIEW&eid=cW s1cjZ1NjFmYzFoZTZvcW9yZnVkM2hka2cgZ29vZ29sbGVlQGhvdG1haWwuY29t&tok=MTkjZ29v Z29sbGVlQGdtYWlsLmNvbWNhMmQ4MjA3NjNlMDhlMGIyNDJlY2IzODVlYzQzYzFlYjA5MTFmM2I &ctz=Asia/Shanghai&hl=zh_CN 查看您的活动。\",\"URL\":\"\",\"Summary\":\"测试\",\"Attendees\":[{\"Name\":\"+8613488802891@exfe.com\",\"Email\":\"+8613488802891@exfe.c om\",\"PartStat\":\"NEEDS-ACTION\"},{\"Name\":\"googollee@hotmail.com\",\"Email\":\"googollee@hotmail.com\",\"PartStat\":\"NEEDS-ACTION\"},{\"Name\":\"googollee/twitter@exfe.com\",\"Email\":\"googollee/twitter@ exfe.com\",\"PartStat\":\"NEEDS-ACTION\"},{\"Name\":\"Googol Lee\",\"Email\":\"googollee@gmail.com\",\"PartStat\":\"ACCEPTED\"}]}]}"

	reader := bytes.NewBufferString(google)
	calendar, err := ParseCalendar(reader)
	assert.Equal(t, err, nil)
	j, _ := json.Marshal(calendar)
	assert.Equal(t, string(j), expect)
}

func TestParseICloud(t *testing.T) {
	icloud := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Apple Inc.//Mac OS X 10.8.2//EN
CALSCALE:GREGORIAN
BEGIN:VEVENT
CREATED:20130311T140841Z
UID:A97BCE9C-70F9-4C31-8A5B-C08A5B8F147D
DTEND;VALUE=DATE:20130327
ATTENDEE;CN=GoogolLee LiZhaoHai;CUTYPE=INDIVIDUAL;
 EMAIL=googollee@gmail.com;PARTSTAT=ACCEPTED;
 X-CALENDARSERVER-DTSTAMP=20130311T140935Z:/1342978541/principal/
ATTENDEE;CN=+8613488802891;CUTYPE=INDIVIDUAL;PARTSTAT=NEEDS-ACTION;
 SCHEDULE-STATUS=5.1:invalid:nomail
ATTENDEE;CN=googollee@163.com;CUTYPE=INDIVIDUAL;PARTSTAT=TENTATIVE;
 SCHEDULE-STATUS=2.0:mailto:googollee@163.com
ATTENDEE;CN=lzh@exfe.com;CUTYPE=INDIVIDUAL;PARTSTAT=NEEDS-ACTION;
 SCHEDULE-STATUS=1.1:mailto:lzh@exfe.com
ATTENDEE;CN=李 兆海;CUTYPE=INDIVIDUAL;EMAIL=googollee@hotmail.com;
 ROLE=REQ-PARTICIPANT;PARTSTAT=DECLINED;SCHEDULE-STATUS=2.0:mailto:
 googollee@hotmail.com
SUMMARY:New Event
DTSTART;VALUE=DATE:20130326
ORGANIZER;CN=GoogolLee LiZhaoHai;EMAIL=googollee@gmail.com:
 /1342978541/principal/
SEQUENCE:3
DTSTAMP:20130311T141045Z
TRANSP:TRANSPARENT
BEGIN:VALARM
X-WR-ALARMUID:E2F546B3-A866-4F77-87C7-9CA46A477EE5
UID:E2F546B3-A866-4F77-87C7-9CA46A477EE5
TRIGGER;VALUE=DURATION:-PT15H
X-APPLE-DEFAULT-ALARM:TRUE
ATTACH;VALUE=URI:Basso
ACTION:AUDIO
END:VALARM
END:VEVENT
END:VCALENDAR`
	expect := `{"Event":[{"ID":"E2F546B3-A866-4F77-87C7-9CA46A477EE5","Organizer":{"Name":"GoogolLee LiZhaoHai","Email":"/1342978541/principal/","PartStat":""},"Start":"2013-03-26T00:00:00Z","DateStart":true,"End":"2013-03-27T00:00:00Z","DateEnd":true,"Location":"","Description":"","URL":"","Summary":"New Event","Attendees":[{"Name":"GoogolLee LiZhaoHai","Email":"/1342978541/principal/","PartStat":"ACCEPTED"},{"Name":"+8613488802891","Email":"+8613488802891","PartStat":"NEEDS-ACTION"},{"Name":"googollee@163.com","Email":"googollee@163.com","PartStat":"TENTATIVE"},{"Name":"lzh@exfe.com","Email":"lzh@exfe.com","PartStat":"NEEDS-ACTION"},{"Name":"李 兆海","Email":" googollee@hotmail.com","PartStat":"DECLINED"}]}]}`

	reader := bytes.NewBufferString(icloud)
	calendar, err := ParseCalendar(reader)
	assert.Equal(t, err, nil)
	j, _ := json.Marshal(calendar)
	assert.Equal(t, string(j), expect)
}

func TestParseOutlook(t *testing.T) {
	{
		outlook := `BEGIN:VCALENDAR
METHOD:REQUEST
VERSION:2.0
PRODID:-//Microsoft Corporation//Windows Live Calendar//EN
BEGIN:VTIMEZONE
TZID:China Standard Time
BEGIN:STANDARD
DTSTART:20080101T000000
TZOFFSETTO:+0800
TZOFFSETFROM:+0800
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
UID:d1521bd3-03f0-49ef-a881-731b7067452a
DTSTAMP:20130312T050829Z
CLASS:PUBLIC
X-MICROSOFT-CDO-BUSYSTATUS:BUSY
TRANSP:OPAQUE
SEQUENCE:0
DTSTART;TZID=China Standard Time:20130312T090000
DTEND;TZID=China Standard Time:20130312T100000
SUMMARY:内容
LOCATION:地点
PRIORITY:0
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=
 TRUE:MAILTO:googollee@gmail.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=
 TRUE:MAILTO:googollee/twitter@exfe.com
ORGANIZER;CN=Lee Googol Lee:MAILTO:googollee@hotmail.com
BEGIN:VALARM
ACTION:DISPLAY
TRIGGER:-PT15M
END:VALARM
END:VEVENT
END:VCALENDAR`
		expect := "{\"Event\":[{\"ID\":\"d1521bd3-03f0-49ef-a881-731b7067452a\",\"Organizer\":{\"Name\":\"Lee Googol Lee\",\"Email\":\"googollee@hotmail.com\",\"PartStat\":\"\"},\"Start\":\"2013-03-12T09:00:00+08:00\",\"DateStart\":false,\"End\":\"2013-03-12T10:00:00+08:00\",\"DateEnd\":false,\"Location\":\"地点\",\"Description\":\"\",\"URL\":\"\",\"Summary\":\"内容\",\"Attendees\":[{\"Name\":\"\",\"Email\":\"googollee@gmail.com\",\"PartStat\":\"NEEDS-ACTION\"},{\"Name\":\"\",\"Email\":\"googollee/twitter@exfe.com\",\"PartStat\":\"NEEDS-ACTION\"}]}]}"

		reader := bytes.NewBufferString(outlook)
		calendar, err := ParseCalendar(reader)
		assert.Equal(t, err, nil)
		j, _ := json.Marshal(calendar)
		assert.Equal(t, string(j), expect)
	}

	{
		outlook_date := `BEGIN:VCALENDAR
METHOD:REQUEST
VERSION:2.0
PRODID:-//Microsoft Corporation//Windows Live Calendar//EN
BEGIN:VTIMEZONE
TZID:China Standard Time
BEGIN:STANDARD
DTSTART:20080101T000000
TZOFFSETTO:+0800
TZOFFSETFROM:+0800
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
UID:d1521bd3-03f0-49ef-a881-731b7067452a
DTSTAMP:20130313T035317Z
CLASS:PUBLIC
X-MICROSOFT-CDO-BUSYSTATUS:FREE
TRANSP:TRANSPARENT
SEQUENCE:1
DTSTART;VALUE=DATE:20130312
DTEND;VALUE=DATE:20130313
SUMMARY:内容
LOCATION:地点
PRIORITY:0
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=
 TRUE:MAILTO:googollee/twitter@exfe.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION;RSVP=
 TRUE:MAILTO:googollee@gmail.com
ORGANIZER;CN=Lee Googol Lee:MAILTO:googollee@hotmail.com
BEGIN:VALARM
ACTION:DISPLAY
TRIGGER:-PT15M
END:VALARM
END:VEVENT
END:VCALENDAR`
		expect := "{\"Event\":[{\"ID\":\"d1521bd3-03f0-49ef-a881-731b7067452a\",\"Organizer\":{\"Name\":\"Lee Googol Lee\",\"Email\":\"googollee@hotmail.com\",\"PartStat\":\"\"},\"Start\":\"2013-03-12T00:00:00Z\",\"DateStart\":true,\"End\":\"2013-03-13T00:00:00Z\",\"DateEnd\":true,\"Location\":\"地点\",\"Description\":\"\",\"URL\":\"\",\"Summary\":\"内容\",\"Attendees\":[{\"Name\":\"\",\"Email\":\"googollee/twitter@exfe.com\",\"PartStat\":\"NEEDS-ACTION\"},{\"Name\":\"\",\"Email\":\"googollee@gmail.com\",\"PartStat\":\"NEEDS-ACTION\"}]}]}"

		reader := bytes.NewBufferString(outlook_date)
		calendar, err := ParseCalendar(reader)
		assert.Equal(t, err, nil)
		j, _ := json.Marshal(calendar)
		assert.Equal(t, string(j), expect)
	}
}
