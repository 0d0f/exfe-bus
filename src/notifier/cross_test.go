package notifier

import (
	"formatter"
	"github.com/stretchrcom/testify/assert"
	"model"
	"testing"
)

func init() {
	config.Email.Name = "business tester"
	config.Email.Domain = "test.com"
	config.TemplatePath = "../../templates"
	var err error
	localTemplate, err = formatter.NewLocalTemplate(config.TemplatePath, "en_US")
	if err != nil {
		panic(err)
	}
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
	cross1.Time = &time1
	updates = append(updates, model.CrossUpdate{
		To:       rtwitter1,
		OldCross: cross2,
		Cross:    cross1,
		By:       email1,
	})

	cross2 = cross1
	cross1.Place = &place1
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
	text, err := c.getSummaryContent(updates)
	assert.Equal(t, err, nil)
	t.Logf("text:-----start------\n%s\n-------end-------", text)
	expect := ""
	assert.Equal(t, text, expect)
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
	cross1.Time = &time1
	updates = append(updates, model.CrossUpdate{
		To:       remail1,
		OldCross: cross2,
		Cross:    cross1,
		By:       email1,
	})

	cross2 = cross1
	cross1.Place = &place1
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
	text, err := c.getSummaryContent(updates)
	assert.Equal(t, err, nil)
	t.Logf("text:-----start------\n%s\n-------end-------", text)
	expect := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <+123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <to_email_address>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <+123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nVXBkYXRlcyBvZiDCt1jCtyDigJxUZXN0IENyb3Nz4oCdIGJ5IGZhY2Vib29rNCBuYW1lLCBlbWFpbDEg\r\nbmFtZSwgZW1haWwyIG5hbWUsIGV0Yy4KCipOZXcgVGl0bGUqCj09PT09PT0KaHR0cDovL3NpdGUvdXJs\r\nLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbgoKKjQ6NDVQTSBvbiBUdWUsIE9jdCAyMyAyMDEy\r\nKgo9PT09PT09CgoqVGVzdCBQbGFjZTEqCj09PT09PT0KICAqdGVzdCBwbGFjZSAxKgoKCsK3IDUgQWNj\r\nZXB0ZWQ6IGVtYWlsMiBuYW1lLCBmYWNlYm9vazUgbmFtZSBhbmQgMSBvdGhlcnMuCsK3IFVuYXZhaWxh\r\nYmxlOiB0d2l0dGVyMyBuYW1lLgrCtyBOZXdseSBpbnZpdGVkOiBmYWNlYm9vazUgbmFtZS4KwrcgUmVt\r\nb3ZlZDogZmFjZWJvb2s2IG5hbWUuCgojIFJlcGx5IHRoaXMgZW1haWwgZGlyZWN0bHkgYXMgY29udmVy\r\nc2F0aW9uLiAj\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9pZGVudGl0\r\neV9uYW1lIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfbXNnX2lkZW50aXR5X25hbWUgewogICAgICAgICAgICAgICAgY29sb3I6\r\nICM2NjY2NjY7CiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfYXQgewogICAg\r\nICAgICAgICAgICAgZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICAgICAgY29sb3I6ICM5OTk5OTk7\r\nCiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfdGltZSB7CiAgICAgICAgICAg\r\nICAgICBmb250LXNpemU6IDEycHg7CiAgICAgICAgICAgICAgICBjb2xvcjogIzY2NjY2NjsKICAgICAg\r\nICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4KICAgICAgICA8dGFi\r\nbGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAiIHN0eWxlPSJmb250LWZh\r\nbWlseTogVmVyZGFuYTsgZm9udC1zaXplOiAxM3B4OyBsaW5lLWhlaWdodDogMjBweDsgY29sb3I6ICMx\r\nOTE5MTk7IGZvbnQtd2VpZ2h0OiBub3JtYWw7IHdpZHRoOiA2NDBweDsgcGFkZGluZzogMjBweDsgYmFj\r\na2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsiPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8\r\ndGQgY29sc3Bhbj0iNSIgc3R5bGU9ImNvbG9yOiAjMzMzMzMzOyI+CiAgICAgICAgICAgICAgICAgICAg\r\nPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5\r\nbGU9ImNvbG9yOiAjMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5VcGRhdGVzIG9mIDxzcGFu\r\nIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsiPsK3WMK3PC9zcGFuPiDigJw8c3BhbiBzdHlsZT0iY29sb3I6\r\nICMxOTE5MTk7Ij5UZXN0IENyb3NzPC9zcGFuPuKAnSBieSBmYWNlYm9vazQgbmFtZSwgZW1haWwxIG5h\r\nbWUsIGVtYWlsMiBuYW1lLCBldGMuPC9hPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAg\r\nPC90cj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjEwIj48L3RkPjwvdHI+\r\nCiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1IiBzdHlsZT0iZm9u\r\ndC1zaXplOiAyMHB4OyBsaW5lLWhlaWdodDogMjZweDsiPgogICAgICAgICAgICAgICAgICAgIDxhIGhy\r\nZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJj\r\nb2xvcjojM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IGZvbnQtd2VpZ2h0OiBsaWdodGVyOyI+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIE5ldyBUaXRsZQogICAgICAgICAgICAgICAgICAgIDwvYT4K\r\nICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQg\r\nY29sc3Bhbj0iNSIgaGVpZ2h0PSIxMCI+PC90ZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAg\r\nICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSIxODAiPgogICAgICAgICAgICAgICAgICAgIAog\r\nICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJmb250LXNpemU6IDIwcHg7IGxpbmUtaGVpZ2h0OiAy\r\nNnB4OyBtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3Np\r\ndGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1\r\nOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij40OjQ1UE0gb24gVHVlLCBPY3QgMjMgMjAxMigrMDgwMCk8\r\nL2E+CiAgICAgICAgICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAg\r\nICAgICAgPC90ZD4KICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMTAiPjwvdGQ+CiAgICAgICAgICAg\r\nICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSIxOTAiIHN0eWxlPSJ3b3JkLWJyZWFrOiBicmVhay1h\r\nbGw7Ij4KICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0iZm9u\r\ndC1zaXplOiAyMHB4OyBsaW5lLWhlaWdodDogMjZweDsgbWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgIDxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFf\r\ndG9rZW4iIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+VGVzdCBQ\r\nbGFjZTE8L2E+CiAgICAgICAgICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgIDxwIHN0\r\neWxlPSJtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3Np\r\ndGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1\r\nOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij50ZXN0IHBsYWNlIDE8L2E+CiAgICAgICAgICAgICAgICAg\r\nICAgPC9wPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAg\r\nICAgICAgIDx0ZCB3aWR0aD0iMTAiPjwvdGQ+CiAgICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3Ai\r\nIHdpZHRoPSIyMTAiPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4KICAg\r\nICAgICAgICAgPC90cj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjEwIj48\r\nL3RkPjwvdHI+CiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1Ij4K\r\nICAgICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNw\r\nYWNpbmc9IjAiIHN0eWxlPSJmb250LWZhbWlseTogVmVyZGFuYTsgZm9udC1zaXplOiAxM3B4OyBsaW5l\r\nLWhlaWdodDogMjBweDsgY29sb3I6ICMxOTE5MTk7IGZvbnQtd2VpZ2h0OiBub3JtYWw7IHdpZHRoOiAx\r\nMDAlOyBiYWNrZ3JvdW5kLWNvbG9yOiAjZmJmYmZiOyI+CiAgICAgICAgICAgICAgICAgICAgCQogICAg\r\nICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQgd2lk\r\ndGg9IjE1Ij48aW1nIHNyYz0iaHR0cDovL3NpdGUvaW1nL2VtYWlsL3JzdnBfYWNjZXB0ZWRfMTJfYmx1\r\nZS5wbmciIC8+PC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGll\r\nbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMxOTE5MTk7IHRleHQtZGVjb3JhdGlvbjogbm9u\r\nZTsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0ibWFyZ2luOiAw\r\nOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8c3BhbiBjbGFzcz0iZXhm\r\nZV9tYWlsX2lkZW50aXR5X25hbWUiPjU8L3NwYW4+IEFjY2VwdGVkOiA8c3BhbiBjbGFzcz0nZXhmZV9t\r\nYWlsX2lkZW50aXR5X25hbWUnPmVtYWlsMiBuYW1lPC9zcGFuPiwgPHNwYW4gY2xhc3M9J2V4ZmVfbWFp\r\nbF9pZGVudGl0eV9uYW1lJz5mYWNlYm9vazUgbmFtZTwvc3Bhbj4sIGFuZCAxIG90aGVycy4KICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDwvYT4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAg\r\nICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgPHRkIHdpZHRoPSIxNSI+PGltZyBzcmM9Imh0dHA6Ly9zaXRlL2ltZy9lbWFpbC9yc3ZwX2RlY2xp\r\nbmVkXzEyLnBuZyIgLz48L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkPgogICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgIDxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJl\r\nY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJjb2xvcjogIzE5MTkxOTsgdGV4dC1kZWNvcmF0aW9u\r\nOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJtYXJn\r\naW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIFVuYXZhaWxhYmxl\r\nOiB0d2l0dGVyMyBuYW1lCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvcD4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L2E+CiAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAg\r\nICAKICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMTUiPjxpbWcgc3JjPSJodHRwOi8vc2l0\r\nZS9pbWcvZW1haWwvcGx1c18xMl9ibHVlLnBuZyIgLz48L3RkPgogICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgPHRkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxhIGhyZWY9Imh0dHA6Ly9z\r\naXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJjb2xvcjogIzE5MTkx\r\nOTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgIDxwIHN0eWxlPSJtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIE5ld2x5IGludml0ZWQ6IDxzcGFuIGNsYXNzPSJleGZlX21haWxfaWRlbnRpdHlfbmFtZSI+\r\nZmFjZWJvb2s1IG5hbWU8L3NwYW4+LgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8\r\nL3A+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICA8\r\ndHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQgd2lkdGg9IjE1Ij48aW1nIHNyYz0iaHR0\r\ncDovL3NpdGUvaW1nL2VtYWlsL21pbnVzXzEyLnBuZyIgLz48L3RkPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgPHRkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxhIGhyZWY9Imh0dHA6\r\nLy9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJjb2xvcjogIzE5\r\nMTkxOTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDxwIHN0eWxlPSJtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgIFJlbW92ZWQ6IGZhY2Vib29rNiBuYW1lLgogICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICA8L3A+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgPC90YWJsZT4KICAgICAgICAg\r\nICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQgY29sc3Bhbj0i\r\nNSIgaGVpZ2h0PSIxMCI+PC90ZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8\r\ndGQgY29sc3Bhbj0iNSI+CiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgCiAg\r\nICAgICAgICAgICAgICAgICAgPGltZyBzdHlsZT0icGFkZGluZy1yaWdodDogNXB4OyIgd2lkdGg9IjQw\r\nIiBoZWlnaHQ9IjQwIiBhbHQ9ImVtYWlsMiBuYW1lIiB0aXRsZT0iZW1haWwyIG5hbWUiIHNyYz0iaHR0\r\ncDovL3NpdGUvYXBpL3YyL2F2YXRhci9yZW5kZXI/cmVzb2x1dGlvbj0yeCZ1cmw9YUhSMGNEb3ZMM0Jo\r\nZEdndmRHOHZaVzFoYVd3eUxtRjJZWFJoY2clM0QlM0Qmd2lkdGg9NDAmaGVpZ2h0PTQwIj4KICAgICAg\r\nICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICA8aW1nIHN0eWxlPSJwYWRkaW5nLXJpZ2h0\r\nOiA1cHg7IiB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIGFsdD0idHdpdHRlcjMgbmFtZSIgdGl0bGU9InR3\r\naXR0ZXIzIG5hbWUiIHNyYz0iaHR0cDovL3NpdGUvYXBpL3YyL2F2YXRhci9yZW5kZXI/cmVzb2x1dGlv\r\nbj0yeCZ1cmw9YUhSMGNEb3ZMM0JoZEdndmRHOHZkSGRwZEhSbGNqTXVZWFpoZEdGeSZ3aWR0aD00MCZo\r\nZWlnaHQ9NDAmYWxwaGE9MC4zMyI+CiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAg\r\nICAgPGltZyBzdHlsZT0icGFkZGluZy1yaWdodDogNXB4OyIgd2lkdGg9IjQwIiBoZWlnaHQ9IjQwIiBh\r\nbHQ9ImZhY2Vib29rNCBuYW1lIiB0aXRsZT0iZmFjZWJvb2s0IG5hbWUiIHNyYz0iaHR0cDovL3NpdGUv\r\nYXBpL3YyL2F2YXRhci9yZW5kZXI/cmVzb2x1dGlvbj0yeCZ1cmw9YUhSMGNEb3ZMM0JoZEdndmRHOHZa\r\nbUZqWldKdmIyczBMbUYyWVhSaGNnJTNEJTNEJndpZHRoPTQwJmhlaWdodD00MCI+CiAgICAgICAgICAg\r\nICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgPGltZyBzdHlsZT0icGFkZGluZy1yaWdodDogNXB4\r\nOyIgd2lkdGg9IjQwIiBoZWlnaHQ9IjQwIiBhbHQ9InR3aXR0ZXIxIG5hbWUiIHRpdGxlPSJ0d2l0dGVy\r\nMSBuYW1lIiBzcmM9Imh0dHA6Ly9zaXRlL2FwaS92Mi9hdmF0YXIvcmVuZGVyP3Jlc29sdXRpb249Mngm\r\ndXJsPWFIUjBjRG92TDNCaGRHZ3ZkRzh2ZEhkcGRIUmxjakV1WVhaaGRHRnkmd2lkdGg9NDAmaGVpZ2h0\r\nPTQwJmFscGhhPTAuMzMmaXNob3N0PXRydWUmbWF0ZXM9MiI+CiAgICAgICAgICAgICAgICAgICAgCiAg\r\nICAgICAgICAgICAgICAgICAgPGltZyBzdHlsZT0icGFkZGluZy1yaWdodDogNXB4OyIgd2lkdGg9IjQw\r\nIiBoZWlnaHQ9IjQwIiBhbHQ9ImZhY2Vib29rNSBuYW1lIiB0aXRsZT0iZmFjZWJvb2s1IG5hbWUiIHNy\r\nYz0iaHR0cDovL3NpdGUvYXBpL3YyL2F2YXRhci9yZW5kZXI/cmVzb2x1dGlvbj0yeCZ1cmw9YUhSMGNE\r\nb3ZMM0JoZEdndmRHOHZabUZqWldKdmIyczFMbUYyWVhSaGNnJTNEJTNEJndpZHRoPTQwJmhlaWdodD00\r\nMCZtYXRlcz0yIj4KICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAg\r\nICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQgY29sc3Bhbj0iNSIgaGVpZ2h0PSIxMCI+PC90\r\nZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQgY29sc3Bhbj0iNSI+CiAg\r\nICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50\r\nX2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7\r\nIj50ZXN0IGNyb3NzIGRlc2NyaXB0aW9uPC9hPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAg\r\nICAgPC90cj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjIwIj48L3RkPjwv\r\ndHI+CiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1IiBzdHlsZT0i\r\nZm9udC1zaXplOiAxMXB4OyBsaW5lLWhlaWdodDogMTVweDsgY29sb3I6ICM3RjdGN0Y7Ij4KICAgICAg\r\nICAgICAgICAgICAgICBSZXBseSB0aGlzIGVtYWlsIGRpcmVjdGx5IGFzIGNvbnZlcnNhdGlvbiwgb3Ig\r\ndHJ5IDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJlZj0i\r\naHR0cDovL2FwcC91cmwiPkVYRkU8L2E+IGFwcC4KICAgICAgICAgICAgICAgICAgICA8YnIgLz4KICAg\r\nICAgICAgICAgICAgICAgICA8c3BhbiBzdHlsZT0iY29sb3I6ICNCMkIyQjIiPlRoaXMgdXBkYXRlIGlz\r\nIHNlbnQgZnJvbSA8YSBzdHlsZT0iY29sb3I6ICMzYTZlYTU7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsi\r\nIGhyZWY9Imh0dHA6Ly9zaXRlL3VybCI+RVhGRTwvYT4gYXV0b21hdGljYWxseS4gPGEgc3R5bGU9ImNv\r\nbG9yOiAjRTZFNkU2OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8vc2l0ZS91cmwv\r\nbXV0ZS9jcm9zcz90b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIj5VbnN1YnNjcmliZT88L2E+CiAg\r\nICAgICAgICAgICAgICAgICAgPCEtLQogICAgICAgICAgICAgICAgICAgIFlvdSBjYW4gY2hhbmdlCiAg\r\nICAgICAgICAgICAgICAgICAgPGEgc3R5bGU9ImNvbG9yOiAjQjJCMkIyOyB0ZXh0LWRlY29yYXRpb246\r\nIHVuZGVsaW5lOyIgaHJlZj0iIj5ub3RpZmljYXRpb24gcHJlZmVyZW5jZTwvYT4uCiAgICAgICAgICAg\r\nICAgICAgICAgLS0+CiAgICAgICAgICAgICAgICAgICAgPC9zcGFuPgogICAgICAgICAgICAgICAgPC90\r\nZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICA8L3RhYmxlPgogICAgPC9ib2R5Pgo8L2h0bWw+Cg==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6TmV3IFRpdGxlDQpYLVdSLUNBTERFU0M6ZXhmZSBjYWwNClgtV1ItVElNRVpPTkU6\r\nKzA4MDANCkJFR0lOOlZFVkVOVA0KVUlEOiExMjNAZXhmZQ0KRFRTVEFNUDoyMDEyMTAyM1QwODQ1MDBa\r\nDQpERVNDUklQVElPTjp0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uDQpEVFNUQVJUOjIwMTIxMDIzVDA4NDUw\r\nMFoNCkxPQ0FUSU9OOlRlc3QgUGxhY2UxXG50ZXN0IHBsYWNlIDENClNVTU1BUlk6TmV3IFRpdGxlDQpV\r\nUkw6aHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbg0KRU5EOlZFVkVO\r\nVA0KRU5EOlZDQUxFTkRBUg0K\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
	assert.Equal(t, text, expect)
}

