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

	c := NewCross(localTemplate, &config)
	private, public, err := c.getContent(updates)
	t.Logf("err: %s", err)
	t.Errorf("private:-----start------\n%s\n-------end-------", private)
	t.Errorf("public:-----start------\n%s\n-------end-------", public)
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

	c := NewCross(localTemplate, &config)
	private, public, err := c.getContent(updates)
	assert.Equal(t, err, nil)
	t.Logf("private:-----start------\n%s\n-------end-------", private)
	expectPrivate := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <x+123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <email1@domain.com>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x+123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nVXBkYXRlcyBvZiDCt1jCtyDigJxUZXN0IENyb3Nz4oCdIGJ5IGZhY2Vib29rNCBuYW1lLCBlbWFpbDEg\r\nbmFtZSwgZW1haWwyIG5hbWUsIGV0Yy4KCipOZXcgVGl0bGUqCj09PT09PT0KaHR0cDovL3NpdGUvdXJs\r\nLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbgoKKjQ6NDVQTSBvbiBUdWUsIE9jdCAyMyoKPT09\r\nPT09PQoKKlRlc3QgUGxhY2UxKgo9PT09PT09CiAgKnRlc3QgcGxhY2UgMSoKCgrCtyA1IEFjY2VwdGVk\r\nOiBlbWFpbDIgbmFtZSwgZmFjZWJvb2s1IG5hbWUgYW5kIDEgb3RoZXJzLgrCtyBVbmF2YWlsYWJsZTog\r\ndHdpdHRlcjMgbmFtZS4KwrcgTmV3bHkgaW52aXRlZDogZmFjZWJvb2s1IG5hbWUuCsK3IFJlbW92ZWQ6\r\nIGZhY2Vib29rNiBuYW1lLgoKIyBSZXBseSB0aGlzIGVtYWlsIGRpcmVjdGx5IGFzIGNvbnZlcnNhdGlv\r\nbi4gIw==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9pZGVudGl0\r\neV9uYW1lIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfbXNnX2lkZW50aXR5X25hbWUgewogICAgICAgICAgICAgICAgY29sb3I6\r\nICM2NjY2NjY7CiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfYXQgewogICAg\r\nICAgICAgICAgICAgZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICAgICAgY29sb3I6ICM5OTk5OTk7\r\nCiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfdGltZSB7CiAgICAgICAgICAg\r\nICAgICBmb250LXNpemU6IDEycHg7CiAgICAgICAgICAgICAgICBjb2xvcjogIzY2NjY2NjsKICAgICAg\r\nICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4KICAgICAgICA8dGFi\r\nbGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAiIHN0eWxlPSJmb250LWZh\r\nbWlseTogVmVyZGFuYTsgZm9udC1zaXplOiAxM3B4OyBsaW5lLWhlaWdodDogMjBweDsgY29sb3I6ICMx\r\nOTE5MTk7IGZvbnQtd2VpZ2h0OiBub3JtYWw7IHdpZHRoOiA2NDBweDsgcGFkZGluZzogMjBweDsgYmFj\r\na2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsiPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8\r\ndGQgY29sc3Bhbj0iNSIgc3R5bGU9ImNvbG9yOiAjMzMzMzMzOyI+CiAgICAgICAgICAgICAgICAgICAg\r\nPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5\r\nbGU9ImNvbG9yOiAjMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5VcGRhdGVzIG9mIDxzcGFu\r\nIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsiPsK3WMK3PC9zcGFuPiDigJw8c3BhbiBzdHlsZT0iY29sb3I6\r\nICMxOTE5MTk7Ij5UZXN0IENyb3NzPC9zcGFuPuKAnSBieSBmYWNlYm9vazQgbmFtZSwgZW1haWwxIG5h\r\nbWUsIGVtYWlsMiBuYW1lLCBldGMuPC9hPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAg\r\nPC90cj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjEwIj48L3RkPjwvdHI+\r\nCiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1IiBzdHlsZT0iZm9u\r\ndC1zaXplOiAyMHB4OyBsaW5lLWhlaWdodDogMjZweDsiPgogICAgICAgICAgICAgICAgICAgIDxhIGhy\r\nZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJj\r\nb2xvcjojM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IGZvbnQtd2VpZ2h0OiBsaWdodGVyOyI+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIE5ldyBUaXRsZQogICAgICAgICAgICAgICAgICAgIDwvYT4K\r\nICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQg\r\nY29sc3Bhbj0iNSIgaGVpZ2h0PSIxMCI+PC90ZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAg\r\nICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSIxODAiPgogICAgICAgICAgICAgICAgICAgIAog\r\nICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJmb250LXNpemU6IDIwcHg7IGxpbmUtaGVpZ2h0OiAy\r\nNnB4OyBtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3Np\r\ndGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1\r\nOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij40OjQ1UE0gb24gVHVlLCBPY3QgMjM8L2E+CiAgICAgICAg\r\nICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4K\r\nICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMTAiPjwvdGQ+CiAgICAgICAgICAgICAgICA8dGQgdmFs\r\naWduPSJ0b3AiIHdpZHRoPSIxOTAiIHN0eWxlPSJ3b3JkLWJyZWFrOiBicmVhay1hbGw7Ij4KICAgICAg\r\nICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0iZm9udC1zaXplOiAyMHB4\r\nOyBsaW5lLWhlaWdodDogMjZweDsgbWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgIDxh\r\nIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxl\r\nPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+VGVzdCBQbGFjZTE8L2E+CiAg\r\nICAgICAgICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJtYXJnaW46\r\nIDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9r\r\nZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29y\r\nYXRpb246IG5vbmU7Ij50ZXN0IHBsYWNlIDE8L2E+CiAgICAgICAgICAgICAgICAgICAgPC9wPgogICAg\r\nICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgIDx0ZCB3\r\naWR0aD0iMTAiPjwvdGQ+CiAgICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSIyMTAi\r\nPgogICAgICAgICAgICAgICAgICAgIDwhLS1NYXAtLT4KICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAg\r\nICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQgY29sc3Bhbj0iNSIgaGVpZ2h0PSIxMCI+PC90\r\nZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQgY29sc3Bhbj0iNSI+CiAg\r\nICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFj\r\naW5nPSIwIiBzdHlsZT0iZm9udC1mYW1pbHk6IFZlcmRhbmE7IGZvbnQtc2l6ZTogMTNweDsgbGluZS1o\r\nZWlnaHQ6IDIwcHg7IGNvbG9yOiAjMTkxOTE5OyBmb250LXdlaWdodDogbm9ybWFsOyB3aWR0aDogMTAw\r\nJTsgYmFja2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsiPgogICAgICAgICAgICAgICAgICAgIAkKICAgICAg\r\nICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRo\r\nPSIxNSI+PGltZyBzcmM9Imh0dHA6Ly9zaXRlL2ltZy9lbWFpbC9yc3ZwX2FjY2VwdGVkXzEyX2JsdWUu\r\ncG5nIiAvPjwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50\r\nX2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjMTkxOTE5OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7\r\nIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9Im1hcmdpbjogMDsi\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHNwYW4gY2xhc3M9ImV4ZmVf\r\nbWFpbF9pZGVudGl0eV9uYW1lIj41PC9zcGFuPiBBY2NlcHRlZDogPHNwYW4gY2xhc3M9J2V4ZmVfbWFp\r\nbF9pZGVudGl0eV9uYW1lJz5lbWFpbDIgbmFtZTwvc3Bhbj4sIDxzcGFuIGNsYXNzPSdleGZlX21haWxf\r\naWRlbnRpdHlfbmFtZSc+ZmFjZWJvb2s1IG5hbWU8L3NwYW4+LCBhbmQgMSBvdGhlcnMuCiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvcD4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICA8L2E+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAg\r\nICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAg\r\nICAgCiAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nIDx0ZCB3aWR0aD0iMTUiPjxpbWcgc3JjPSJodHRwOi8vc2l0ZS9pbWcvZW1haWwvcnN2cF9kZWNsaW5l\r\nZF8xMi5wbmciIC8+PC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNp\r\ncGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMxOTE5MTk7IHRleHQtZGVjb3JhdGlvbjog\r\nbm9uZTsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0ibWFyZ2lu\r\nOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBVbmF2YWlsYWJsZTog\r\ndHdpdHRlcjMgbmFtZQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3A+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAg\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICA8dGQgd2lkdGg9IjE1Ij48aW1nIHNyYz0iaHR0cDovL3NpdGUv\r\naW1nL2VtYWlsL3BsdXNfMTJfYmx1ZS5wbmciIC8+PC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMxOTE5MTk7\r\nIHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICA8cCBzdHlsZT0ibWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICBOZXdseSBpbnZpdGVkOiA8c3BhbiBjbGFzcz0iZXhmZV9tYWlsX2lkZW50aXR5X25hbWUiPmZh\r\nY2Vib29rNSBuYW1lPC9zcGFuPi4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9w\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvYT4KICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAg\r\nICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgPHRy\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIxNSI+PGltZyBzcmM9Imh0dHA6\r\nLy9zaXRlL2ltZy9lbWFpbC9taW51c18xMi5wbmciIC8+PC90ZD4KICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8v\r\nc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMxOTE5\r\nMTk7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICA8cCBzdHlsZT0ibWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICBSZW1vdmVkOiBmYWNlYm9vazYgbmFtZS4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvYT4KICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgIDwvdGFibGU+CiAgICAgICAgICAg\r\nICAgICA8L3RkPgogICAgICAgICAgICA8L3RyPgogICAgICAgICAgICA8dHI+PHRkIGNvbHNwYW49IjUi\r\nIGhlaWdodD0iMTAiPjwvdGQ+PC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRk\r\nIGNvbHNwYW49IjUiPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgIAogICAg\r\nICAgICAgICAgICAgICAgIDxpbWcgc3R5bGU9InBhZGRpbmctcmlnaHQ6IDVweDsiIHdpZHRoPSI0MCIg\r\naGVpZ2h0PSI0MCIgYWx0PSJlbWFpbDIgbmFtZSIgdGl0bGU9ImVtYWlsMiBuYW1lIiBzcmM9Imh0dHA6\r\nLy9zaXRlL2FwaS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249MngmdXJsPWFIUjBjRG92TDNCaGRH\r\nZ3ZkRzh2WlcxaGFXd3lMbUYyWVhSaGNnJTNEJTNEJndpZHRoPTQwJmhlaWdodD00MCI+CiAgICAgICAg\r\nICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgPGltZyBzdHlsZT0icGFkZGluZy1yaWdodDog\r\nNXB4OyIgd2lkdGg9IjQwIiBoZWlnaHQ9IjQwIiBhbHQ9InR3aXR0ZXIzIG5hbWUiIHRpdGxlPSJ0d2l0\r\ndGVyMyBuYW1lIiBzcmM9Imh0dHA6Ly9zaXRlL2FwaS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249\r\nMngmdXJsPWFIUjBjRG92TDNCaGRHZ3ZkRzh2ZEhkcGRIUmxjak11WVhaaGRHRnkmd2lkdGg9NDAmaGVp\r\nZ2h0PTQwJmFscGhhPTAuMzMiPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAg\r\nIDxpbWcgc3R5bGU9InBhZGRpbmctcmlnaHQ6IDVweDsiIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgYWx0\r\nPSJmYWNlYm9vazQgbmFtZSIgdGl0bGU9ImZhY2Vib29rNCBuYW1lIiBzcmM9Imh0dHA6Ly9zaXRlL2Fw\r\naS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249MngmdXJsPWFIUjBjRG92TDNCaGRHZ3ZkRzh2Wm1G\r\nalpXSnZiMnMwTG1GMllYUmhjZyUzRCUzRCZ3aWR0aD00MCZoZWlnaHQ9NDAiPgogICAgICAgICAgICAg\r\nICAgICAgIAogICAgICAgICAgICAgICAgICAgIDxpbWcgc3R5bGU9InBhZGRpbmctcmlnaHQ6IDVweDsi\r\nIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgYWx0PSJ0d2l0dGVyMSBuYW1lIiB0aXRsZT0idHdpdHRlcjEg\r\nbmFtZSIgc3JjPSJodHRwOi8vc2l0ZS9hcGkvdjIvYXZhdGFyL3JlbmRlcj9yZXNvbHV0aW9uPTJ4JnVy\r\nbD1hSFIwY0RvdkwzQmhkR2d2ZEc4dmRIZHBkSFJsY2pFdVlYWmhkR0Z5JndpZHRoPTQwJmhlaWdodD00\r\nMCZhbHBoYT0wLjMzJmlzaG9zdD10cnVlJm1hdGVzPTIiPgogICAgICAgICAgICAgICAgICAgIAogICAg\r\nICAgICAgICAgICAgICAgIDxpbWcgc3R5bGU9InBhZGRpbmctcmlnaHQ6IDVweDsiIHdpZHRoPSI0MCIg\r\naGVpZ2h0PSI0MCIgYWx0PSJmYWNlYm9vazUgbmFtZSIgdGl0bGU9ImZhY2Vib29rNSBuYW1lIiBzcmM9\r\nImh0dHA6Ly9zaXRlL2FwaS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249MngmdXJsPWFIUjBjRG92\r\nTDNCaGRHZ3ZkRzh2Wm1GalpXSnZiMnMxTG1GMllYUmhjZyUzRCUzRCZ3aWR0aD00MCZoZWlnaHQ9NDAm\r\nbWF0ZXM9MiI+CiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICA8L3RkPgogICAgICAg\r\nICAgICA8L3RyPgogICAgICAgICAgICA8dHI+PHRkIGNvbHNwYW49IjUiIGhlaWdodD0iMjAiPjwvdGQ+\r\nPC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjUiIHN0eWxl\r\nPSJmb250LXNpemU6IDExcHg7IGxpbmUtaGVpZ2h0OiAxNXB4OyBjb2xvcjogIzdGN0Y3RjsiPgogICAg\r\nICAgICAgICAgICAgICAgIFJlcGx5IHRoaXMgZW1haWwgZGlyZWN0bHkgYXMgY29udmVyc2F0aW9uLjwh\r\nLS0sIG9yIHRyeSA8YSBzdHlsZT0iY29sb3I6ICMzYTZlYTU7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsi\r\nIGhyZWY9Imh0dHA6Ly9hcHAvdXJsIj5FWEZFPC9hPiBhcHAuLS0+CiAgICAgICAgICAgICAgICAgICAg\r\nPGJyIC8+CiAgICAgICAgICAgICAgICAgICAgPHNwYW4gc3R5bGU9ImNvbG9yOiAjQjJCMkIyIj5UaGlz\r\nIHVwZGF0ZSBpcyBzZW50IGZyb20gPGEgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRp\r\nb246IG5vbmU7IiBocmVmPSJodHRwOi8vc2l0ZS91cmwiPkVYRkU8L2E+IGF1dG9tYXRpY2FsbHkuIDxh\r\nIHN0eWxlPSJjb2xvcjogI0U2RTZFNjsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJlZj0iaHR0cDov\r\nL3NpdGUvdXJsL3MvcmVwb3J0U3BhbT90b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIj5VbnN1YnNj\r\ncmliZT88L2E+CiAgICAgICAgICAgICAgICAgICAgPCEtLQogICAgICAgICAgICAgICAgICAgIFlvdSBj\r\nYW4gY2hhbmdlCiAgICAgICAgICAgICAgICAgICAgPGEgc3R5bGU9ImNvbG9yOiAjQjJCMkIyOyB0ZXh0\r\nLWRlY29yYXRpb246IHVuZGVsaW5lOyIgaHJlZj0iIj5ub3RpZmljYXRpb24gcHJlZmVyZW5jZTwvYT4u\r\nCiAgICAgICAgICAgICAgICAgICAgLS0+CiAgICAgICAgICAgICAgICAgICAgPC9zcGFuPgogICAgICAg\r\nICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICA8L3RhYmxlPgogICAgPC9ib2R5\r\nPgo8L2h0bWw+Cg==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6TmV3IFRpdGxlDQpYLVdSLUNBTERFU0M6ZXhmZSBjYWwNClgtV1ItVElNRVpPTkU6\r\nKzA4MDANCkJFR0lOOlZFVkVOVA0KVUlEOiExMjNAZXhmZQ0KRFRTVEFNUDoyMDEyMTAyM1QwODQ1MDBa\r\nDQpERVNDUklQVElPTjp0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uDQpEVFNUQVJUOjIwMTIxMDIzVDA4NDUw\r\nMFoNCkxPQ0FUSU9OOlRlc3QgUGxhY2UxXG50ZXN0IHBsYWNlIDENClNVTU1BUlk6TmV3IFRpdGxlDQpV\r\nUkw6aHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbg0KRU5EOlZFVkVO\r\nVA0KRU5EOlZDQUxFTkRBUg0K\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
	assert.Equal(t, private, expectPrivate)
	assert.Equal(t, public, "")
}
