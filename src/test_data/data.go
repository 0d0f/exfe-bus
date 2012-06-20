package exfe_data

import (
	"exfe/model"
)

var Leo_email = exfe_model.Identity{
	Id: 1,
	Name: "Leonard",
	Nickname: "Leo.",
	Bio: "I am a physicist at CalTech and live with my best friend Sheldon.",
	Provider: "email",
	Timezone: "-07:00 PST",
	Connected_user_id: 1,
	External_id: "tester_leonard@0d0f.com",
	External_username: "tester_leonard@0d0f.com",
	Avatar_filename: "https://twimg0-a.akamaihd.net/profile_images/1204136991/johnny-galecki-as-leonard-hofstadter.jpg",
}

var Leo_twitter = exfe_model.Identity{
	Id: 2,
	Name: "Leonard",
	Nickname: "Leo.",
	Bio: "I am a physicist at CalTech and live with my best friend Sheldon.",
	Provider: "twitter",
	Timezone: "-06:00",
	Connected_user_id: 1,
	External_id: "575129929",
	External_username: "0d0f_tester_leo",
	Avatar_filename: "https://twimg0-a.akamaihd.net/profile_images/1204136991/johnny-galecki-as-leonard-hofstadter.jpg",
}

var Raj_email = exfe_model.Identity{
	Id: 3,
	Name: "Rajesh",
	Nickname: "Raj.",
	Bio: "Give me a grasshopper and I'm ready to go!",
	Provider: "email",
	Timezone: "-07:00 PST",
	Connected_user_id: 2,
	External_id: "tester_raj@0d0f.com",
	External_username: "tester_raj@0d0f.com",
	Avatar_filename: "https://twimg0-a.akamaihd.net/profile_images/198817661/200px-Kunal-Nayyar.jpg",
}

var Raj_twitter = exfe_model.Identity{
	Id: 4,
	Name: "Rajesh",
	Nickname: "Raj.",
	Bio: "Give me a grasshopper and I'm ready to go!",
	Provider: "twitter",
	Timezone: "-06:00",
	Connected_user_id: 2,
	External_id: "575215638",
	External_username: "0d0f_tester_raj",
	Avatar_filename: "https://twimg0-a.akamaihd.net/profile_images/198817661/200px-Kunal-Nayyar.jpg",
}

var How_email = exfe_model.Identity{
	Id: 5,
	Name: "Howard",
	Nickname: "How.",
	Bio: "My mom lives with me! Got that?",
	Provider: "email",
	Timezone: "-07:00 PST",
	Connected_user_id: 3,
	External_id: "tester_how@rd0d0f.com",
	External_username: "tester_howard@0d0f.com",
	Avatar_filename: "https://twimg0-a.akamaihd.net/profile_images/198833248/simon-helberg.jpg",
}

var How_twitter = exfe_model.Identity{
	Id: 6,
	Name: "Howard",
	Nickname: "How.",
	Bio: "My mom lives with me! Got that?",
	Provider: "twitter",
	Timezone: "-06:00",
	Connected_user_id: 3,
	External_id: "575216679",
	External_username: "0d0f_tester_how",
	Avatar_filename: "https://twimg0-a.akamaihd.net/profile_images/198833248/simon-helberg.jpg",
}

var She_email = exfe_model.Identity{
	Id: 7,
	Name: "Sheldon",
	Nickname: "She.",
	Bio: "Quite possibly the most intelligent human on the planet. Brilliant theoretical physicist.",
	Provider: "email",
	Timezone: "-07:00 PST",
	Connected_user_id: 4,
	External_id: "tester_sheldon@0d0f.com",
	External_username: "tester_sheldon@0d0f.com",
	Avatar_filename: "https://twimg0-a.akamaihd.net/profile_images/365042597/sheldon.jpg",
}

var She_twitter = exfe_model.Identity{
	Id: 8,
	Name: "Sheldon",
	Nickname: "She.",
	Bio: "Quite possibly the most intelligent human on the planet. Brilliant theoretical physicist.",
	Provider: "twitter",
	Timezone: "-06:00",
	Connected_user_id: 4,
	External_id: "575131718",
	External_username: "0d0f_tester_she",
	Avatar_filename: "https://twimg0-a.akamaihd.net/profile_images/365042597/sheldon.jpg",
}