func TestCrossInvitationEmail(t *testing.T) {
	cross1 := cross
	cross1.Time = &time1
	cross1.Exfee = exfee1

	inv := model.CrossInvitation{}
	inv.To = remail1
	inv.Cross = cross1

	c := NewCross(localTemplate, &config, nil)
	text, err := c.getInvitationContent(inv)
	assert.Equal(t, err, nil)
	t.Logf("text:---------start---------\n%s\n---------end----------", text)
	expect := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <+123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <to_email_address>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <+123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nWW91J3JlIGdhdGhlcmluZyB0aGlzIMK3WMK3LgoKClRlc3QgQ3Jvc3MKPT09PT09PQpodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuCgo0OjQ1UE0gb24gVHVlLCBPY3QgMjMg\r\nMjAxMgo9PT09PT09CgpQbGFjZQo9PT09PT09CiAgVG8gYmUgZGVjaWRlZC4KCgpJJ20gaW4uIENoZWNr\r\nIGl0IG91dDogaHR0cDovL3NpdGUvdXJsLz90b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuJnJzdnA9\r\nYWNjZXB0CgoKNiBJbnZpdGVkOgrCtyBlbWFpbDEgbmFtZSAoSG9zdCkgd2l0aCAyIHBlb3BsZQrCtyBl\r\nbWFpbDIgbmFtZQrCtyB0d2l0dGVyMyBuYW1lCsK3IGZhY2Vib29rNCBuYW1lCgoKRGVzY3JpcHRpb24K\r\nLS0tLS0tLQogIHRlc3QgY3Jvc3MgZGVzY3JpcHRpb24KCgojIFJlcGx5IHRoaXMgZW1haWwgZGlyZWN0\r\nbHkgYXMgY29udmVyc2F0aW9uLiAj\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9sYWJlbCB7\r\nCiAgICAgICAgICAgICAgICBiYWNrZ3JvdW5kLWNvbG9yOiAjRDVFOEYyOwogICAgICAgICAgICAgICAg\r\nY29sb3I6ICMzYTZlYTU7CiAgICAgICAgICAgICAgICBmb250LXNpemU6IDExcHg7CiAgICAgICAgICAg\r\nICAgICBwYWRkaW5nOiAwIDJweCAwIDJweDsKICAgICAgICAgICAgfQogICAgICAgICAgICAuZXhmZV9t\r\nYWlsX21hdGVzIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICAgICAg\r\nZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICB9CiAgICAgICAgICAgIC5leGZlX21haWxfaWRlbnRp\r\ndHkgewogICAgICAgICAgICAgICAgZm9udC1zdHlsZTogaXRhbGljOwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfaWRlbnRpdHlfbmFtZSB7CiAgICAgICAgICAgICAgICBjb2xvcjogIzE5\r\nMTkxOTsKICAgICAgICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4K\r\nICAgICAgICA8dGFibGUgd2lkdGg9IjY0MCIgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNw\r\nYWNpbmc9IjAiIHN0eWxlPSJmb250LWZhbWlseTogSGVsdmV0aWNhOyBmb250LXNpemU6IDEzcHg7IGxp\r\nbmUtaGVpZ2h0OiAxOXB4OyBjb2xvcjogIzE5MTkxOTsgZm9udC13ZWlnaHQ6IG5vcm1hbDsgcGFkZGlu\r\nZzogMzBweCA0MHB4IDMwcHggNDBweDsgYmFja2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsgbWluLWhlaWdo\r\ndDogNTYycHg7Ij4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjMi\r\nIHZhbGlnbj0idG9wIiBzdHlsZT0iZm9udC1zaXplOiAzMnB4OyBsaW5lLWhlaWdodDogMzhweDsgcGFk\r\nZGluZy1ib3R0b206IDE4cHg7Ij4KICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMzYTZlYTU7\r\nIHRleHQtZGVjb3JhdGlvbjogbm9uZTsgZm9udC13ZWlnaHQ6IDMwMDsiPgogICAgICAgICAgICAgICAg\r\nICAgICAgICBUZXN0IENyb3NzCiAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAg\r\nPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRk\r\nIHdpZHRoPSIzNDAiIHN0eWxlPSJ2ZXJ0aWNhbC1hbGlnbjogYmFzZWxpbmU7IGZvbnQtd2VpZ2h0OiAz\r\nMDA7Ij4KICAgICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIg\r\nY2VsbHNwYWNpbmc9IjAiPgogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHN0eWxlPSJwYWRkaW5nLWJvdHRvbTogMjBweDsg\r\nZm9udC1zaXplOiAyMHB4OyB2ZXJ0aWNhbC1hbGlnbjogYmFzZWxpbmU7Ij4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBZb3UncmUgZ2F0\r\naGVyaW5nIHRoaXMgPHNwYW4gc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5v\r\nbmU7Ij7Ct1jCtzwvc3Bhbj4uCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAg\r\nICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2VsbHBhZGRpbmc9\r\nIjAiIGNlbGxzcGFjaW5nPSIwIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRy\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIiB3\r\naWR0aD0iMTYwIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBo\r\ncmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0i\r\ndGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgCQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9ImZv\r\nbnQtc2l6ZTogMjBweDsgbGluZS1oZWlnaHQ6IDI2cHg7IG1hcmdpbjogMDsgY29sb3I6ICMzMzMzMzM7\r\nIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDQ6NDVQ\r\nTSBvbiBUdWUsIE9jdCAyMyAyMDEyKCswODAwKQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8L3A+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvYT4K\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHN0eWxlPSJwYWRkaW5nLWxl\r\nZnQ6IDEwcHg7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBo\r\ncmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0i\r\ndGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICA8cCBzdHlsZT0iZm9udC1zaXplOiAyMHB4OyBsaW5lLWhlaWdodDogMjZweDsgbWFyZ2luOiAwOyBj\r\nb2xvcjogIzMzMzMzMzsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgUGxhY2UKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgPC9wPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8cCBz\r\ndHlsZT0iY29sb3I6ICMxOTE5MTk7IG1hcmdpbjogMDsiPgogICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgVG8gYmUgZGVjaWRlZAogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3A+IAogICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICA8L2E+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgIDwvdGFibGU+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAg\r\nICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAg\r\nICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0i\r\ndG9wIiBzdHlsZT0icGFkZGluZy10b3A6IDMwcHg7IHBhZGRpbmctYm90dG9tOiAzMHB4OyI+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgPGEgc3R5bGU9ImZsb2F0OiBsZWZ0OyBkaXNwbGF5OiBi\r\nbG9jazsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyBib3JkZXI6IDFweCBzb2xpZCAjYmViZWJlOyBiYWNr\r\nZ3JvdW5kLWNvbG9yOiAjM0E2RUE1OyBjb2xvcjogI0ZGRkZGRjsgcGFkZGluZzogNXB4IDMwcHggNXB4\r\nIDMwcHg7IG1hcmdpbi1sZWZ0OiAyNXB4OyIgYWx0PSJBY2NlcHQiIGhyZWY9Imh0dHA6Ly9zaXRlL3Vy\r\nbC8/dG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiZyc3ZwPWFjY2VwdCI+SSdtIGluPC9hPgogICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxhIHN0eWxlPSJmbG9hdDogbGVmdDsgZGlzcGxheTog\r\nYmxvY2s7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsgYm9yZGVyOiAxcHggc29saWQgI2JlYmViZTsgYmFj\r\na2dyb3VuZC1jb2xvcjogI0U2RTZFNjsgY29sb3I6ICMxOTE5MTk7IHBhZGRpbmc6IDVweCAyNXB4IDVw\r\neCAyNXB4OyBtYXJnaW4tbGVmdDogMTVweDsiIGFsdD0iQ2hlY2sgaXQgb3V0IiBocmVmPSJodHRwOi8v\r\nc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIj5DaGVjayBpdCBvdXQuLi48L2E+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8\r\nL3RyPgogICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICB0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uCiAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAg\r\nICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgIDwvdGFibGU+CiAgICAgICAgICAgICAg\r\nICA8L3RkPgogICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIzMCI+PC90ZD4KICAgICAgICAgICAgICAg\r\nIDx0ZCB2YWxpZ249InRvcCI+CiAgICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2Vs\r\nbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIj4KICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQg\r\naGVpZ2h0PSI2OCIgdmFsaWduPSJ0b3AiIGFsaWduPSJyaWdodCI+CiAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAg\r\nICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVy\r\nPSIwIiBzdHlsZT0iY29sb3I6ICMzMzMzMzM7IiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAi\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nPHRkIHdpZHRoPSIyNSIgaGVpZ2h0PSIyNSIgYWxpZ249ImxlZnQiIHZhbGlnbj0idG9wIj4KICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8aW1nIHdpZHRoPSIyMCIgaGVpZ2h0\r\nPSIyMCIgdGl0bGU9ImVtYWlsMSBuYW1lIiBhbHQ9ImVtYWlsMSBuYW1lIiBzcmM9Imh0dHA6Ly9wYXRo\r\nL3RvL2VtYWlsMS5hdmF0YXIiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHNwYW4+ZW1haWwxIG5hbWU8L3NwYW4+IDxz\r\ncGFuIGNsYXNzPSJleGZlX21haWxfbWF0ZXMiPisyPC9zcGFuPiA8c3BhbiBjbGFzcz0iZXhmZV9tYWls\r\nX2xhYmVsIj5ob3N0PC9zcGFuPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMjUi\r\nIGhlaWdodD0iMjUiIGFsaWduPSJsZWZ0IiB2YWxpZ249InRvcCI+CiAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgPGltZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIHRpdGxlPSJl\r\nbWFpbDIgbmFtZSIgYWx0PSJlbWFpbDIgbmFtZSIgc3JjPSJodHRwOi8vcGF0aC90by9lbWFpbDIuYXZh\r\ndGFyIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgIDxzcGFuPmVtYWlsMiBuYW1lPC9zcGFuPgogICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgIDx0ZCB3aWR0aD0iMjUiIGhlaWdodD0iMjUiIGFsaWduPSJsZWZ0IiB2YWxpZ249InRv\r\ncCI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPGltZyB3aWR0aD0i\r\nMjAiIGhlaWdodD0iMjAiIHRpdGxlPSJ0d2l0dGVyMyBuYW1lIiBhbHQ9InR3aXR0ZXIzIG5hbWUiIHNy\r\nYz0iaHR0cDovL3BhdGgvdG8vdHdpdHRlcjMuYXZhdGFyIj4KICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICA8dGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxzcGFuPnR3aXR0\r\nZXIzIG5hbWU8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3Rk\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRy\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIyNSIgaGVp\r\nZ2h0PSIyNSIgYWxpZ249ImxlZnQiIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8aW1nIHdpZHRoPSIyMCIgaGVpZ2h0PSIyMCIgdGl0bGU9ImZhY2Vi\r\nb29rNCBuYW1lIiBhbHQ9ImZhY2Vib29rNCBuYW1lIiBzcmM9Imh0dHA6Ly9wYXRoL3RvL2ZhY2Vib29r\r\nNC5hdmF0YXIiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgPHNwYW4+ZmFjZWJvb2s0IG5hbWU8L3NwYW4+CiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RhYmxlPgogICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAg\r\nICA8L3RhYmxlPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICAg\r\nICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjMiIHN0eWxlPSJmb250LXNpemU6IDEx\r\ncHg7IGxpbmUtaGVpZ2h0OiAxNXB4OyBjb2xvcjogIzdGN0Y3RjsgcGFkZGluZy10b3A6IDQwcHg7Ij4K\r\nICAgICAgICAgICAgICAgICAgICBSZXBseSB0aGlzIGVtYWlsIGRpcmVjdGx5IGFzIGNvbnZlcnNhdGlv\r\nbiwgb3IgVHJ5IDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIg\r\naHJlZj0iaHR0cDovL2FwcC91cmwiPkVYRkU8L2E+IGFwcC4KICAgICAgICAgICAgICAgICAgICA8YnIg\r\nLz4KICAgICAgICAgICAgICAgICAgICBUaGlzIDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1k\r\nZWNvcmF0aW9uOiBub25lOyIgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2Vt\r\nYWlsMV90b2tlbiI+wrdYwrc8L2E+IGludml0YXRpb24gaXMgc2VudCBieSA8c3BhbiBjbGFzcz0iZXhm\r\nZV9tYWlsX2lkZW50aXR5X25hbWUiPmVtYWlsMSBuYW1lPC9zcGFuPiBmcm9tIDxhIHN0eWxlPSJjb2xv\r\ncjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJlZj0iaHR0cDovL3NpdGUvdXJsIj5F\r\nWEZFPC9hPi4KICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgPC90\r\nYWJsZT4KICAgIDwvYm9keT4KPC9odG1sPgo=\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6VGVzdCBDcm9zcw0KWC1XUi1DQUxERVNDOmV4ZmUgY2FsDQpYLVdSLVRJTUVaT05F\r\nOiswODAwDQpCRUdJTjpWRVZFTlQNClVJRDohMTIzQGV4ZmUNCkRUU1RBTVA6MjAxMjEwMjNUMDg0NTAw\r\nWg0KREVTQ1JJUFRJT046dGVzdCBjcm9zcyBkZXNjcmlwdGlvbg0KRFRTVEFSVDoyMDEyMTAyM1QwODQ1\r\nMDBaDQpMT0NBVElPTjoNClNVTU1BUlk6VGVzdCBDcm9zcw0KVVJMOmh0dHA6Ly9zaXRlL3VybC8jIXRv\r\na2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4NCkVORDpWRVZFTlQNCkVORDpWQ0FMRU5EQVINCg==\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
	assert.Equal(t, text, expect)
}

func TestCrossInvitationTwitter(t *testing.T) {
	cross1 := cross
	cross1.Time = &time1
	cross1.Exfee = exfee1

	inv := model.CrossInvitation{}
	inv.To = rtwitter1
	inv.Cross = cross1

	c := NewCross(localTemplate, &config, nil)
	text, err := c.getInvitationContent(inv)
	assert.Equal(t, err, nil)
	t.Logf("text:---------start---------\n%s\n---------end----------", text)
	expect := ""
	assert.Equal(t, text, expect)
}
