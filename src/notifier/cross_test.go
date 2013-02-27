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
	expectPrivate := "\n\n\n\n\n\n\n\n\\(\"Test Cross\"\\) update: \\(\"New Title\"\\). 4:45PM on Tue, Oct 23 2012 at \\(Test Place1\\). 5 people invited. http://site/url/#!token=recipient_twitter1_token\n\n\\(facebook5 name\\) is invited to \\(\"New Title\"\\) by facebook4 name, email1 name, etc. http://site/url/#!token=recipient_twitter1_token\n\n\\(facebook6 name\\) left \\(\"New Title\"\\). http://site/url/#!token=recipient_twitter1_token\n\n\n\n\n\n\n\n\\(email2 name\\) and \\(facebook5 name\\) accepted \\(\"New Title\"\\), \\(twitter3 name\\) is unavailable, 5 of 9 accepted. http://site/url/#!token=recipient_twitter1_token\n\n\n\n\n"
	assert.Equal(t, private, expectPrivate)
	t.Logf("private:-----start------\n%s\n-------end-------", private)
	expectPublic := `Updates: http://site/url/#!123/ecip (Please follow @EXFE to receive details PRIVATELY through Direct Message.)`
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
	private, _, err := c.getSummaryContent(updates)
	assert.Equal(t, err, nil)
	t.Logf("private:-----start------\n%s\n-------end-------", private)
	expectPrivate := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <x+e123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <to_email_address>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x+e123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nVXBkYXRlcyBvZiDCt1jCtyDigJxUZXN0IENyb3Nz4oCdIGJ5IGZhY2Vib29rNCBuYW1lLCBlbWFpbDEg\r\nbmFtZSwgZW1haWwyIG5hbWUsIGV0Yy4KCipOZXcgVGl0bGUqCj09PT09PT0KaHR0cDovL3NpdGUvdXJs\r\nLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbgoKKjQ6NDVQTSBvbiBUdWUsIE9jdCAyMyAyMDEy\r\nKgo9PT09PT09CgoqVGVzdCBQbGFjZTEqCj09PT09PT0KICAqdGVzdCBwbGFjZSAxKgoKCsK3IDUgQWNj\r\nZXB0ZWQ6IGVtYWlsMiBuYW1lLCBmYWNlYm9vazUgbmFtZSBhbmQgMSBvdGhlcnMuCsK3IFVuYXZhaWxh\r\nYmxlOiB0d2l0dGVyMyBuYW1lLgrCtyBOZXdseSBpbnZpdGVkOiBmYWNlYm9vazUgbmFtZS4KwrcgUmVt\r\nb3ZlZDogZmFjZWJvb2s2IG5hbWUuCgojIFJlcGx5IHRoaXMgZW1haWwgZGlyZWN0bHkgYXMgY29udmVy\r\nc2F0aW9uLiAj\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9pZGVudGl0\r\neV9uYW1lIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfbXNnX2lkZW50aXR5X25hbWUgewogICAgICAgICAgICAgICAgY29sb3I6\r\nICM2NjY2NjY7CiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfYXQgewogICAg\r\nICAgICAgICAgICAgZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICAgICAgY29sb3I6ICM5OTk5OTk7\r\nCiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfdGltZSB7CiAgICAgICAgICAg\r\nICAgICBmb250LXNpemU6IDEycHg7CiAgICAgICAgICAgICAgICBjb2xvcjogIzY2NjY2NjsKICAgICAg\r\nICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4KICAgICAgICA8dGFi\r\nbGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAiIHN0eWxlPSJmb250LWZh\r\nbWlseTogVmVyZGFuYTsgZm9udC1zaXplOiAxM3B4OyBsaW5lLWhlaWdodDogMjBweDsgY29sb3I6ICMx\r\nOTE5MTk7IGZvbnQtd2VpZ2h0OiBub3JtYWw7IHdpZHRoOiA2NDBweDsgcGFkZGluZzogMjBweDsgYmFj\r\na2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsiPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8\r\ndGQgY29sc3Bhbj0iNSIgc3R5bGU9ImNvbG9yOiAjMzMzMzMzOyI+CiAgICAgICAgICAgICAgICAgICAg\r\nPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5\r\nbGU9ImNvbG9yOiAjMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5VcGRhdGVzIG9mIDxzcGFu\r\nIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsiPsK3WMK3PC9zcGFuPiDigJw8c3BhbiBzdHlsZT0iY29sb3I6\r\nICMxOTE5MTk7Ij5UZXN0IENyb3NzPC9zcGFuPuKAnSBieSBmYWNlYm9vazQgbmFtZSwgZW1haWwxIG5h\r\nbWUsIGVtYWlsMiBuYW1lLCBldGMuPC9hPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAg\r\nPC90cj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjEwIj48L3RkPjwvdHI+\r\nCiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1IiBzdHlsZT0iZm9u\r\ndC1zaXplOiAyMHB4OyBsaW5lLWhlaWdodDogMjZweDsiPgogICAgICAgICAgICAgICAgICAgIDxhIGhy\r\nZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJj\r\nb2xvcjojM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IGZvbnQtd2VpZ2h0OiBsaWdodGVyOyI+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIE5ldyBUaXRsZQogICAgICAgICAgICAgICAgICAgIDwvYT4K\r\nICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQg\r\nY29sc3Bhbj0iNSIgaGVpZ2h0PSIxMCI+PC90ZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAg\r\nICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSIxODAiPgogICAgICAgICAgICAgICAgICAgIAog\r\nICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJmb250LXNpemU6IDIwcHg7IGxpbmUtaGVpZ2h0OiAy\r\nNnB4OyBtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3Np\r\ndGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1\r\nOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij40OjQ1UE0gb24gVHVlLCBPY3QgMjMgMjAxMjwvYT4KICAg\r\nICAgICAgICAgICAgICAgICA8L3A+CiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICA8\r\nL3RkPgogICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIxMCI+PC90ZD4KICAgICAgICAgICAgICAgIDx0\r\nZCB2YWxpZ249InRvcCIgd2lkdGg9IjE5MCIgc3R5bGU9IndvcmQtYnJlYWs6IGJyZWFrLWFsbDsiPgog\r\nICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJmb250LXNpemU6\r\nIDIwcHg7IGxpbmUtaGVpZ2h0OiAyNnB4OyBtYXJnaW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAg\r\nICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIg\r\nc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5UZXN0IFBsYWNlMTwv\r\nYT4KICAgICAgICAgICAgICAgICAgICA8L3A+CiAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9Im1h\r\ncmdpbjogMDsiPgogICAgICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0ZS91cmwv\r\nIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMzYTZlYTU7IHRleHQt\r\nZGVjb3JhdGlvbjogbm9uZTsiPnRlc3QgcGxhY2UgMTwvYT4KICAgICAgICAgICAgICAgICAgICA8L3A+\r\nCiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAg\r\nPHRkIHdpZHRoPSIxMCI+PC90ZD4KICAgICAgICAgICAgICAgIDx0ZCB2YWxpZ249InRvcCIgd2lkdGg9\r\nIjIxMCI+CiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAg\r\nICA8L3RyPgogICAgICAgICAgICA8dHI+PHRkIGNvbHNwYW49IjUiIGhlaWdodD0iMTAiPjwvdGQ+PC90\r\ncj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjUiPgogICAgICAg\r\nICAgICAgICAgICAgIDx0YWJsZSBib3JkZXI9IjAiIGNlbGxwYWRkaW5nPSIwIiBjZWxsc3BhY2luZz0i\r\nMCIgc3R5bGU9ImZvbnQtZmFtaWx5OiBWZXJkYW5hOyBmb250LXNpemU6IDEzcHg7IGxpbmUtaGVpZ2h0\r\nOiAyMHB4OyBjb2xvcjogIzE5MTkxOTsgZm9udC13ZWlnaHQ6IG5vcm1hbDsgd2lkdGg6IDEwMCU7IGJh\r\nY2tncm91bmQtY29sb3I6ICNmYmZiZmI7Ij4KICAgICAgICAgICAgICAgICAgICAJCiAgICAgICAgICAg\r\nICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMTUi\r\nPjxpbWcgc3JjPSJodHRwOi8vc2l0ZS9pbWcvZW1haWwvcnN2cF9hY2NlcHRlZF8xMl9ibHVlLnBuZyIg\r\nLz48L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgIDxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFp\r\nbDFfdG9rZW4iIHN0eWxlPSJjb2xvcjogIzE5MTkxOTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJtYXJnaW46IDA7Ij4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxzcGFuIGNsYXNzPSJleGZlX21haWxf\r\naWRlbnRpdHlfbmFtZSI+NTwvc3Bhbj4gQWNjZXB0ZWQ6IDxzcGFuIGNsYXNzPSdleGZlX21haWxfaWRl\r\nbnRpdHlfbmFtZSc+ZW1haWwyIG5hbWU8L3NwYW4+LCA8c3BhbiBjbGFzcz0nZXhmZV9tYWlsX2lkZW50\r\naXR5X25hbWUnPmZhY2Vib29rNSBuYW1lPC9zcGFuPiwgYW5kIDEgb3RoZXJzLgogICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICA8L3A+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nPC9hPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAg\r\nICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgIAog\r\nICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQg\r\nd2lkdGg9IjE1Ij48aW1nIHNyYz0iaHR0cDovL3NpdGUvaW1nL2VtYWlsL3JzdnBfZGVjbGluZWRfMTIu\r\ncG5nIiAvPjwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50\r\nX2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjMTkxOTE5OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7\r\nIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9Im1hcmdpbjogMDsi\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgVW5hdmFpbGFibGU6IHR3aXR0\r\nZXIzIG5hbWUKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9wPgogICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgIDwvYT4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgICAgIAogICAg\r\nICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIxNSI+PGltZyBzcmM9Imh0dHA6Ly9zaXRlL2ltZy9l\r\nbWFpbC9wbHVzXzEyX2JsdWUucG5nIiAvPjwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8\r\ndGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJs\r\nLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjMTkxOTE5OyB0ZXh0\r\nLWRlY29yYXRpb246IG5vbmU7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHAg\r\nc3R5bGU9Im1hcmdpbjogMDsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nTmV3bHkgaW52aXRlZDogPHNwYW4gY2xhc3M9ImV4ZmVfbWFpbF9pZGVudGl0eV9uYW1lIj5mYWNlYm9v\r\nazUgbmFtZTwvc3Bhbj4uCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvcD4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L2E+CiAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAg\r\nICAKICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMTUiPjxpbWcgc3JjPSJodHRwOi8vc2l0\r\nZS9pbWcvZW1haWwvbWludXNfMTIucG5nIiAvPjwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICA8dGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUv\r\ndXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiIgc3R5bGU9ImNvbG9yOiAjMTkxOTE5OyB0\r\nZXh0LWRlY29yYXRpb246IG5vbmU7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nPHAgc3R5bGU9Im1hcmdpbjogMDsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgUmVtb3ZlZDogZmFjZWJvb2s2IG5hbWUuCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgIDwvcD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L2E+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAg\r\nICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICA8L3RhYmxlPgogICAgICAgICAgICAgICAg\r\nPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWln\r\naHQ9IjEwIj48L3RkPjwvdHI+CiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xz\r\ncGFuPSI1Ij4KICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAKICAgICAgICAg\r\nICAgICAgICAgICA8aW1nIHN0eWxlPSJwYWRkaW5nLXJpZ2h0OiA1cHg7IiB3aWR0aD0iNDAiIGhlaWdo\r\ndD0iNDAiIGFsdD0iZW1haWwyIG5hbWUiIHRpdGxlPSJlbWFpbDIgbmFtZSIgc3JjPSJodHRwOi8vc2l0\r\nZS9hcGkvdjIvYXZhdGFyL3JlbmRlcj9yZXNvbHV0aW9uPTJ4JnVybD1hSFIwY0RvdkwzQmhkR2d2ZEc4\r\ndlpXMWhhV3d5TG1GMllYUmhjZyUzRCUzRCZ3aWR0aD00MCZoZWlnaHQ9NDAiPgogICAgICAgICAgICAg\r\nICAgICAgIAogICAgICAgICAgICAgICAgICAgIDxpbWcgc3R5bGU9InBhZGRpbmctcmlnaHQ6IDVweDsi\r\nIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgYWx0PSJ0d2l0dGVyMyBuYW1lIiB0aXRsZT0idHdpdHRlcjMg\r\nbmFtZSIgc3JjPSJodHRwOi8vc2l0ZS9hcGkvdjIvYXZhdGFyL3JlbmRlcj9yZXNvbHV0aW9uPTJ4JnVy\r\nbD1hSFIwY0RvdkwzQmhkR2d2ZEc4dmRIZHBkSFJsY2pNdVlYWmhkR0Z5JndpZHRoPTQwJmhlaWdodD00\r\nMCZhbHBoYT0wLjMzIj4KICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICA8aW1n\r\nIHN0eWxlPSJwYWRkaW5nLXJpZ2h0OiA1cHg7IiB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIGFsdD0iZmFj\r\nZWJvb2s0IG5hbWUiIHRpdGxlPSJmYWNlYm9vazQgbmFtZSIgc3JjPSJodHRwOi8vc2l0ZS9hcGkvdjIv\r\nYXZhdGFyL3JlbmRlcj9yZXNvbHV0aW9uPTJ4JnVybD1hSFIwY0RvdkwzQmhkR2d2ZEc4dlptRmpaV0p2\r\nYjJzMExtRjJZWFJoY2clM0QlM0Qmd2lkdGg9NDAmaGVpZ2h0PTQwIj4KICAgICAgICAgICAgICAgICAg\r\nICAKICAgICAgICAgICAgICAgICAgICA8aW1nIHN0eWxlPSJwYWRkaW5nLXJpZ2h0OiA1cHg7IiB3aWR0\r\naD0iNDAiIGhlaWdodD0iNDAiIGFsdD0idHdpdHRlcjEgbmFtZSIgdGl0bGU9InR3aXR0ZXIxIG5hbWUi\r\nIHNyYz0iaHR0cDovL3NpdGUvYXBpL3YyL2F2YXRhci9yZW5kZXI/cmVzb2x1dGlvbj0yeCZ1cmw9YUhS\r\nMGNEb3ZMM0JoZEdndmRHOHZkSGRwZEhSbGNqRXVZWFpoZEdGeSZ3aWR0aD00MCZoZWlnaHQ9NDAmYWxw\r\naGE9MC4zMyZpc2hvc3Q9dHJ1ZSZtYXRlcz0yIj4KICAgICAgICAgICAgICAgICAgICAKICAgICAgICAg\r\nICAgICAgICAgICA8aW1nIHN0eWxlPSJwYWRkaW5nLXJpZ2h0OiA1cHg7IiB3aWR0aD0iNDAiIGhlaWdo\r\ndD0iNDAiIGFsdD0iZmFjZWJvb2s1IG5hbWUiIHRpdGxlPSJmYWNlYm9vazUgbmFtZSIgc3JjPSJodHRw\r\nOi8vc2l0ZS9hcGkvdjIvYXZhdGFyL3JlbmRlcj9yZXNvbHV0aW9uPTJ4JnVybD1hSFIwY0RvdkwzQmhk\r\nR2d2ZEc4dlptRmpaV0p2YjJzMUxtRjJZWFJoY2clM0QlM0Qmd2lkdGg9NDAmaGVpZ2h0PTQwJm1hdGVz\r\nPTIiPgogICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAg\r\nPC90cj4KICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSI1IiBoZWlnaHQ9IjEwIj48L3RkPjwvdHI+\r\nCiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSI1Ij4KICAgICAgICAg\r\nICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwx\r\nX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMzMzMzMzM7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPnRlc3Qg\r\nY3Jvc3MgZGVzY3JpcHRpb248L2E+CiAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICA8L3Ry\r\nPgogICAgICAgICAgICA8dHI+PHRkIGNvbHNwYW49IjUiIGhlaWdodD0iMjAiPjwvdGQ+PC90cj4KICAg\r\nICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjUiIHN0eWxlPSJmb250LXNp\r\nemU6IDExcHg7IGxpbmUtaGVpZ2h0OiAxNXB4OyBjb2xvcjogIzdGN0Y3RjsiPgogICAgICAgICAgICAg\r\nICAgICAgIFJlcGx5IHRoaXMgZW1haWwgZGlyZWN0bHkgYXMgY29udmVyc2F0aW9uLCBvciB0cnkgPGEg\r\nc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8v\r\nYXBwL3VybCI+RVhGRTwvYT4gYXBwLgogICAgICAgICAgICAgICAgICAgIDxiciAvPgogICAgICAgICAg\r\nICAgICAgICAgIDxzcGFuIHN0eWxlPSJjb2xvcjogI0IyQjJCMiI+VGhpcyB1cGRhdGUgaXMgc2VudCBm\r\ncm9tIDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIgaHJlZj0i\r\naHR0cDovL3NpdGUvdXJsIj5FWEZFPC9hPiBhdXRvbWF0aWNhbGx5LiA8YSBzdHlsZT0iY29sb3I6ICNF\r\nNkU2RTY7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC9tdXRlL2Ny\r\nb3NzP3Rva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iPlVuc3Vic2NyaWJlPzwvYT4KICAgICAgICAg\r\nICAgICAgICAgICA8IS0tCiAgICAgICAgICAgICAgICAgICAgWW91IGNhbiBjaGFuZ2UKICAgICAgICAg\r\nICAgICAgICAgICA8YSBzdHlsZT0iY29sb3I6ICNCMkIyQjI7IHRleHQtZGVjb3JhdGlvbjogdW5kZWxp\r\nbmU7IiBocmVmPSIiPm5vdGlmaWNhdGlvbiBwcmVmZXJlbmNlPC9hPi4KICAgICAgICAgICAgICAgICAg\r\nICAtLT4KICAgICAgICAgICAgICAgICAgICA8L3NwYW4+CiAgICAgICAgICAgICAgICA8L3RkPgogICAg\r\nICAgICAgICA8L3RyPgogICAgICAgIDwvdGFibGU+CiAgICA8L2JvZHk+CjwvaHRtbD4K\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?TmV3IFRpdGxlLmljcw==?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6TmV3IFRpdGxlDQpYLVdSLUNBTERFU0M6ZXhmZSBjYWwNClgtV1ItVElNRVpPTkU6\r\nKzA4MDANCkJFR0lOOlZFVkVOVA0KVUlEOiExMjNAZXhmZQ0KRFRTVEFNUDoyMDEyMTAyM1QwODQ1MDBa\r\nDQpERVNDUklQVElPTjp0ZXN0IGNyb3NzIGRlc2NyaXB0aW9uDQpEVFNUQVJUOjIwMTIxMDIzVDA4NDUw\r\nMFoNCkxPQ0FUSU9OOlRlc3QgUGxhY2UxXG50ZXN0IHBsYWNlIDENClNVTU1BUlk6TmV3IFRpdGxlDQpV\r\nUkw6aHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbg0KRU5EOlZFVkVO\r\nVA0KRU5EOlZDQUxFTkRBUg0K\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
	assert.Equal(t, private, expectPrivate)
}

