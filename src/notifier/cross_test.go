package notifier

import (
	"formatter"
	"github.com/stretchrcom/testify/assert"
	"model"
	"testing"
)

func init() {
	var err error
	localTemplate, err = formatter.NewLocalTemplate("../../templates", "en_US")
	if err != nil {
		panic(err)
	}
	config.Email.Name = "business tester"
	config.Email.Domain = "test.com"
}

var localTemplate *formatter.LocalTemplate

var config = model.Config{
	SiteUrl: "http://site/url",
	SiteApi: "http://site/api",
	SiteImg: "http://site/img",
	AppUrl:  "http://app/url",
}

var email1 = model.Identity{
	ID:               11,
	Name:             "email1 name",
	Nickname:         "email1 nick",
	Bio:              "email1 bio",
	Timezone:         "+0800",
	Avatar:           "http://path/to/email1.avatar",
	UserID:           1,
	Provider:         "email",
	ExternalID:       "email1@domain.com",
	ExternalUsername: "email1@domain.com",
}

var remail1 = model.Recipient{
	IdentityID:       11,
	UserID:           1,
	Name:             email1.Name,
	Timezone:         email1.Timezone,
	Token:            "recipient_email1_token",
	Language:         "en_US",
	Provider:         email1.Provider,
	ExternalID:       email1.ExternalID,
	ExternalUsername: email1.ExternalUsername,
}

var email2 = model.Identity{
	ID:               12,
	Name:             "email2 name",
	Nickname:         "email2 nick",
	Bio:              "email2 bio",
	Timezone:         "+0800",
	Avatar:           "http://path/to/email2.avatar",
	UserID:           2,
	Provider:         "email",
	ExternalID:       "email2@domain.com",
	ExternalUsername: "email2@domain.com",
}

var remail2 = model.Recipient{
	IdentityID:       12,
	UserID:           2,
	Name:             email2.Name,
	Timezone:         email2.Timezone,
	Token:            "recipient_email2_token",
	Language:         "en_US",
	Provider:         email2.Provider,
	ExternalID:       email2.ExternalID,
	ExternalUsername: email2.ExternalUsername,
}

var twitter1 = model.Identity{
	ID:               21,
	Name:             "twitter1 name",
	Nickname:         "twitter1 nick",
	Bio:              "twitter1 bio",
	Timezone:         "+0800",
	Avatar:           "http://path/to/twitter1.avatar",
	UserID:           1,
	Provider:         "twitter",
	ExternalID:       "twitter1@domain.com",
	ExternalUsername: "twitter1@domain.com",
}

var rtwitter1 = model.Recipient{
	IdentityID:       21,
	UserID:           1,
	Name:             twitter1.Name,
	Timezone:         twitter1.Timezone,
	Token:            "recipient_twitter1_token",
	Language:         "en_US",
	Provider:         twitter1.Provider,
	ExternalID:       twitter1.ExternalID,
	ExternalUsername: twitter1.ExternalUsername,
}

var twitter3 = model.Identity{
	ID:               22,
	Name:             "twitter3 name",
	Nickname:         "twitter3 nick",
	Bio:              "twitter3 bio",
	Timezone:         "+0800",
	Avatar:           "http://path/to/twitter3.avatar",
	UserID:           3,
	Provider:         "twitter",
	ExternalID:       "twitter3@domain.com",
	ExternalUsername: "twitter3@domain.com",
}

var rtwitter3 = model.Recipient{
	IdentityID:       22,
	UserID:           3,
	Name:             twitter3.Name,
	Timezone:         twitter3.Timezone,
	Token:            "recipient_twitter3_token",
	Language:         "en_US",
	Provider:         twitter3.Provider,
	ExternalID:       twitter3.ExternalID,
	ExternalUsername: twitter3.ExternalUsername,
}

var facebook1 = model.Identity{
	ID:               31,
	Name:             "facebook1 name",
	Nickname:         "facebook1 nick",
	Bio:              "facebook1 bio",
	Timezone:         "+0800",
	Avatar:           "http://path/to/facebook1.avatar",
	UserID:           1,
	Provider:         "facebook",
	ExternalID:       "facebook1@domain.com",
	ExternalUsername: "facebook1@domain.com",
}

var rfacebbok1 = model.Recipient{
	IdentityID:       31,
	UserID:           1,
	Name:             facebook1.Name,
	Timezone:         facebook1.Timezone,
	Token:            "recipient_facebook1_token",
	Language:         "en_US",
	Provider:         facebook1.Provider,
	ExternalID:       facebook1.ExternalID,
	ExternalUsername: facebook1.ExternalUsername,
}

var facebook4 = model.Identity{
	ID:               32,
	Name:             "facebook4 name",
	Nickname:         "facebook4 nick",
	Bio:              "facebook4 bio",
	Timezone:         "+0800",
	Avatar:           "http://path/to/facebook4.avatar",
	UserID:           4,
	Provider:         "facebook",
	ExternalID:       "facebook4@domain.com",
	ExternalUsername: "facebook4@domain.com",
}

var facebook5 = model.Identity{
	ID:               33,
	Name:             "facebook5 name",
	Nickname:         "facebook5 nick",
	Bio:              "facebook5 bio",
	Timezone:         "+0800",
	Avatar:           "http://path/to/facebook5.avatar",
	UserID:           5,
	Provider:         "facebook",
	ExternalID:       "facebook5@domain.com",
	ExternalUsername: "facebook5@domain.com",
}

var facebook6 = model.Identity{
	ID:               34,
	Name:             "facebook6 name",
	Nickname:         "facebook6 nick",
	Bio:              "facebook6 bio",
	Timezone:         "+0800",
	Avatar:           "http://path/to/facebook6.avatar",
	UserID:           6,
	Provider:         "facebook",
	ExternalID:       "facebook6@domain.com",
	ExternalUsername: "facebook6@domain.com",
}

