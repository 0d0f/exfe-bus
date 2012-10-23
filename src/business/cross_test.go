package business

import (
	"formatter"
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

func TestCrossSummary(t *testing.T) {
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

	updates := []model.Update{
		model.Update{
			To:       remail1,
			OldCross: cross2,
			Cross:    cross1,
			By:       facebook4,
		},
	}

	cross2 = cross1
	cross1.Time = time1
	updates = append(updates, model.Update{
		To:       remail1,
		OldCross: cross2,
		Cross:    cross1,
		By:       email1,
	})

	cross2 = cross1
	cross1.Place = place1
	updates = append(updates, model.Update{
		To:       remail1,
		OldCross: cross2,
		Cross:    cross1,
		By:       email2,
	})

	cross2 = cross1
	cross1.Title = "New Title"
	updates = append(updates, model.Update{
		To:       remail1,
		OldCross: cross2,
		Cross:    cross1,
		By:       twitter3,
	})

	t.Errorf("%+v\n", updates)

	c := NewCross(localTemplate, &config)
	private, public, err := c.getContent(updates)
	t.Logf("err: %s", err)
	t.Errorf("private:-----start------\n%s\n-------end-------", private)
	t.Errorf("public:-----start------\n%s\n-------end-------", public)
}
