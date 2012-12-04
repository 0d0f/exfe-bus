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
	expectPrivate := "\n\n\n\n\n\n\n\n\\(“Test Cross”\\) update: \\(“New Title”\\). 4:45PM on Tue, Oct 23 at \\(Test Place1\\). 5 people invited. http://site/url/#!token=recipient_twitter1_token\n\n\\(facebook5 name\\) is invited to \\(“New Title”\\) by facebook4 name, email1 name, etc. http://site/url/#!token=recipient_twitter1_token\n\n\\(facebook6 name\\) left \\(“New Title”\\). http://site/url/#!token=recipient_twitter1_token\n\n\n\n\n\n\n\n\\(email2 name\\) and \\(facebook5 name\\) accepted \\(“New Title”\\), \\(twitter3 name\\) is unavailable, 5 of 9 accepted. http://site/url/#!token=recipient_twitter1_token\n\n\n\n\n"
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
	expectPrivate := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <x+123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <to_email_address>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x+123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nVXBkYXRlcyBvZiDCt1jCtyDigJxUZXN0IENyb3Nz4oCdIGJ5IGZhY2Vib29rNCBuYW1lLCBlbWFpbDEg\r\nbmFtZSwgZW1haWwyIG5hbWUsIGV0Yy4KCipOZXcgVGl0bGUqCj09PT09PT0KaHR0cDovL3NpdGUvdXJs\r\nLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbgoKKjQ6NDVQTSBvbiBUdWUsIE9jdCAyMyoKPT09\r\nPT09PQoKKlRlc3QgUGxhY2UxKgo9PT09PT09CiAgKnRlc3QgcGxhY2UgMSoKCgrCtyA1IEFjY2VwdGVk\r\nOiBlbWFpbDIgbmFtZSwgZmFjZWJvb2s1IG5hbWUgYW5kIDEgb3RoZXJzLgrCtyBVbmF2YWlsYWJsZTog\r\ndHdpdHRlcjMgbmFtZS4KwrcgTmV3bHkgaW52aXRlZDogZmFjZWJvb2s1IG5hbWUuCsK3IFJlbW92ZWQ6\r\nIGZhY2Vib29rNiBuYW1lLgoKIyBSZXBseSB0aGlzIGVtYWlsIGRpcmVjdGx5IGFzIGNvbnZlcnNhdGlv\r\nbi4gIw==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9pZGVudGl0\r\neV9uYW1lIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfbXNnX2lkZW50aXR5X25hbWUgewogICAgICAgICAgICAgICAgY29sb3I6\r\nICM2NjY2NjY7CiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfYXQgewogICAg\r\nICAgICAgICAgICAgZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICAgICAgY29sb3I6ICM5OTk5OTk7\r\nCiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfdGltZSB7CiAgICAgICAgICAg\r\nICAgICBmb250LXNpemU6IDEycHg7CiAgICAgICAgICAgICAgICBjb2xvcjogIzY2NjY2NjsKICAgICAg\r\nICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4KICAgICAgICA8dGFi\r\nbGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAiIHN0eWxlPSJmb250LWZh\r\nbWlseTogVmVyZGFuYTsgZm9udC1zaXplOiAxM3B4OyBsaW5lLWhlaWdodDogMjBweDsgY29sb3I6ICMx\r\nOTE5MTk7IGZvbnQtd2VpZ2h0OiBub3JtYWw7IHdpZHRoOiA2NDBweDsgcGFkZGluZzogMjBweDsgYmFj\r\na2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsiPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8\r\ndGQgY29sc3Bhbj0iNSIgc3R5bGU9ImNvbG9yOiAjMzMzMzMzOyI+CiAgICAgICAgICAgICAgICAgICAg\r\nPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5\r\nbGU9ImNvbG9yOiAjMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5VcGRhdGVzIG9mIDxzcGFu\r\nIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsiPsK3WMK3PC9zcGFuPiDigJw8c3BhbiBzdHlsZT0iY29sb3I6\r\nICMxOTE5MTk7Ij5UZXN0IENyb3NzPC9zcGFuPuKAnSBieSBmYWNlYm9vazQgbmFtZSwgZW1haWwxIG5h\r\nbWUsIGVtYWlsMiBuYW1lLCBldGMuPC9hPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAg\r\nPC90cj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjEwIj48L3RkPjwvdHI+\r\nCiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1IiBzdHlsZT0iZm9u\r\ndC1zaXplOiAyMHB4OyBsaW5lLWhlaWdodDogMjZweDsiPgogICAgICAgICAgICAgICAgICAgIDxhIGhy\r\nZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJj\r\nb2xvcjojM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IGZvbnQtd2VpZ2h0OiBsaWdodGVyOyI+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIE5ldyBUaXRsZQogICAgICAgICAgICAgICAgICAgIDwvYT4K\r\nICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQg\r\nY29sc3Bhbj0iNSIgaGVpZ2h0PSIxMCI+PC90ZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAg\r\nICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSIxODAiPgogICAgICAgICAgICAgICAgICAgIAog\r\nICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJmb250LXNpemU6IDIwcHg7IGxpbmUtaGVpZ2h0OiAy\r\nNnB4OyBtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3Np\r\ndGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1\r\nOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij40OjQ1UE0gb24gVHVlLCBPY3QgMjM8L2E+CiAgICAgICAg\r\nICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4K\r\nICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMTAiPjwvdGQ+CiAgICAgICAgICAgICAgICA8dGQgdmFs\r\naWduPSJ0b3AiIHdpZHRoPSIxOTAiIHN0eWxlPSJ3b3JkLWJyZWFrOiBicmVhay1hbGw7Ij4KICAgICAg\r\nICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0iZm9udC1zaXplOiAyMHB4\r\nOyBsaW5lLWhlaWdodDogMjZweDsgbWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgIDxh\r\nIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxl\r\nPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+VGVzdCBQbGFjZTE8L2E+CiAg\r\nICAgICAgICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJtYXJnaW46\r\nIDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9r\r\nZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29y\r\nYXRpb246IG5vbmU7Ij50ZXN0IHBsYWNlIDE8L2E+CiAgICAgICAgICAgICAgICAgICAgPC9wPgogICAg\r\nICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgIDx0ZCB3\r\naWR0aD0iMTAiPjwvdGQ+CiAgICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSIyMTAi\r\nPgogICAgICAgICAgICAgICAgICAgIDwhLS1NYXAtLT4KICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAg\r\nICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQgY29sc3Bhbj0iNSIgaGVpZ2h0PSIxMCI+PC90\r\nZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQgY29sc3Bhbj0iNSI+CiAg\r\nICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFj\r\naW5nPSIwIiBzdHlsZT0iZm9udC1mYW1pbHk6IFZlcmRhbmE7IGZvbnQtc2l6ZTogMTNweDsgbGluZS1o\r\nZWlnaHQ6IDIwcHg7IGNvbG9yOiAjMTkxOTE5OyBmb250LXdlaWdodDogbm9ybWFsOyB3aWR0aDogMTAw\r\nJTsgYmFja2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsiPgogICAgICAgICAgICAgICAgICAgIAkKICAgICAg\r\nICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRo\r\nPSIxNSI+PGltZyBzcmM9Imh0dHA6Ly9zaXRlL2ltZy9lbWFpbC9yc3ZwX2FjY2VwdGVkXzEyX2JsdWUu\r\ncG5nIiAvPjwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50\r\nX2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjMTkxOTE5OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7\r\nIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9Im1hcmdpbjogMDsi\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHNwYW4gY2xhc3M9ImV4ZmVf\r\nbWFpbF9pZGVudGl0eV9uYW1lIj41PC9zcGFuPiBBY2NlcHRlZDogPHNwYW4gY2xhc3M9J2V4ZmVfbWFp\r\nbF9pZGVudGl0eV9uYW1lJz5lbWFpbDIgbmFtZTwvc3Bhbj4sIDxzcGFuIGNsYXNzPSdleGZlX21haWxf\r\naWRlbnRpdHlfbmFtZSc+ZmFjZWJvb2s1IG5hbWU8L3NwYW4+LCBhbmQgMSBvdGhlcnMuCiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvcD4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICA8L2E+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAg\r\nICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAg\r\nICAgCiAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nIDx0ZCB3aWR0aD0iMTUiPjxpbWcgc3JjPSJodHRwOi8vc2l0ZS9pbWcvZW1haWwvcnN2cF9kZWNsaW5l\r\nZF8xMi5wbmciIC8+PC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNp\r\ncGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMxOTE5MTk7IHRleHQtZGVjb3JhdGlvbjog\r\nbm9uZTsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0ibWFyZ2lu\r\nOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBVbmF2YWlsYWJsZTog\r\ndHdpdHRlcjMgbmFtZQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3A+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAg\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICA8dGQgd2lkdGg9IjE1Ij48aW1nIHNyYz0iaHR0cDovL3NpdGUv\r\naW1nL2VtYWlsL3BsdXNfMTJfYmx1ZS5wbmciIC8+PC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMxOTE5MTk7\r\nIHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICA8cCBzdHlsZT0ibWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICBOZXdseSBpbnZpdGVkOiA8c3BhbiBjbGFzcz0iZXhmZV9tYWlsX2lkZW50aXR5X25hbWUiPmZh\r\nY2Vib29rNSBuYW1lPC9zcGFuPi4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9w\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvYT4KICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAg\r\nICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgPHRy\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIxNSI+PGltZyBzcmM9Imh0dHA6\r\nLy9zaXRlL2ltZy9lbWFpbC9taW51c18xMi5wbmciIC8+PC90ZD4KICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8v\r\nc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMxOTE5\r\nMTk7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICA8cCBzdHlsZT0ibWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICBSZW1vdmVkOiBmYWNlYm9vazYgbmFtZS4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvYT4KICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgIDwvdGFibGU+CiAgICAgICAgICAg\r\nICAgICA8L3RkPgogICAgICAgICAgICA8L3RyPgogICAgICAgICAgICA8dHI+PHRkIGNvbHNwYW49IjUi\r\nIGhlaWdodD0iMTAiPjwvdGQ+PC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRk\r\nIGNvbHNwYW49IjUiPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgIAogICAg\r\nICAgICAgICAgICAgICAgIDxpbWcgc3R5bGU9InBhZGRpbmctcmlnaHQ6IDVweDsiIHdpZHRoPSI0MCIg\r\naGVpZ2h0PSI0MCIgYWx0PSJlbWFpbDIgbmFtZSIgdGl0bGU9ImVtYWlsMiBuYW1lIiBzcmM9Imh0dHA6\r\nLy9zaXRlL2FwaS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249MngmdXJsPWFIUjBjRG92TDNCaGRH\r\nZ3ZkRzh2WlcxaGFXd3lMbUYyWVhSaGNnJTNEJTNEJndpZHRoPTQwJmhlaWdodD00MCI+CiAgICAgICAg\r\nICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgPGltZyBzdHlsZT0icGFkZGluZy1yaWdodDog\r\nNXB4OyIgd2lkdGg9IjQwIiBoZWlnaHQ9IjQwIiBhbHQ9InR3aXR0ZXIzIG5hbWUiIHRpdGxlPSJ0d2l0\r\ndGVyMyBuYW1lIiBzcmM9Imh0dHA6Ly9zaXRlL2FwaS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249\r\nMngmdXJsPWFIUjBjRG92TDNCaGRHZ3ZkRzh2ZEhkcGRIUmxjak11WVhaaGRHRnkmd2lkdGg9NDAmaGVp\r\nZ2h0PTQwJmFscGhhPTAuMzMiPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAg\r\nIDxpbWcgc3R5bGU9InBhZGRpbmctcmlnaHQ6IDVweDsiIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgYWx0\r\nPSJmYWNlYm9vazQgbmFtZSIgdGl0bGU9ImZhY2Vib29rNCBuYW1lIiBzcmM9Imh0dHA6Ly9zaXRlL2Fw\r\naS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249MngmdXJsPWFIUjBjRG92TDNCaGRHZ3ZkRzh2Wm1G\r\nalpXSnZiMnMwTG1GMllYUmhjZyUzRCUzRCZ3aWR0aD00MCZoZWlnaHQ9NDAiPgogICAgICAgICAgICAg\r\nICAgICAgIAogICAgICAgICAgICAgICAgICAgIDxpbWcgc3R5bGU9InBhZGRpbmctcmlnaHQ6IDVweDsi\r\nIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgYWx0PSJ0d2l0dGVyMSBuYW1lIiB0aXRsZT0idHdpdHRlcjEg\r\nbmFtZSIgc3JjPSJodHRwOi8vc2l0ZS9hcGkvdjIvYXZhdGFyL3JlbmRlcj9yZXNvbHV0aW9uPTJ4JnVy\r\nbD1hSFIwY0RvdkwzQmhkR2d2ZEc4dmRIZHBkSFJsY2pFdVlYWmhkR0Z5JndpZHRoPTQwJmhlaWdodD00\r\nMCZhbHBoYT0wLjMzJmlzaG9zdD10cnVlJm1hdGVzPTIiPgogICAgICAgICAgICAgICAgICAgIAogICAg\r\nICAgICAgICAgICAgICAgIDxpbWcgc3R5bGU9InBhZGRpbmctcmlnaHQ6IDVweDsiIHdpZHRoPSI0MCIg\r\naGVpZ2h0PSI0MCIgYWx0PSJmYWNlYm9vazUgbmFtZSIgdGl0bGU9ImZhY2Vib29rNSBuYW1lIiBzcmM9\r\nImh0dHA6Ly9zaXRlL2FwaS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249MngmdXJsPWFIUjBjRG92\r\nTDNCaGRHZ3ZkRzh2Wm1GalpXSnZiMnMxTG1GMllYUmhjZyUzRCUzRCZ3aWR0aD00MCZoZWlnaHQ9NDAm\r\nbWF0ZXM9MiI+CiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICA8L3RkPgogICAgICAg\r\nICAgICA8L3RyPgogICAgICAgICAgICA8dHI+PHRkIGNvbHNwYW49IjUiIGhlaWdodD0iMTAiPjwvdGQ+\r\nPC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjUiPgogICAg\r\nICAgICAgICAgICAgICAgIDxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9l\r\nbWFpbDFfdG9rZW4iIHN0eWxlPSJjb2xvcjogIzMzMzMzMzsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+\r\ndGVzdCBjcm9zcyBkZXNjcmlwdGlvbjwvYT4KICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAg\r\nIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQgY29sc3Bhbj0iNSIgaGVpZ2h0PSIyMCI+PC90ZD48L3Ry\r\nPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQgY29sc3Bhbj0iNSIgc3R5bGU9ImZv\r\nbnQtc2l6ZTogMTFweDsgbGluZS1oZWlnaHQ6IDE1cHg7IGNvbG9yOiAjN0Y3RjdGOyI+CiAgICAgICAg\r\nICAgICAgICAgICAgUmVwbHkgdGhpcyBlbWFpbCBkaXJlY3RseSBhcyBjb252ZXJzYXRpb24uPCEtLSwg\r\nb3IgdHJ5IDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJl\r\nZj0iaHR0cDovL2FwcC91cmwiPkVYRkU8L2E+IGFwcC4tLT4KICAgICAgICAgICAgICAgICAgICA8YnIg\r\nLz4KICAgICAgICAgICAgICAgICAgICA8c3BhbiBzdHlsZT0iY29sb3I6ICNCMkIyQjIiPlRoaXMgdXBk\r\nYXRlIGlzIHNlbnQgZnJvbSA8YSBzdHlsZT0iY29sb3I6ICMzYTZlYTU7IHRleHQtZGVjb3JhdGlvbjog\r\nbm9uZTsiIGhyZWY9Imh0dHA6Ly9zaXRlL3VybCI+RVhGRTwvYT4gYXV0b21hdGljYWxseS4gPGEgc3R5\r\nbGU9ImNvbG9yOiAjRTZFNkU2OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8vc2l0\r\nZS91cmwvcy9yZXBvcnRTcGFtP3Rva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iPlVuc3Vic2NyaWJl\r\nPzwvYT4KICAgICAgICAgICAgICAgICAgICA8IS0tCiAgICAgICAgICAgICAgICAgICAgWW91IGNhbiBj\r\naGFuZ2UKICAgICAgICAgICAgICAgICAgICA8YSBzdHlsZT0iY29sb3I6ICNCMkIyQjI7IHRleHQtZGVj\r\nb3JhdGlvbjogdW5kZWxpbmU7IiBocmVmPSIiPm5vdGlmaWNhdGlvbiBwcmVmZXJlbmNlPC9hPi4KICAg\r\nICAgICAgICAgICAgICAgICAtLT4KICAgICAgICAgICAgICAgICAgICA8L3NwYW4+CiAgICAgICAgICAg\r\nICAgICA8L3RkPgogICAgICAgICAgICA8L3RyPgogICAgICAgIDwvdGFibGU+CiAgICA8L2JvZHk+Cjwv\r\naHRtbD4K\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6TmV3IFRpdGxlDQpYLVdSLUNBTERFU0M6ZXhmZSBjYWwNClgtV1ItVElNRVpPTkU6\r\nKzA4MDANCkJFR0lOOlZFVkVOVA0KVUlEOiExMjNAZXhmZQ0KRFRTVEFNUDoyMDEyMTAyM1QwODQ1MDBa\r\nDQpERVNDUklQVElPTjp0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uDQpEVFNUQVJUOjIwMTIxMDIzVDA4NDUw\r\nMFoNCkxPQ0FUSU9OOlRlc3QgUGxhY2UxXG50ZXN0IHBsYWNlIDENClNVTU1BUlk6TmV3IFRpdGxlDQpV\r\nUkw6aHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbg0KRU5EOlZFVkVO\r\nVA0KRU5EOlZDQUxFTkRBUg0K\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
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
	expectPrivate := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <x+123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <to_email_address>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x+123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nWW91J3JlIGdhdGhlcmluZyB0aGlzIMK3WMK3LgoKClRlc3QgQ3Jvc3MKPT09PT09PQpodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuCgo0OjQ1UE0gb24gVHVlLCBPY3QgMjMK\r\nPT09PT09PQoKUGxhY2UKPT09PT09PQogIFRvIGJlIGRlY2lkZWQuCgoKSSdtIGluLiBDaGVjayBpdCBv\r\ndXQ6IGh0dHA6Ly9zaXRlL3VybC8/dG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiZyc3ZwPWFjY2Vw\r\ndAoKCjYgSW52aXRlZDoKwrcgZW1haWwxIG5hbWUgKEhvc3QpIHdpdGggMiBwZW9wbGUKwrcgZW1haWwy\r\nIG5hbWUKwrcgdHdpdHRlcjMgbmFtZQrCtyBmYWNlYm9vazQgbmFtZQoKCkRlc2NyaXB0aW9uCi0tLS0t\r\nLS0KICB0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uCgoKIyBSZXBseSB0aGlzIGVtYWlsIGRpcmVjdGx5IGFz\r\nIGNvbnZlcnNhdGlvbi4gIw==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9sYWJlbCB7\r\nCiAgICAgICAgICAgICAgICBiYWNrZ3JvdW5kLWNvbG9yOiAjRDVFOEYyOwogICAgICAgICAgICAgICAg\r\nY29sb3I6ICMzYTZlYTU7CiAgICAgICAgICAgICAgICBmb250LXNpemU6IDExcHg7CiAgICAgICAgICAg\r\nICAgICBwYWRkaW5nOiAwIDJweCAwIDJweDsKICAgICAgICAgICAgfQogICAgICAgICAgICAuZXhmZV9t\r\nYWlsX21hdGVzIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICAgICAg\r\nZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICB9CiAgICAgICAgICAgIC5leGZlX21haWxfaWRlbnRp\r\ndHkgewogICAgICAgICAgICAgICAgZm9udC1zdHlsZTogaXRhbGljOwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfaWRlbnRpdHlfbmFtZSB7CiAgICAgICAgICAgICAgICBjb2xvcjogIzE5\r\nMTkxOTsKICAgICAgICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4K\r\nICAgICAgICA8dGFibGUgd2lkdGg9IjY0MCIgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNw\r\nYWNpbmc9IjAiIHN0eWxlPSJmb250LWZhbWlseTogSGVsdmV0aWNhOyBmb250LXNpemU6IDEzcHg7IGxp\r\nbmUtaGVpZ2h0OiAxOXB4OyBjb2xvcjogIzE5MTkxOTsgZm9udC13ZWlnaHQ6IG5vcm1hbDsgcGFkZGlu\r\nZzogMzBweCA0MHB4IDMwcHggNDBweDsgYmFja2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsgbWluLWhlaWdo\r\ndDogNTYycHg7Ij4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjMi\r\nIHZhbGlnbj0idG9wIiBzdHlsZT0iZm9udC1zaXplOiAzMnB4OyBsaW5lLWhlaWdodDogMzhweDsgcGFk\r\nZGluZy1ib3R0b206IDE4cHg7Ij4KICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMzYTZlYTU7\r\nIHRleHQtZGVjb3JhdGlvbjogbm9uZTsgZm9udC13ZWlnaHQ6IDMwMDsiPgogICAgICAgICAgICAgICAg\r\nICAgICAgICBUZXN0IENyb3NzCiAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAg\r\nPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRk\r\nIHdpZHRoPSIzNDAiIHN0eWxlPSJ2ZXJ0aWNhbC1hbGlnbjogYmFzZWxpbmU7IGZvbnQtd2VpZ2h0OiAz\r\nMDA7Ij4KICAgICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIg\r\nY2VsbHNwYWNpbmc9IjAiPgogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHN0eWxlPSJwYWRkaW5nLWJvdHRvbTogMjBweDsg\r\nZm9udC1zaXplOiAyMHB4OyB2ZXJ0aWNhbC1hbGlnbjogYmFzZWxpbmU7Ij4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBZb3UncmUgZ2F0\r\naGVyaW5nIHRoaXMgPHNwYW4gc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5v\r\nbmU7Ij7Ct1jCtzwvc3Bhbj4uCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAg\r\nICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2VsbHBhZGRpbmc9\r\nIjAiIGNlbGxzcGFjaW5nPSIwIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRy\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIiB3\r\naWR0aD0iMTYwIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBo\r\ncmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0i\r\ndGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgCQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9ImZv\r\nbnQtc2l6ZTogMjBweDsgbGluZS1oZWlnaHQ6IDI2cHg7IG1hcmdpbjogMDsgY29sb3I6ICMzMzMzMzM7\r\nIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDQ6NDVQ\r\nTSBvbiBUdWUsIE9jdCAyMwogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICA8L3A+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvYT4KICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHN0eWxlPSJwYWRkaW5nLWxlZnQ6IDEwcHg7Ij4K\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8v\r\nc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0idGV4dC1kZWNvcmF0\r\naW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0i\r\nZm9udC1zaXplOiAyMHB4OyBsaW5lLWhlaWdodDogMjZweDsgbWFyZ2luOiAwOyBjb2xvcjogIzMzMzMz\r\nMzsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgUGxh\r\nY2UKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9wPgogICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0iY29sb3I6\r\nICMxOTE5MTk7IG1hcmdpbjogMDsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgVG8gYmUgZGVjaWRlZAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8L3A+IAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L2E+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nIDwvdGFibGU+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAg\r\nICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAg\r\nICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIiBzdHlsZT0i\r\ncGFkZGluZy10b3A6IDMwcHg7IHBhZGRpbmctYm90dG9tOiAzMHB4OyI+CiAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgPGEgc3R5bGU9ImZsb2F0OiBsZWZ0OyBkaXNwbGF5OiBibG9jazsgdGV4dC1k\r\nZWNvcmF0aW9uOiBub25lOyBib3JkZXI6IDFweCBzb2xpZCAjYmViZWJlOyBiYWNrZ3JvdW5kLWNvbG9y\r\nOiAjM0E2RUE1OyBjb2xvcjogI0ZGRkZGRjsgcGFkZGluZzogNXB4IDMwcHggNXB4IDMwcHg7IG1hcmdp\r\nbi1sZWZ0OiAyNXB4OyIgYWx0PSJBY2NlcHQiIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8/dG9rZW49cmVj\r\naXBpZW50X2VtYWlsMV90b2tlbiZyc3ZwPWFjY2VwdCI+SSdtIGluPC9hPgogICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgIDxhIHN0eWxlPSJmbG9hdDogbGVmdDsgZGlzcGxheTogYmxvY2s7IHRleHQt\r\nZGVjb3JhdGlvbjogbm9uZTsgYm9yZGVyOiAxcHggc29saWQgI2JlYmViZTsgYmFja2dyb3VuZC1jb2xv\r\ncjogI0U2RTZFNjsgY29sb3I6ICMxOTE5MTk7IHBhZGRpbmc6IDVweCAyNXB4IDVweCAyNXB4OyBtYXJn\r\naW4tbGVmdDogMTVweDsiIGFsdD0iQ2hlY2sgaXQgb3V0IiBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0\r\nb2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIj5DaGVjayBpdCBvdXQuLi48L2E+CiAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAg\r\nICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICB0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nCiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8\r\nL3RyPgogICAgICAgICAgICAgICAgICAgIDwvdGFibGU+CiAgICAgICAgICAgICAgICA8L3RkPgogICAg\r\nICAgICAgICAgICAgPHRkIHdpZHRoPSIzMCI+PC90ZD4KICAgICAgICAgICAgICAgIDx0ZCB2YWxpZ249\r\nInRvcCI+CiAgICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2VsbHBhZGRpbmc9IjAi\r\nIGNlbGxzcGFjaW5nPSIwIj4KICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgPHRkIGhlaWdodD0iNjgiIHZhbGlnbj0iYm90dG9tIiBhbGlnbj0icmlnaHQi\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwhLS08c3BhbiBzdHlsZT0iY29sb3I6ICMz\r\nYTZlYTU7Ij4zPC9zcGFuPiBjb25maXJtZWQtLT4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwv\r\ndGQ+CiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgICAgIDx0\r\ncj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBzdHlsZT0iY29sb3I6ICMzMzMzMzM7IiBjZWxscGFkZGlu\r\nZz0iMCIgY2VsbHNwYWNpbmc9IjAiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIyNSIgaGVpZ2h0PSIyNSIgYWxpZ249ImxlZnQi\r\nIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8\r\naW1nIHdpZHRoPSIyMCIgaGVpZ2h0PSIyMCIgdGl0bGU9ImVtYWlsMSBuYW1lIiBhbHQ9ImVtYWlsMSBu\r\nYW1lIiBzcmM9Imh0dHA6Ly9wYXRoL3RvL2VtYWlsMS5hdmF0YXIiPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHNwYW4+\r\nZW1haWwxIG5hbWU8L3NwYW4+IDxzcGFuIGNsYXNzPSJleGZlX21haWxfbWF0ZXMiPisyPC9zcGFuPiA8\r\nc3BhbiBjbGFzcz0iZXhmZV9tYWlsX2xhYmVsIj5ob3N0PC9zcGFuPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDx0ZCB3aWR0aD0iMjUiIGhlaWdodD0iMjUiIGFsaWduPSJsZWZ0IiB2YWxpZ249InRvcCI+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPGltZyB3aWR0aD0iMjAi\r\nIGhlaWdodD0iMjAiIHRpdGxlPSJlbWFpbDIgbmFtZSIgYWx0PSJlbWFpbDIgbmFtZSIgc3JjPSJodHRw\r\nOi8vcGF0aC90by9lbWFpbDIuYXZhdGFyIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxzcGFuPmVtYWlsMiBuYW1lPC9z\r\ncGFuPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMjUiIGhlaWdodD0iMjUiIGFs\r\naWduPSJsZWZ0IiB2YWxpZ249InRvcCI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgPGltZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIHRpdGxlPSJ0d2l0dGVyMyBuYW1lIiBh\r\nbHQ9InR3aXR0ZXIzIG5hbWUiIHNyYz0iaHR0cDovL3BhdGgvdG8vdHdpdHRlcjMuYXZhdGFyIj4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgIDxzcGFuPnR3aXR0ZXIzIG5hbWU8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8\r\nL3RyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgPHRkIHdpZHRoPSIyNSIgaGVpZ2h0PSIyNSIgYWxpZ249ImxlZnQiIHZhbGlnbj0idG9wIj4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8aW1nIHdpZHRoPSIyMCIgaGVp\r\nZ2h0PSIyMCIgdGl0bGU9ImZhY2Vib29rNCBuYW1lIiBhbHQ9ImZhY2Vib29rNCBuYW1lIiBzcmM9Imh0\r\ndHA6Ly9wYXRoL3RvL2ZhY2Vib29rNC5hdmF0YXIiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0\r\nZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHNwYW4+ZmFjZWJvb2s0\r\nIG5hbWU8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RhYmxlPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90\r\ncj4KICAgICAgICAgICAgICAgICAgICA8L3RhYmxlPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAg\r\nICAgICAgPC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjMi\r\nIHN0eWxlPSJmb250LXNpemU6IDExcHg7IGxpbmUtaGVpZ2h0OiAxNXB4OyBjb2xvcjogIzdGN0Y3Rjsg\r\ncGFkZGluZy10b3A6IDQwcHg7Ij4KICAgICAgICAgICAgICAgICAgICBSZXBseSB0aGlzIGVtYWlsIGRp\r\ncmVjdGx5IGFzIGNvbnZlcnNhdGlvbi48IS0tICwgb3IgVHJ5IDxhIHN0eWxlPSJjb2xvcjogIzNhNmVh\r\nNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJlZj0iaHR0cDovL2FwcC91cmwiPkVYRkU8L2E+IGFw\r\ncC4tLT4KICAgICAgICAgICAgICAgICAgICA8YnIgLz4KICAgICAgICAgICAgICAgICAgICBUaGlzIDxh\r\nIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJlZj0iaHR0cDov\r\nL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiI+wrdYwrc8L2E+IGludml0YXRp\r\nb24gaXMgc2VudCBieSA8c3BhbiBjbGFzcz0iZXhmZV9tYWlsX2lkZW50aXR5X25hbWUiPmVtYWlsMSBu\r\nYW1lPC9zcGFuPiBmcm9tIDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBu\r\nb25lOyIgaHJlZj0iaHR0cDovL3NpdGUvdXJsIj5FWEZFPC9hPi4KICAgICAgICAgICAgICAgIDwvdGQ+\r\nCiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgPC90YWJsZT4KICAgIDwvYm9keT4KPC9odG1sPgo=\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6VGVzdCBDcm9zcw0KWC1XUi1DQUxERVNDOmV4ZmUgY2FsDQpYLVdSLVRJTUVaT05F\r\nOiswODAwDQpCRUdJTjpWRVZFTlQNClVJRDohMTIzQGV4ZmUNCkRUU1RBTVA6MjAxMjEwMjNUMDg0NTAw\r\nWg0KREVTQ1JJUFRJT046dGVzdCBjcm9zcyBkZXNjcmlwdGlvbg0KRFRTVEFSVDoyMDEyMTAyM1QwODQ1\r\nMDBaDQpMT0NBVElPTjoNClNVTU1BUlk6VGVzdCBDcm9zcw0KVVJMOmh0dHA6Ly9zaXRlL3VybC8jIXRv\r\na2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4NCkVORDpWRVZFTlQNCkVORDpWQ0FMRU5EQVINCg==\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
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
	expectPrivate := "\n\n\n\nSuccessfully gathering \\(“Test Cross”\\), \\(4:45PM on Tue, Oct 23\\). 6 invited: email1 name, email2 name, twitter3 name… http://site/url/#!token=recipient_twitter1_token"
	assert.Equal(t, private, expectPrivate)
	assert.Equal(t, public, "Invitation: http://site/url/#!123/eci (Please follow @EXFE to receive details PRIVATELY through Direct Message.)")
}