var rfacebook4 = model.Recipient{
	IdentityID:       32,
	UserID:           4,
	Name:             facebook4.Name,
	Timezone:         facebook4.Timezone,
	Token:            "recipient_facebook4_token",
	Language:         "en_US",
	Provider:         facebook4.Provider,
	ExternalID:       facebook4.ExternalID,
	ExternalUsername: facebook4.ExternalUsername,
}

var exfee1 = model.Exfee{
	ID: 123,
	Invitations: []model.Invitation{
		model.Invitation{
			ID:         11,
			Host:       true,
			Mates:      2,
			Identity:   email1,
			RsvpStatus: model.RsvpNoresponse,
			By:         email1,
		},
		model.Invitation{
			ID:         22,
			Identity:   email2,
			RsvpStatus: model.RsvpNoresponse,
			By:         email1,
		},
		model.Invitation{
			ID:         33,
			Identity:   twitter3,
			RsvpStatus: model.RsvpNoresponse,
			By:         email1,
		},
		model.Invitation{
			ID:         44,
			Identity:   facebook4,
			RsvpStatus: model.RsvpNoresponse,
			By:         twitter3,
		},
	},
}

var exfee2 = model.Exfee{
	ID: 123,
	Invitations: []model.Invitation{
		model.Invitation{
			ID:         11,
			Host:       true,
			Mates:      2,
			Identity:   email1,
			RsvpStatus: model.RsvpNoresponse,
			By:         email1,
		},
		model.Invitation{
			ID:         22,
			Identity:   email2,
			RsvpStatus: model.RsvpNoresponse,
			By:         email1,
		},
		model.Invitation{
			ID:         33,
			Identity:   twitter3,
			RsvpStatus: model.RsvpNoresponse,
			By:         email1,
		},
		model.Invitation{
			ID:         44,
			Identity:   facebook4,
			RsvpStatus: model.RsvpNoresponse,
			By:         twitter3,
		},
	},
}

var time1 = model.CrossTime{
	Origin:       "2012-10-23 16:45:00",
	OutputFormat: model.TimeFormat,
	BeginAt: model.EFTime{
		Date:     "2012-10-23",
		Time:     "08:45:00",
		Timezone: "+0800",
	},
}

var time2 = model.CrossTime{
	Origin:       "2012-10-23 16:45:00",
	OutputFormat: model.TimeFormat,
	BeginAt: model.EFTime{
		Date:     "2012-10-23",
		Time:     "16:45:00",
		Timezone: "+0000",
	},
}

var place1 = model.Place{
	Title:       "Test Place1",
	Description: "test place 1",
}

var place2 = model.Place{
	Title:       "Test Place2",
	Description: "test place 2",
}

var cross = model.Cross{
	ID:          123,
	By:          email1,
	Title:       "Test Cross",
	Description: "test cross description",
}

func TestCrossSummaryTwitter(t *testing.T) {
	cross1 := cross
	cross1.Exfee = exfee1
	cross1.Exfee.Invitations = append(cross1.Exfee.Invitations[1:], model.Invitation{
		ID:         55,
		Host:       true,
		Mates:      2,
		Identity:   twitter1,
		RsvpStatus: model.RsvpNoresponse,
		By:         email1,
	})
	cross1.Exfee.Invitations = append(cross1.Exfee.Invitations, model.Invitation{
		ID:         66,
		Mates:      2,
		Identity:   facebook5,
		RsvpStatus: model.RsvpAccepted,
		By:         facebook4,
	})
	cross1.Exfee.Invitations[0].RsvpStatus = model.RsvpAccepted
	cross1.Exfee.Invitations[1].RsvpStatus = model.RsvpDeclined
	cross1.Exfee.Invitations[2].RsvpStatus = model.RsvpAccepted

	cross2 := cross
	cross2.Exfee = exfee2
	cross2.Exfee.Invitations = append(cross2.Exfee.Invitations, model.Invitation{
		ID:         77,
		Identity:   facebook6,
		RsvpStatus: model.RsvpNoresponse,
		By:         facebook4,
	})
	cross2.Exfee.Invitations[3].RsvpStatus = model.RsvpAccepted

	updates := []model.CrossUpdate{
		model.CrossUpdate{
			To:       rtwitter1,
			OldCross: cross2,
			Cross:    cross1,
			By:       facebook4,
		},
	}

	cross2 = cross1
	cross1.Time = time1
	updates = append(updates, model.CrossUpdate{
		To:       rtwitter1,
		OldCross: cross2,
		Cross:    cross1,
		By:       email1,
	})

	cross2 = cross1
	cross1.Place = place1
	updates = append(updates, model.CrossUpdate{
		To:       rtwitter1,
		OldCross: cross2,
		Cross:    cross1,
		By:       email2,
	})

	cross2 = cross1
	cross1.Title = "New Title"
	updates = append(updates, model.CrossUpdate{
		To:       rtwitter1,
		OldCross: cross2,
		Cross:    cross1,
		By:       twitter3,
	})

	c := NewCross(localTemplate, &config, nil)
	private, public, err := c.getSummaryContent(updates)
	assert.Equal(t, err, nil)
	t.Logf("private:-----start------\n%s\n-------end-------", private)
	expectPrivate := "\n\n\n\n\n\n\n\n\\(\"Test Cross\"\\) update: \\(\"New Title\"\\). 4:45PM on Tue, Oct 23 at \\(Test Place1\\). 5 people invited. http://site/url/#!token=recipient_twitter1_token\n\n\\(facebook5 name\\) is invited to \\(\"New Title\"\\) by facebook4 name, email1 name, etc. http://site/url/#!token=recipient_twitter1_token\n\n\\(facebook6 name\\) left \\(\"New Title\"\\). http://site/url/#!token=recipient_twitter1_token\n\n\n\n\n\n\n\n\\(email2 name\\) and \\(facebook5 name\\) accepted \\(\"New Title\"\\), \\(twitter3 name\\) is unavailable, 5 of 9 accepted. http://site/url/#!token=recipient_twitter1_token\n\n\n\n\n"
	assert.Equal(t, private, expectPrivate)
	t.Logf("private:-----start------\n%s\n-------end-------", private)
	expectPublic := `Updates: http://site/url/#!123/eci (Please follow @EXFE to receive details PRIVATELY through Direct Message.)`
	assert.Equal(t, public, expectPublic)
}

