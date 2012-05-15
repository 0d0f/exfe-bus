package exfe_service

import (
	"exfe/model"
	"testing"
	"log/syslog"
	"encoding/json"
	"bytes"
)

var exfee = exfe_model.Exfee{
	Invitations: []exfe_model.Invitation{
		exfe_model.Invitation{
			Rsvp_status: "ACCEPTED",
			Identity: exfe_model.Identity{
				Id: 1,
				Name: "Tester1",
				Connected_user_id: 1,
			},
		},
		exfe_model.Invitation{
			Rsvp_status: "ACCEPTED",
			Identity: exfe_model.Identity{
				Id: 2,
				Name: "Tester2",
				Connected_user_id: 2,
			},
		},
		exfe_model.Invitation{
			Rsvp_status: "ACCEPTED",
			Identity: exfe_model.Identity{
				Id: 3,
				Name: "Tester3",
				Connected_user_id: 3,
			},
		},
		exfe_model.Invitation{
			Rsvp_status: "ACCEPTED",
			Identity: exfe_model.Identity{
				Id: 4,
				Name: "Tester4",
				Connected_user_id: 4,
			},
		},
	},
}

func TestExfeeStatus(t *testing.T) {
	log, _ := syslog.New(syslog.LOG_INFO, "test")
	var new_, old exfe_model.Exfee

	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	e.Encode(exfee)

	buf1 := bytes.NewBufferString(buf.String())
	d1 := json.NewDecoder(buf1)
	d1.Decode(&new_)
	buf2 := bytes.NewBufferString(buf.String())
	d2 := json.NewDecoder(buf2)
	d2.Decode(&old)

	new_.Invitations[0].Identity.Connected_user_id = 5
	old.Invitations[1].Identity.Connected_user_id = 6
	old.Invitations[2].Rsvp_status = "DECLINED"
	new_.Invitations[3].Rsvp_status = "DECLINED"

	accepted, declined, newlyInvited, removed := newStatusUser(log, &old, &new_)
	t.Logf("accepted: %v", accepted)
	t.Logf("declined: %v", declined)
	t.Logf("newly invited: %v", newlyInvited)
	t.Logf("removed: %v", removed)
	t.Errorf("print")
}
