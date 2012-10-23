package model

import (
	"github.com/stretchrcom/testify/assert"
	"testing"
)

var email1 = Identity{
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

var email2 = Identity{
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

var twitter1 = Identity{
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

var twitter3 = Identity{
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

var facebook1 = Identity{
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

var facebook4 = Identity{
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

var exfee = Exfee{
	ID: 123,
	Invitations: []Invitation{
		Invitation{
			ID:         11,
			Host:       true,
			Mates:      2,
			Identity:   email1,
			RsvpStatus: RsvpAccepted,
			By:         email1,
		},
		Invitation{
			ID:         22,
			Identity:   email2,
			RsvpStatus: RsvpNoresponse,
			By:         email1,
		},
		Invitation{
			ID:         33,
			Identity:   twitter3,
			RsvpStatus: RsvpDeclined,
			By:         email1,
		},
		Invitation{
			ID:         44,
			Identity:   facebook4,
			RsvpStatus: RsvpInterested,
			By:         twitter3,
		},
	},
}

func TestExfeeParse(t *testing.T) {
	e := exfee
	e.Parse()
	assert.Equal(t, e.Accepted[0].ID, uint64(11))
	assert.Equal(t, e.Declined[0].ID, uint64(33))
	assert.Equal(t, e.Pending[0].ID, uint64(22))
	assert.Equal(t, e.Interested[0].ID, uint64(44))
}