func TestCrossSummaryEmail(t *testing.T) {
	cross1 := cross
	cross1.Exfee = exfee1
	cross1.Exfee.Invitations = append(cross1.Exfee.Invitations[1:], model.Invitation{
		ID:         55,
		Host:       true,
		Mates:      2,
		Identity:   twitter1,
		RsvpStatus: model.RsvpNoresponse,
		By:         email1,
	})
	cross1.Exfee.Invitations = append(cross1.Exfee.Invitations, model.Invitation{
		ID:         66,
		Mates:      2,
		Identity:   facebook5,
		RsvpStatus: model.RsvpAccepted,
		By:         facebook4,
	})
	cross1.Exfee.Invitations[0].RsvpStatus = model.RsvpAccepted
	cross1.Exfee.Invitations[1].RsvpStatus = model.RsvpDeclined
	cross1.Exfee.Invitations[2].RsvpStatus = model.RsvpAccepted

	cross2 := cross
	cross2.Exfee = exfee2
	cross2.Exfee.Invitations = append(cross2.Exfee.Invitations, model.Invitation{
		ID:         77,
		Identity:   facebook6,
		RsvpStatus: model.RsvpNoresponse,
		By:         facebook4,
	})
	cross2.Exfee.Invitations[3].RsvpStatus = model.RsvpAccepted

	updates := []model.CrossUpdate{
		model.CrossUpdate{
			To:       remail1,
			OldCross: cross2,
			Cross:    cross1,
			By:       facebook4,
		},
	}

	cross2 = cross1
	cross1.Time = time1
	updates = append(updates, model.CrossUpdate{
		To:       remail1,
		OldCross: cross2,
		Cross:    cross1,
		By:       email1,
	})

	cross2 = cross1
	cross1.Place = place1
	updates = append(updates, model.CrossUpdate{
		To:       remail1,
		OldCross: cross2,
		Cross:    cross1,
		By:       email2,
	})

	cross2 = cross1
	cross1.Title = "New Title"
	updates = append(updates, model.CrossUpdate{
		To:       remail1,
		OldCross: cross2,
		Cross:    cross1,
		By:       twitter3,
	})

	c := NewCross(localTemplate, &config, nil)
	private, public, err := c.getSummaryContent(updates)
	assert.Equal(t, err, nil)
	t.Logf("private:-----start------\n%s\n-------end-------", private)
	expectPrivate := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <x+123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <to_email_address>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x+123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nVXBkYXRlcyBvZiDCt1jCtyDigJxUZXN0IENyb3Nz4oCdIGJ5IGZhY2Vib29rNCBuYW1lLCBlbWFpbDEg\r\nbmFtZSwgZW1haWwyIG5hbWUsIGV0Yy4KCipOZXcgVGl0bGUqCj09PT09PT0KaHR0cDovL3NpdGUvdXJs\r\nLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbgoKKjQ6NDVQTSBvbiBUdWUsIE9jdCAyMyoKPT09\r\nPT09PQoKKlRlc3QgUGxhY2UxKgo9PT09PT09CiAgKnRlc3QgcGxhY2UgMSoKCgrCtyA1IEFjY2VwdGVk\r\nOiBlbWFpbDIgbmFtZSwgZmFjZWJvb2s1IG5hbWUgYW5kIDEgb3RoZXJzLgrCtyBVbmF2YWlsYWJsZTog\r\ndHdpdHRlcjMgbmFtZS4KwrcgTmV3bHkgaW52aXRlZDogZmFjZWJvb2s1IG5hbWUuCsK3IFJlbW92ZWQ6\r\nIGZhY2Vib29rNiBuYW1lLgoKIyBSZXBseSB0aGlzIGVtYWlsIGRpcmVjdGx5IGFzIGNvbnZlcnNhdGlv\r\nbi4gIw==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9pZGVudGl0\r\neV9uYW1lIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfbXNnX2lkZW50aXR5X25hbWUgewogICAgICAgICAgICAgICAgY29sb3I6\r\nICM2NjY2NjY7CiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfYXQgewogICAg\r\nICAgICAgICAgICAgZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICAgICAgY29sb3I6ICM5OTk5OTk7\r\nCiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfdGltZSB7CiAgICAgICAgICAg\r\nICAgICBmb250LXNpemU6IDEycHg7CiAgICAgICAgICAgICAgICBjb2xvcjogIzY2NjY2NjsKICAgICAg\r\nICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4KICAgICAgICA8dGFi\r\nbGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAiIHN0eWxlPSJmb250LWZh\r\nbWlseTogVmVyZGFuYTsgZm9udC1zaXplOiAxM3B4OyBsaW5lLWhlaWdodDogMjBweDsgY29sb3I6ICMx\r\nOTE5MTk7IGZvbnQtd2VpZ2h0OiBub3JtYWw7IHdpZHRoOiA2NDBweDsgcGFkZGluZzogMjBweDsgYmFj\r\na2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsiPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8\r\ndGQgY29sc3Bhbj0iNSIgc3R5bGU9ImNvbG9yOiAjMzMzMzMzOyI+CiAgICAgICAgICAgICAgICAgICAg\r\nPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5\r\nbGU9ImNvbG9yOiAjMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5VcGRhdGVzIG9mIDxzcGFu\r\nIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsiPsK3WMK3PC9zcGFuPiDigJw8c3BhbiBzdHlsZT0iY29sb3I6\r\nICMxOTE5MTk7Ij5UZXN0IENyb3NzPC9zcGFuPuKAnSBieSBmYWNlYm9vazQgbmFtZSwgZW1haWwxIG5h\r\nbWUsIGVtYWlsMiBuYW1lLCBldGMuPC9hPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAg\r\nPC90cj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjEwIj48L3RkPjwvdHI+\r\nCiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1IiBzdHlsZT0iZm9u\r\ndC1zaXplOiAyMHB4OyBsaW5lLWhlaWdodDogMjZweDsiPgogICAgICAgICAgICAgICAgICAgIDxhIGhy\r\nZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJj\r\nb2xvcjojM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IGZvbnQtd2VpZ2h0OiBsaWdodGVyOyI+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIE5ldyBUaXRsZQogICAgICAgICAgICAgICAgICAgIDwvYT4K\r\nICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQg\r\nY29sc3Bhbj0iNSIgaGVpZ2h0PSIxMCI+PC90ZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAg\r\nICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSIxODAiPgogICAgICAgICAgICAgICAgICAgIAog\r\nICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJmb250LXNpemU6IDIwcHg7IGxpbmUtaGVpZ2h0OiAy\r\nNnB4OyBtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3Np\r\ndGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1\r\nOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij40OjQ1UE0gb24gVHVlLCBPY3QgMjM8L2E+CiAgICAgICAg\r\nICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4K\r\nICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMTAiPjwvdGQ+CiAgICAgICAgICAgICAgICA8dGQgdmFs\r\naWduPSJ0b3AiIHdpZHRoPSIxOTAiIHN0eWxlPSJ3b3JkLWJyZWFrOiBicmVhay1hbGw7Ij4KICAgICAg\r\nICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0iZm9udC1zaXplOiAyMHB4\r\nOyBsaW5lLWhlaWdodDogMjZweDsgbWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgIDxh\r\nIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxl\r\nPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+VGVzdCBQbGFjZTE8L2E+CiAg\r\nICAgICAgICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJtYXJnaW46\r\nIDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9r\r\nZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29y\r\nYXRpb246IG5vbmU7Ij50ZXN0IHBsYWNlIDE8L2E+CiAgICAgICAgICAgICAgICAgICAgPC9wPgogICAg\r\nICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgIDx0ZCB3\r\naWR0aD0iMTAiPjwvdGQ+CiAgICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSIyMTAi\r\nPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90\r\ncj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjEwIj48L3RkPjwvdHI+CiAg\r\nICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1Ij4KICAgICAgICAgICAg\r\nICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAiIHN0\r\neWxlPSJmb250LWZhbWlseTogVmVyZGFuYTsgZm9udC1zaXplOiAxM3B4OyBsaW5lLWhlaWdodDogMjBw\r\neDsgY29sb3I6ICMxOTE5MTk7IGZvbnQtd2VpZ2h0OiBub3JtYWw7IHdpZHRoOiAxMDAlOyBiYWNrZ3Jv\r\ndW5kLWNvbG9yOiAjZmJmYmZiOyI+CiAgICAgICAgICAgICAgICAgICAgCQogICAgICAgICAgICAgICAg\r\nICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQgd2lkdGg9IjE1Ij48aW1n\r\nIHNyYz0iaHR0cDovL3NpdGUvaW1nL2VtYWlsL3JzdnBfYWNjZXB0ZWRfMTJfYmx1ZS5wbmciIC8+PC90\r\nZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rv\r\na2VuIiBzdHlsZT0iY29sb3I6ICMxOTE5MTk7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPgogICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0ibWFyZ2luOiAwOyI+CiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8c3BhbiBjbGFzcz0iZXhmZV9tYWlsX2lkZW50\r\naXR5X25hbWUiPjU8L3NwYW4+IEFjY2VwdGVkOiA8c3BhbiBjbGFzcz0nZXhmZV9tYWlsX2lkZW50aXR5\r\nX25hbWUnPmVtYWlsMiBuYW1lPC9zcGFuPiwgPHNwYW4gY2xhc3M9J2V4ZmVfbWFpbF9pZGVudGl0eV9u\r\nYW1lJz5mYWNlYm9vazUgbmFtZTwvc3Bhbj4sIGFuZCAxIG90aGVycy4KICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvYT4K\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgIDwv\r\ndHI+CiAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAg\r\nICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRo\r\nPSIxNSI+PGltZyBzcmM9Imh0dHA6Ly9zaXRlL2ltZy9lbWFpbC9yc3ZwX2RlY2xpbmVkXzEyLnBuZyIg\r\nLz48L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgIDxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFp\r\nbDFfdG9rZW4iIHN0eWxlPSJjb2xvcjogIzE5MTkxOTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJtYXJnaW46IDA7Ij4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIFVuYXZhaWxhYmxlOiB0d2l0dGVyMyBu\r\nYW1lCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvcD4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8L2E+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAg\r\nICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAg\r\nICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMTUiPjxpbWcgc3JjPSJodHRwOi8vc2l0ZS9pbWcvZW1haWwv\r\ncGx1c18xMl9ibHVlLnBuZyIgLz48L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRv\r\na2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJjb2xvcjogIzE5MTkxOTsgdGV4dC1kZWNv\r\ncmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0eWxl\r\nPSJtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIE5ld2x5\r\nIGludml0ZWQ6IDxzcGFuIGNsYXNzPSJleGZlX21haWxfaWRlbnRpdHlfbmFtZSI+ZmFjZWJvb2s1IG5h\r\nbWU8L3NwYW4+LgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3A+CiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90\r\nZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAgCiAg\r\nICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8dGQgd2lkdGg9IjE1Ij48aW1nIHNyYz0iaHR0cDovL3NpdGUvaW1n\r\nL2VtYWlsL21pbnVzXzEyLnBuZyIgLz48L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRk\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8j\r\nIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJjb2xvcjogIzE5MTkxOTsgdGV4dC1k\r\nZWNvcmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0\r\neWxlPSJtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIFJl\r\nbW92ZWQ6IGZhY2Vib29rNiBuYW1lLgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8\r\nL3A+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgPC90YWJsZT4KICAgICAgICAgICAgICAgIDwvdGQ+\r\nCiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQgY29sc3Bhbj0iNSIgaGVpZ2h0PSIx\r\nMCI+PC90ZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQgY29sc3Bhbj0i\r\nNSI+CiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAg\r\nICAgICAgPGltZyBzdHlsZT0icGFkZGluZy1yaWdodDogNXB4OyIgd2lkdGg9IjQwIiBoZWlnaHQ9IjQw\r\nIiBhbHQ9ImVtYWlsMiBuYW1lIiB0aXRsZT0iZW1haWwyIG5hbWUiIHNyYz0iaHR0cDovL3NpdGUvYXBp\r\nL3YyL2F2YXRhci9yZW5kZXI/cmVzb2x1dGlvbj0yeCZ1cmw9YUhSMGNEb3ZMM0JoZEdndmRHOHZaVzFo\r\nYVd3eUxtRjJZWFJoY2clM0QlM0Qmd2lkdGg9NDAmaGVpZ2h0PTQwIj4KICAgICAgICAgICAgICAgICAg\r\nICAKICAgICAgICAgICAgICAgICAgICA8aW1nIHN0eWxlPSJwYWRkaW5nLXJpZ2h0OiA1cHg7IiB3aWR0\r\naD0iNDAiIGhlaWdodD0iNDAiIGFsdD0idHdpdHRlcjMgbmFtZSIgdGl0bGU9InR3aXR0ZXIzIG5hbWUi\r\nIHNyYz0iaHR0cDovL3NpdGUvYXBpL3YyL2F2YXRhci9yZW5kZXI/cmVzb2x1dGlvbj0yeCZ1cmw9YUhS\r\nMGNEb3ZMM0JoZEdndmRHOHZkSGRwZEhSbGNqTXVZWFpoZEdGeSZ3aWR0aD00MCZoZWlnaHQ9NDAmYWxw\r\naGE9MC4zMyI+CiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgPGltZyBzdHls\r\nZT0icGFkZGluZy1yaWdodDogNXB4OyIgd2lkdGg9IjQwIiBoZWlnaHQ9IjQwIiBhbHQ9ImZhY2Vib29r\r\nNCBuYW1lIiB0aXRsZT0iZmFjZWJvb2s0IG5hbWUiIHNyYz0iaHR0cDovL3NpdGUvYXBpL3YyL2F2YXRh\r\nci9yZW5kZXI/cmVzb2x1dGlvbj0yeCZ1cmw9YUhSMGNEb3ZMM0JoZEdndmRHOHZabUZqWldKdmIyczBM\r\nbUYyWVhSaGNnJTNEJTNEJndpZHRoPTQwJmhlaWdodD00MCI+CiAgICAgICAgICAgICAgICAgICAgCiAg\r\nICAgICAgICAgICAgICAgICAgPGltZyBzdHlsZT0icGFkZGluZy1yaWdodDogNXB4OyIgd2lkdGg9IjQw\r\nIiBoZWlnaHQ9IjQwIiBhbHQ9InR3aXR0ZXIxIG5hbWUiIHRpdGxlPSJ0d2l0dGVyMSBuYW1lIiBzcmM9\r\nImh0dHA6Ly9zaXRlL2FwaS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249MngmdXJsPWFIUjBjRG92\r\nTDNCaGRHZ3ZkRzh2ZEhkcGRIUmxjakV1WVhaaGRHRnkmd2lkdGg9NDAmaGVpZ2h0PTQwJmFscGhhPTAu\r\nMzMmaXNob3N0PXRydWUmbWF0ZXM9MiI+CiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAg\r\nICAgICAgPGltZyBzdHlsZT0icGFkZGluZy1yaWdodDogNXB4OyIgd2lkdGg9IjQwIiBoZWlnaHQ9IjQw\r\nIiBhbHQ9ImZhY2Vib29rNSBuYW1lIiB0aXRsZT0iZmFjZWJvb2s1IG5hbWUiIHNyYz0iaHR0cDovL3Np\r\ndGUvYXBpL3YyL2F2YXRhci9yZW5kZXI/cmVzb2x1dGlvbj0yeCZ1cmw9YUhSMGNEb3ZMM0JoZEdndmRH\r\nOHZabUZqWldKdmIyczFMbUYyWVhSaGNnJTNEJTNEJndpZHRoPTQwJmhlaWdodD00MCZtYXRlcz0yIj4K\r\nICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+\r\nCiAgICAgICAgICAgIDx0cj48dGQgY29sc3Bhbj0iNSIgaGVpZ2h0PSIxMCI+PC90ZD48L3RyPgogICAg\r\nICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQgY29sc3Bhbj0iNSI+CiAgICAgICAgICAgICAg\r\nICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tl\r\nbiIgc3R5bGU9ImNvbG9yOiAjMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij50ZXN0IGNyb3Nz\r\nIGRlc2NyaXB0aW9uPC9hPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4KICAg\r\nICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjIwIj48L3RkPjwvdHI+CiAgICAgICAg\r\nICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1IiBzdHlsZT0iZm9udC1zaXplOiAx\r\nMXB4OyBsaW5lLWhlaWdodDogMTVweDsgY29sb3I6ICM3RjdGN0Y7Ij4KICAgICAgICAgICAgICAgICAg\r\nICBSZXBseSB0aGlzIGVtYWlsIGRpcmVjdGx5IGFzIGNvbnZlcnNhdGlvbiwgb3IgdHJ5IDxhIHN0eWxl\r\nPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJlZj0iaHR0cDovL2FwcC91\r\ncmwiPkVYRkU8L2E+IGFwcC4KICAgICAgICAgICAgICAgICAgICA8YnIgLz4KICAgICAgICAgICAgICAg\r\nICAgICA8c3BhbiBzdHlsZT0iY29sb3I6ICNCMkIyQjIiPlRoaXMgdXBkYXRlIGlzIHNlbnQgZnJvbSA8\r\nYSBzdHlsZT0iY29sb3I6ICMzYTZlYTU7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiIGhyZWY9Imh0dHA6\r\nLy9zaXRlL3VybCI+RVhGRTwvYT4gYXV0b21hdGljYWxseS4gPGEgc3R5bGU9ImNvbG9yOiAjRTZFNkU2\r\nOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8vc2l0ZS91cmwvbXV0ZS9jcm9zcz90\r\nb2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIj5VbnN1YnNjcmliZT88L2E+CiAgICAgICAgICAgICAg\r\nICAgICAgPCEtLQogICAgICAgICAgICAgICAgICAgIFlvdSBjYW4gY2hhbmdlCiAgICAgICAgICAgICAg\r\nICAgICAgPGEgc3R5bGU9ImNvbG9yOiAjQjJCMkIyOyB0ZXh0LWRlY29yYXRpb246IHVuZGVsaW5lOyIg\r\naHJlZj0iIj5ub3RpZmljYXRpb24gcHJlZmVyZW5jZTwvYT4uCiAgICAgICAgICAgICAgICAgICAgLS0+\r\nCiAgICAgICAgICAgICAgICAgICAgPC9zcGFuPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAg\r\nICAgPC90cj4KICAgICAgICA8L3RhYmxlPgogICAgPC9ib2R5Pgo8L2h0bWw+Cg==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6TmV3IFRpdGxlDQpYLVdSLUNBTERFU0M6ZXhmZSBjYWwNClgtV1ItVElNRVpPTkU6\r\nKzA4MDANCkJFR0lOOlZFVkVOVA0KVUlEOiExMjNAZXhmZQ0KRFRTVEFNUDoyMDEyMTAyM1QwODQ1MDBa\r\nDQpERVNDUklQVElPTjp0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uDQpEVFNUQVJUOjIwMTIxMDIzVDA4NDUw\r\nMFoNCkxPQ0FUSU9OOlRlc3QgUGxhY2UxXG50ZXN0IHBsYWNlIDENClNVTU1BUlk6TmV3IFRpdGxlDQpV\r\nUkw6aHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbg0KRU5EOlZFVkVO\r\nVA0KRU5EOlZDQUxFTkRBUg0K\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
	assert.Equal(t, private, expectPrivate)
	assert.Equal(t, public, "")
}