var Exfee1 = exfe_model.Exfee{
	Id: 1,
	Invitations: []exfe_model.Invitation{
		exfe_model.Invitation{
			Id: 1,
			Token: "sheemailtoken",
			Host: true,
			Mates: 0,
			Identity: She_email,
			Rsvp_status: "ACCEPTED",
			By_identity: She_email,
		},
		exfe_model.Invitation{
			Id: 2,
			Token: "shetwittertoken",
			Host: false,
			Mates: 0,
			Identity: She_twitter,
			Rsvp_status: "NOTIFICATION",
			By_identity: She_email,
		},
		exfe_model.Invitation{
			Id: 3,
			Token: "leoemailtoken",
			Host: false,
			Mates: 1,
			Identity: Leo_email,
			Rsvp_status: "NORESPONSE",
			By_identity: She_email,
		},
		exfe_model.Invitation{
			Id: 4,
			Token: "rajemailtoken",
			Host: false,
			Mates: 0,
			Identity: Raj_email,
			Rsvp_status: "NORESPONSE",
			By_identity: She_email,
		},
	},
}

var Exfee2 = exfe_model.Exfee{
	Id: 2,
	Invitations: []exfe_model.Invitation{
		exfe_model.Invitation{
			Id: 1,
			Token: "sheemailtoken",
			Host: true,
			Mates: 0,
			Identity: She_email,
			Rsvp_status: "ACCEPTED",
			By_identity: She_email,
		},
		exfe_model.Invitation{
			Id: 2,
			Token: "shetwittertoken",
			Host: false,
			Mates: 0,
			Identity: She_twitter,
			Rsvp_status: "NOTIFICATION",
			By_identity: She_email,
		},
		exfe_model.Invitation{
			Id: 3,
			Token: "leoemailtoken",
			Host: false,
			Mates: 1,
			Identity: Leo_email,
			Rsvp_status: "ACCEPTED",
			By_identity: She_email,
		},
		exfe_model.Invitation{
			Id: 5,
			Token: "howemailtoken",
			Host: false,
			Mates: 0,
			Identity: How_email,
			Rsvp_status: "DECLINED",
			By_identity: She_email,
		},
	},
}

var Cross = exfe_model.Cross{
	Id: 123,
	Title: "Team dinner in Sav Francisco with Bay Area friends",
	Description: `Lorem ipsum dolor sit amet, consectetur adipiscing elit, set eiusmod tempor incidunt et labore et dolore magna aliquam. Ut enim ad minim veniam, quis nostrud exerc. Irure dolor in reprehend incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

Duis aute irure dolor in reprehenderit in voluptate velit esse molestaie cillum. Tia non ob ea soluad incom dereud facilis est er expedit distinct. Nam liber te conscient to factor tum poen legum odioque civiuda et tam. Neque pecun modut est neque nonor et imper ned libidig met, consectetur adipiscing elit, sed ut labore et dolore magna aliquam is nostrud exercitation ullam mmodo consequet.`,
	Time: exfe_model.CrossTime{
		Begin_at: exfe_model.EFTime{
			Date_word: "Tomorrow",
			Date: "2012-06-18",
			Time_word: "Dinner",
			Time: "12:23:00",
			Timezone: "-08:00 PST",
		},
		Origin: "Tomorrow Dinner",
		OutputFormat: exfe_model.Format,
	},
	Place: exfe_model.Place{
		Id: 123,
		Title: "Crab House Pier 39, 2nd floor",
		Description: "Pier 39, 203 C, San Francisco, http://crabhouse39.com, (555) 434-2722",
	},
	By_identity: She_email,
	Exfee: Exfee2,
}

var Post1 = exfe_model.Post{
	Id: 1,
	By_identity: Leo_email,
	Content: "err... can't make it this fri!!!",
	Created_at: "2012-04-05 23:47:00",
}

var Post2 = exfe_model.Post{
	Id: 2,
	By_identity: How_email,
	Content: "My only missing food in US, dudes! yummy! Lorem ipsum dolor sit amet, ligula suspendisse nulla pretium, rhoncus tempor placerat fermentum, enim integer ad vestibulum volutpat. Nisl rhoncus turpis est, vel elit, congue wisi enim nunc ultricies …",
	Created_at: "2012-04-06 00:01:00",
}