func TestCrossInvitationEmail(t *testing.T) {
	cross1 := cross
	cross1.Time = time1
	cross1.Exfee = exfee1

	inv := model.CrossInvitation{}
	inv.To = remail1
	inv.Cross = cross1

	c := NewCross(localTemplate, &config, nil)
	private, _, err := c.getInvitationContent(inv)
	assert.Equal(t, err, nil)
	t.Logf("private:---------start---------\n%s\n---------end----------", private)
	expectPrivate := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <x+e123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <to_email_address>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x+e123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nWW91J3JlIGdhdGhlcmluZyB0aGlzIMK3WMK3LgoKClRlc3QgQ3Jvc3MKPT09PT09PQpodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuCgo0OjQ1UE0gb24gVHVlLCBPY3QgMjMg\r\nMjAxMgo9PT09PT09CgpQbGFjZQo9PT09PT09CiAgVG8gYmUgZGVjaWRlZC4KCgpJJ20gaW4uIENoZWNr\r\nIGl0IG91dDogaHR0cDovL3NpdGUvdXJsLz90b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuJnJzdnA9\r\nYWNjZXB0CgoKNiBJbnZpdGVkOgrCtyBlbWFpbDEgbmFtZSAoSG9zdCkgd2l0aCAyIHBlb3BsZQrCtyBl\r\nbWFpbDIgbmFtZQrCtyB0d2l0dGVyMyBuYW1lCsK3IGZhY2Vib29rNCBuYW1lCgoKRGVzY3JpcHRpb24K\r\nLS0tLS0tLQogIHRlc3QgY3Jvc3MgZGVzY3JpcHRpb24KCgojIFJlcGx5IHRoaXMgZW1haWwgZGlyZWN0\r\nbHkgYXMgY29udmVyc2F0aW9uLiAj\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9sYWJlbCB7\r\nCiAgICAgICAgICAgICAgICBiYWNrZ3JvdW5kLWNvbG9yOiAjRDVFOEYyOwogICAgICAgICAgICAgICAg\r\nY29sb3I6ICMzYTZlYTU7CiAgICAgICAgICAgICAgICBmb250LXNpemU6IDExcHg7CiAgICAgICAgICAg\r\nICAgICBwYWRkaW5nOiAwIDJweCAwIDJweDsKICAgICAgICAgICAgfQogICAgICAgICAgICAuZXhmZV9t\r\nYWlsX21hdGVzIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICAgICAg\r\nZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICB9CiAgICAgICAgICAgIC5leGZlX21haWxfaWRlbnRp\r\ndHkgewogICAgICAgICAgICAgICAgZm9udC1zdHlsZTogaXRhbGljOwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfaWRlbnRpdHlfbmFtZSB7CiAgICAgICAgICAgICAgICBjb2xvcjogIzE5\r\nMTkxOTsKICAgICAgICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4K\r\nICAgICAgICA8dGFibGUgd2lkdGg9IjY0MCIgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNw\r\nYWNpbmc9IjAiIHN0eWxlPSJmb250LWZhbWlseTogSGVsdmV0aWNhOyBmb250LXNpemU6IDEzcHg7IGxp\r\nbmUtaGVpZ2h0OiAxOXB4OyBjb2xvcjogIzE5MTkxOTsgZm9udC13ZWlnaHQ6IG5vcm1hbDsgcGFkZGlu\r\nZzogMzBweCA0MHB4IDMwcHggNDBweDsgYmFja2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsgbWluLWhlaWdo\r\ndDogNTYycHg7Ij4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGNvbHNwYW49IjMi\r\nIHZhbGlnbj0idG9wIiBzdHlsZT0iZm9udC1zaXplOiAzMnB4OyBsaW5lLWhlaWdodDogMzhweDsgcGFk\r\nZGluZy1ib3R0b206IDE4cHg7Ij4KICAgICAgICAgICAgICAgICAgICA8YSBocmVmPSJodHRwOi8vc2l0\r\nZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0iY29sb3I6ICMzYTZlYTU7\r\nIHRleHQtZGVjb3JhdGlvbjogbm9uZTsgZm9udC13ZWlnaHQ6IDMwMDsiPgogICAgICAgICAgICAgICAg\r\nICAgICAgICBUZXN0IENyb3NzCiAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAg\r\nPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRk\r\nIHdpZHRoPSIzNDAiIHN0eWxlPSJ2ZXJ0aWNhbC1hbGlnbjogYmFzZWxpbmU7IGZvbnQtd2VpZ2h0OiAz\r\nMDA7Ij4KICAgICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIg\r\nY2VsbHNwYWNpbmc9IjAiPgogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHN0eWxlPSJwYWRkaW5nLWJvdHRvbTogMjBweDsg\r\nZm9udC1zaXplOiAyMHB4OyB2ZXJ0aWNhbC1hbGlnbjogYmFzZWxpbmU7Ij4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBZb3UncmUgZ2F0\r\naGVyaW5nIHRoaXMgPHNwYW4gc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5v\r\nbmU7Ij7Ct1jCtzwvc3Bhbj4uCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAg\r\nICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQ+CiAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2VsbHBhZGRpbmc9\r\nIjAiIGNlbGxzcGFjaW5nPSIwIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRy\r\nPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIiB3\r\naWR0aD0iMTYwIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8YSBo\r\ncmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBzdHlsZT0i\r\ndGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgCQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9ImZv\r\nbnQtc2l6ZTogMjBweDsgbGluZS1oZWlnaHQ6IDI2cHg7IG1hcmdpbjogMDsgY29sb3I6ICMzMzMzMzM7\r\nIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDQ6NDVQ\r\nTSBvbiBUdWUsIE9jdCAyMyAyMDEyCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgIDwvcD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgIDx0ZCB2YWxpZ249InRvcCIgc3R5bGU9InBhZGRpbmctbGVmdDogMTBw\r\neDsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxhIGhyZWY9Imh0\r\ndHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iIHN0eWxlPSJ0ZXh0LWRl\r\nY29yYXRpb246IG5vbmU7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0\r\neWxlPSJmb250LXNpemU6IDIwcHg7IGxpbmUtaGVpZ2h0OiAyNnB4OyBtYXJnaW46IDA7IGNvbG9yOiAj\r\nMzMzMzMzOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICBQbGFjZQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3A+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJj\r\nb2xvcjogIzE5MTkxOTsgbWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICBUbyBiZSBkZWNpZGVkCiAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgIDwvcD4gCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nIDwvYT4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgPC90YWJsZT4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAg\r\nICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAg\r\nICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHN0\r\neWxlPSJwYWRkaW5nLXRvcDogMzBweDsgcGFkZGluZy1ib3R0b206IDMwcHg7Ij4KICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8YSBzdHlsZT0iZmxvYXQ6IGxlZnQ7IGRpc3BsYXk6IGJsb2NrOyB0\r\nZXh0LWRlY29yYXRpb246IG5vbmU7IGJvcmRlcjogMXB4IHNvbGlkICNiZWJlYmU7IGJhY2tncm91bmQt\r\nY29sb3I6ICMzQTZFQTU7IGNvbG9yOiAjRkZGRkZGOyBwYWRkaW5nOiA1cHggMzBweCA1cHggMzBweDsg\r\nbWFyZ2luLWxlZnQ6IDI1cHg7IiBhbHQ9IkFjY2VwdCIgaHJlZj0iaHR0cDovL3NpdGUvdXJsLz90b2tl\r\nbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuJnJzdnA9YWNjZXB0Ij5JJ20gaW48L2E+CiAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgPGEgc3R5bGU9ImZsb2F0OiBsZWZ0OyBkaXNwbGF5OiBibG9jazsg\r\ndGV4dC1kZWNvcmF0aW9uOiBub25lOyBib3JkZXI6IDFweCBzb2xpZCAjYmViZWJlOyBiYWNrZ3JvdW5k\r\nLWNvbG9yOiAjRTZFNkU2OyBjb2xvcjogIzE5MTkxOTsgcGFkZGluZzogNXB4IDI1cHggNXB4IDI1cHg7\r\nIG1hcmdpbi1sZWZ0OiAxNXB4OyIgYWx0PSJDaGVjayBpdCBvdXQiIGhyZWY9Imh0dHA6Ly9zaXRlL3Vy\r\nbC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iPkNoZWNrIGl0IG91dC4uLjwvYT4KICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAg\r\nICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiPgogICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgIHRlc3QgY3Jvc3MgZGVzY3JpcHRpb24KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAg\r\nICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgPC90YWJsZT4KICAgICAgICAgICAgICAgIDwvdGQ+\r\nCiAgICAgICAgICAgICAgICA8dGQgd2lkdGg9IjMwIj48L3RkPgogICAgICAgICAgICAgICAgPHRkIHZh\r\nbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGlu\r\nZz0iMCIgY2VsbHNwYWNpbmc9IjAiPgogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZCBoZWlnaHQ9\r\nIjY4IiB2YWxpZ249InRvcCIgYWxpZ249InJpZ2h0Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAg\r\nPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgPHRkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0YWJsZSBib3JkZXI9IjAiIHN0\r\neWxlPSJjb2xvcjogIzMzMzMzMzsiIGNlbGxwYWRkaW5nPSIwIiBjZWxsc3BhY2luZz0iMCI+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQgd2lk\r\ndGg9IjI1IiBoZWlnaHQ9IjI1IiBhbGlnbj0ibGVmdCIgdmFsaWduPSJ0b3AiPgogICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxpbWcgd2lkdGg9IjIwIiBoZWlnaHQ9IjIwIiB0\r\naXRsZT0iZW1haWwxIG5hbWUiIGFsdD0iZW1haWwxIG5hbWUiIHNyYz0iaHR0cDovL3BhdGgvdG8vZW1h\r\naWwxLmF2YXRhciI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkPgogICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICA8c3Bhbj5lbWFpbDEgbmFtZTwvc3Bhbj4gPHNwYW4gY2xh\r\nc3M9ImV4ZmVfbWFpbF9tYXRlcyI+KzI8L3NwYW4+IDxzcGFuIGNsYXNzPSJleGZlX21haWxfbGFiZWwi\r\nPmhvc3Q8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIyNSIgaGVpZ2h0\r\nPSIyNSIgYWxpZ249ImxlZnQiIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8aW1nIHdpZHRoPSIyMCIgaGVpZ2h0PSIyMCIgdGl0bGU9ImVtYWlsMiBu\r\nYW1lIiBhbHQ9ImVtYWlsMiBuYW1lIiBzcmM9Imh0dHA6Ly9wYXRoL3RvL2VtYWlsMi5hdmF0YXIiPgog\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgPHNwYW4+ZW1haWwyIG5hbWU8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8\r\nL3RyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgPHRkIHdpZHRoPSIyNSIgaGVpZ2h0PSIyNSIgYWxpZ249ImxlZnQiIHZhbGlnbj0idG9wIj4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8aW1nIHdpZHRoPSIyMCIgaGVp\r\nZ2h0PSIyMCIgdGl0bGU9InR3aXR0ZXIzIG5hbWUiIGFsdD0idHdpdHRlcjMgbmFtZSIgc3JjPSJodHRw\r\nOi8vcGF0aC90by90d2l0dGVyMy5hdmF0YXIiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4K\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHNwYW4+dHdpdHRlcjMgbmFt\r\nZTwvc3Bhbj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQgd2lkdGg9IjI1IiBoZWlnaHQ9IjI1\r\nIiBhbGlnbj0ibGVmdCIgdmFsaWduPSJ0b3AiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgIDxpbWcgd2lkdGg9IjIwIiBoZWlnaHQ9IjIwIiB0aXRsZT0iZmFjZWJvb2s0IG5h\r\nbWUiIGFsdD0iZmFjZWJvb2s0IG5hbWUiIHNyYz0iaHR0cDovL3BhdGgvdG8vZmFjZWJvb2s0LmF2YXRh\r\nciI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkPgogICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICA8c3Bhbj5mYWNlYm9vazQgbmFtZTwvc3Bhbj4KICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgIDwvdGFibGU+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8\r\nL3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgIDwvdGFi\r\nbGU+CiAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICA8L3RyPgogICAgICAgICAgICA8dHI+\r\nCiAgICAgICAgICAgICAgICA8dGQgY29sc3Bhbj0iMyIgc3R5bGU9ImZvbnQtc2l6ZTogMTFweDsgbGlu\r\nZS1oZWlnaHQ6IDE1cHg7IGNvbG9yOiAjN0Y3RjdGOyBwYWRkaW5nLXRvcDogNDBweDsiPgogICAgICAg\r\nICAgICAgICAgICAgIFJlcGx5IHRoaXMgZW1haWwgZGlyZWN0bHkgYXMgY29udmVyc2F0aW9uLCBvciBU\r\ncnkgPGEgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJo\r\ndHRwOi8vYXBwL3VybCI+RVhGRTwvYT4gYXBwLgogICAgICAgICAgICAgICAgICAgIDxiciAvPgogICAg\r\nICAgICAgICAgICAgICAgIFRoaXMgPGEgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRp\r\nb246IG5vbmU7IiBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rv\r\na2VuIj7Ct1jCtzwvYT4gaW52aXRhdGlvbiBpcyBzZW50IGJ5IDxzcGFuIGNsYXNzPSJleGZlX21haWxf\r\naWRlbnRpdHlfbmFtZSI+ZW1haWwxIG5hbWU8L3NwYW4+IGZyb20gPGEgc3R5bGU9ImNvbG9yOiAjM2E2\r\nZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8vc2l0ZS91cmwiPkVYRkU8L2E+\r\nLgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICA8L3RhYmxlPgog\r\nICAgPC9ib2R5Pgo8L2h0bWw+Cg==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6VGVzdCBDcm9zcw0KWC1XUi1DQUxERVNDOmV4ZmUgY2FsDQpYLVdSLVRJTUVaT05F\r\nOiswODAwDQpCRUdJTjpWRVZFTlQNClVJRDohMTIzQGV4ZmUNCkRUU1RBTVA6MjAxMjEwMjNUMDg0NTAw\r\nWg0KREVTQ1JJUFRJT046dGVzdCBjcm9zcyBkZXNjcmlwdGlvbg0KRFRTVEFSVDoyMDEyMTAyM1QwODQ1\r\nMDBaDQpMT0NBVElPTjoNClNVTU1BUlk6VGVzdCBDcm9zcw0KVVJMOmh0dHA6Ly9zaXRlL3VybC8jIXRv\r\na2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4NCkVORDpWRVZFTlQNCkVORDpWQ0FMRU5EQVINCg==\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
	assert.Equal(t, private, expectPrivate)
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
	expectPrivate := "\n\n\n\nSuccessfully gathering \\(\"Test Cross\"\\), \\(4:45PM on Tue, Oct 23 2012\\). 6 invited: email1 name, email2 name, twitter3 name... http://site/url/#!123/ecip"
	assert.Equal(t, private, expectPrivate)
	assert.Equal(t, public, "Invitation: http://site/url/#!123/ecip (Please follow @EXFE to receive details PRIVATELY through Direct Message.)")
}