func TestCrossInvitationEmail(t *testing.T) {
	cross1 := cross
	cross1.Time = time1
	cross1.Exfee = exfee1

	inv := model.CrossInvitation{}
	inv.To = remail1
	inv.Cross = cross1

	c := NewCross(localTemplate, &config, nil)
	private, public, err := c.getInvitationContent(inv)
	assert.Equal(t, err, nil)
	t.Logf("private:---------start---------\n%s\n---------end----------", private)
	expectPrivate := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <x+123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <to_email_address>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x+123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nWW91J3JlIGdhdGhlcmluZyB0aGlzIMK3WMK3LgoKClRlc3QgQ3Jvc3MKPT09PT09PQpodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuCgo0OjQ1UE0gb24gVHVlLCBPY3QgMjMK\r\nPT09PT09PQoKUGxhY2UKPT09PT09PQogIFRvIGJlIGRlY2lkZWQuCgoKSSdtIGluLiBDaGVjayBpdCBv\r\ndXQ6IGh0dHA6Ly9zaXRlL3VybC8/dG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiZyc3ZwPWFjY2Vw\r\ndAoKCjYgSW52aXRlZDoKwrcgZW1haWwxIG5hbWUgKEhvc3QpIHdpdGggMiBwZW9wbGUKwrcgZW1haWwy\r\nIG5hbWUKwrcgdHdpdHRlcjMgbmFtZQrCtyBmYWNlYm9vazQgbmFtZQoKCkRlc2NyaXB0aW9uCi0tLS0t\r\nLS0KICB0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uCgoKIyBSZXBseSB0aGlzIGVtYWlsIGRpcmVjdGx5IGFz\r\nIGNvbnZlcnNhdGlvbi4gIw==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9sYWJlbCB7\r\nCiAgICAgICAgICAgICAgICBiYWNrZ3JvdW5kLWNvbG9yOiAjRDVFOEYyOwogICAgICAgICAgICAgICAg\r\nY29sb3I6ICMzYTZlYTU7CiAgICAgICAgICAgICAgICBmb250LXNpemU6IDExcHg7CiAgICAgICAgICAg\r\nICAgICBwYWRkaW5nOiAwIDJweCAwIDJweDsKICAgICAgICAgICAgfQogICAgICAgICAgICAuZXhmZV9t\r\nYWlsX21hdGVzIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICAgICAg\r\nZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICB9CiAgICAgICAgICAgIC5leGZlX21haWxfaWRlbnRp\r\ndHkgewogICAgICAgICAgICAgICAgZm9udC1zdHlsZTogaXRhbGljOwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfaWRlbnRpdHlfbmFtZSB7CiAgICAgICAgICAgICAgICBjb2xvcjogIzE5\r\nMTkxOTsKICAgICAgICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4K\r\nICAgICAgICA8dGFibGUgd2lkdGg9IjY0MCIgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNw\r\nYWNpbmc9IjAiIHN0eWxlPSJmb250LWZhbWlseTogSGVsdmV0aWNhOyBmb250LXNpemU6IDEzcHg7IGxp\r\nbmUtaGVpZ2h0OiAxOXB4OyBjb2xvcjogIzE5MTkxOTsgZm9udC13ZWlnaHQ6IG5vcm1hbDsgcGFkZGlu\r\nZzogMzBweCA0MHB4IDMwcHggNDBweDsgYmFja2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsgbWluLWhlaWdo\r\ndDogNTYycHg7Ij4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjMi\r\nIHZhbGlnbj0idG9wIiBzdHlsZT0iZm9udC1zaXplOiAzMnB4OyBsaW5lLWhlaWdodDogMzhweDsgcGFk\r\nZGluZy1ib3R0b206IDE4cHg7Ij4KICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMzYTZlYTU7\r\nIHRleHQtZGVjb3JhdGlvbjogbm9uZTsgZm9udC13ZWlnaHQ6IDMwMDsiPgogICAgICAgICAgICAgICAg\r\nICAgICAgICBUZXN0IENyb3NzCiAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAg\r\nPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRk\r\nIHdpZHRoPSIzNDAiIHN0eWxlPSJ2ZXJ0aWNhbC1hbGlnbjogYmFzZWxpbmU7IGZvbnQtd2VpZ2h0OiAz\r\nMDA7Ij4KICAgICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIg\r\nY2VsbHNwYWNpbmc9IjAiPgogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHN0eWxlPSJwYWRkaW5nLWJvdHRvbTogMjBweDsg\r\nZm9udC1zaXplOiAyMHB4OyB2ZXJ0aWNhbC1hbGlnbjogYmFzZWxpbmU7Ij4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBZb3UncmUgZ2F0\r\naGVyaW5nIHRoaXMgPHNwYW4gc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5v\r\nbmU7Ij7Ct1jCtzwvc3Bhbj4uCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAg\r\nICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2VsbHBhZGRpbmc9\r\nIjAiIGNlbGxzcGFjaW5nPSIwIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRy\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIiB3\r\naWR0aD0iMTYwIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBo\r\ncmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0i\r\ndGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgCQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9ImZv\r\nbnQtc2l6ZTogMjBweDsgbGluZS1oZWlnaHQ6IDI2cHg7IG1hcmdpbjogMDsgY29sb3I6ICMzMzMzMzM7\r\nIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDQ6NDVQ\r\nTSBvbiBUdWUsIE9jdCAyMwogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICA8L3A+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvYT4KICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHN0eWxlPSJwYWRkaW5nLWxlZnQ6IDEwcHg7Ij4K\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8v\r\nc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0idGV4dC1kZWNvcmF0\r\naW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0i\r\nZm9udC1zaXplOiAyMHB4OyBsaW5lLWhlaWdodDogMjZweDsgbWFyZ2luOiAwOyBjb2xvcjogIzMzMzMz\r\nMzsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgUGxh\r\nY2UKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9wPgogICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0iY29sb3I6\r\nICMxOTE5MTk7IG1hcmdpbjogMDsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgVG8gYmUgZGVjaWRlZAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8L3A+IAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L2E+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nIDwvdGFibGU+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAg\r\nICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAg\r\nICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIiBzdHlsZT0i\r\ncGFkZGluZy10b3A6IDMwcHg7IHBhZGRpbmctYm90dG9tOiAzMHB4OyI+CiAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgPGEgc3R5bGU9ImZsb2F0OiBsZWZ0OyBkaXNwbGF5OiBibG9jazsgdGV4dC1k\r\nZWNvcmF0aW9uOiBub25lOyBib3JkZXI6IDFweCBzb2xpZCAjYmViZWJlOyBiYWNrZ3JvdW5kLWNvbG9y\r\nOiAjM0E2RUE1OyBjb2xvcjogI0ZGRkZGRjsgcGFkZGluZzogNXB4IDMwcHggNXB4IDMwcHg7IG1hcmdp\r\nbi1sZWZ0OiAyNXB4OyIgYWx0PSJBY2NlcHQiIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8/dG9rZW49cmVj\r\naXBpZW50X2VtYWlsMV90b2tlbiZyc3ZwPWFjY2VwdCI+SSdtIGluPC9hPgogICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgIDxhIHN0eWxlPSJmbG9hdDogbGVmdDsgZGlzcGxheTogYmxvY2s7IHRleHQt\r\nZGVjb3JhdGlvbjogbm9uZTsgYm9yZGVyOiAxcHggc29saWQgI2JlYmViZTsgYmFja2dyb3VuZC1jb2xv\r\ncjogI0U2RTZFNjsgY29sb3I6ICMxOTE5MTk7IHBhZGRpbmc6IDVweCAyNXB4IDVweCAyNXB4OyBtYXJn\r\naW4tbGVmdDogMTVweDsiIGFsdD0iQ2hlY2sgaXQgb3V0IiBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0\r\nb2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIj5DaGVjayBpdCBvdXQuLi48L2E+CiAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAg\r\nICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICB0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nCiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8\r\nL3RyPgogICAgICAgICAgICAgICAgICAgIDwvdGFibGU+CiAgICAgICAgICAgICAgICA8L3RkPgogICAg\r\nICAgICAgICAgICAgPHRkIHdpZHRoPSIzMCI+PC90ZD4KICAgICAgICAgICAgICAgIDx0ZCB2YWxpZ249\r\nInRvcCI+CiAgICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2VsbHBhZGRpbmc9IjAi\r\nIGNlbGxzcGFjaW5nPSIwIj4KICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQgaGVpZ2h0PSI2OCIg\r\ndmFsaWduPSJ0b3AiIGFsaWduPSJyaWdodCI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3Rk\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0\r\nZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBzdHlsZT0i\r\nY29sb3I6ICMzMzMzMzM7IiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAiPgogICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIy\r\nNSIgaGVpZ2h0PSIyNSIgYWxpZ249ImxlZnQiIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICA8aW1nIHdpZHRoPSIyMCIgaGVpZ2h0PSIyMCIgdGl0bGU9\r\nImVtYWlsMSBuYW1lIiBhbHQ9ImVtYWlsMSBuYW1lIiBzcmM9Imh0dHA6Ly9wYXRoL3RvL2VtYWlsMS5h\r\ndmF0YXIiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgPHNwYW4+ZW1haWwxIG5hbWU8L3NwYW4+IDxzcGFuIGNsYXNzPSJl\r\neGZlX21haWxfbWF0ZXMiPisyPC9zcGFuPiA8c3BhbiBjbGFzcz0iZXhmZV9tYWlsX2xhYmVsIj5ob3N0\r\nPC9zcGFuPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMjUiIGhlaWdodD0iMjUi\r\nIGFsaWduPSJsZWZ0IiB2YWxpZ249InRvcCI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgPGltZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIHRpdGxlPSJlbWFpbDIgbmFtZSIg\r\nYWx0PSJlbWFpbDIgbmFtZSIgc3JjPSJodHRwOi8vcGF0aC90by9lbWFpbDIuYXZhdGFyIj4KICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8dGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDxzcGFuPmVtYWlsMiBuYW1lPC9zcGFuPgogICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4K\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0\r\nZCB3aWR0aD0iMjUiIGhlaWdodD0iMjUiIGFsaWduPSJsZWZ0IiB2YWxpZ249InRvcCI+CiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPGltZyB3aWR0aD0iMjAiIGhlaWdodD0i\r\nMjAiIHRpdGxlPSJ0d2l0dGVyMyBuYW1lIiBhbHQ9InR3aXR0ZXIzIG5hbWUiIHNyYz0iaHR0cDovL3Bh\r\ndGgvdG8vdHdpdHRlcjMuYXZhdGFyIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxzcGFuPnR3aXR0ZXIzIG5hbWU8L3Nw\r\nYW4+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIyNSIgaGVpZ2h0PSIyNSIgYWxp\r\nZ249ImxlZnQiIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICA8aW1nIHdpZHRoPSIyMCIgaGVpZ2h0PSIyMCIgdGl0bGU9ImZhY2Vib29rNCBuYW1lIiBh\r\nbHQ9ImZhY2Vib29rNCBuYW1lIiBzcmM9Imh0dHA6Ly9wYXRoL3RvL2ZhY2Vib29rNC5hdmF0YXIiPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgPHNwYW4+ZmFjZWJvb2s0IG5hbWU8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8L3RhYmxlPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4K\r\nICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAgICA8L3RhYmxlPgog\r\nICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgPHRyPgogICAg\r\nICAgICAgICAgICAgPHRkIGNvbHNwYW49IjMiIHN0eWxlPSJmb250LXNpemU6IDExcHg7IGxpbmUtaGVp\r\nZ2h0OiAxNXB4OyBjb2xvcjogIzdGN0Y3RjsgcGFkZGluZy10b3A6IDQwcHg7Ij4KICAgICAgICAgICAg\r\nICAgICAgICBSZXBseSB0aGlzIGVtYWlsIGRpcmVjdGx5IGFzIGNvbnZlcnNhdGlvbiwgb3IgVHJ5IDxh\r\nIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJlZj0iaHR0cDov\r\nL2FwcC91cmwiPkVYRkU8L2E+IGFwcC4KICAgICAgICAgICAgICAgICAgICA8YnIgLz4KICAgICAgICAg\r\nICAgICAgICAgICBUaGlzIDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBu\r\nb25lOyIgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiI+\r\nwrdYwrc8L2E+IGludml0YXRpb24gaXMgc2VudCBieSA8c3BhbiBjbGFzcz0iZXhmZV9tYWlsX2lkZW50\r\naXR5X25hbWUiPmVtYWlsMSBuYW1lPC9zcGFuPiBmcm9tIDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsg\r\ndGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJlZj0iaHR0cDovL3NpdGUvdXJsIj5FWEZFPC9hPi4KICAg\r\nICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgPC90YWJsZT4KICAgIDwv\r\nYm9keT4KPC9odG1sPgo=\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6VGVzdCBDcm9zcw0KWC1XUi1DQUxERVNDOmV4ZmUgY2FsDQpYLVdSLVRJTUVaT05F\r\nOiswODAwDQpCRUdJTjpWRVZFTlQNClVJRDohMTIzQGV4ZmUNCkRUU1RBTVA6MjAxMjEwMjNUMDg0NTAw\r\nWg0KREVTQ1JJUFRJT046dGVzdCBjcm9zcyBkZXNjcmlwdGlvbg0KRFRTVEFSVDoyMDEyMTAyM1QwODQ1\r\nMDBaDQpMT0NBVElPTjoNClNVTU1BUlk6VGVzdCBDcm9zcw0KVVJMOmh0dHA6Ly9zaXRlL3VybC8jIXRv\r\na2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4NCkVORDpWRVZFTlQNCkVORDpWQ0FMRU5EQVINCg==\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
	assert.Equal(t, private, expectPrivate)
	assert.Equal(t, public, "")
}

func TestCrossInvitationTwitter(t *testing.T) {
	cross1 := cross
	cross1.Time = time1
	cross1.Exfee = exfee1

	inv := model.CrossInvitation{}
	inv.To = rtwitter1
	inv.Cross = cross1

	c := NewCross(localTemplate, &config, nil)
	private, public, err := c.getInvitationContent(inv)
	assert.Equal(t, err, nil)
	t.Logf("private:---------start---------\n%s\n---------end----------", private)
	expectPrivate := "\n\n\n\nSuccessfully gathering \\(\"Test Cross\"\\), \\(4:45PM on Tue, Oct 23\\). 6 invited: email1 name, email2 name, twitter3 name... http://site/url/#!token=recipient_twitter1_token"
	assert.Equal(t, private, expectPrivate)
	assert.Equal(t, public, "Invitation: http://site/url/#!123/eci (Please follow @EXFE to receive details PRIVATELY through Direct Message.)")
}
