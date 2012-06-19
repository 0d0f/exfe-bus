package main

import (
	"test_data"
	"exfe/service"
	"exfe/model"
	"fmt"
	"os"
	"log"
)

func main() {
	c := exfe_service.InitConfig()

	old_cross := exfe_data.Cross
	old_cross.Title = "Team dinner"
	old_cross.Exfee = exfe_data.Exfee1

	arg := exfe_service.ProviderArg{
		Cross: &exfe_data.Cross,
		Old_cross: &old_cross,
		To_identity: &exfe_data.Leo_email,
		By_identities: []*exfe_model.Identity{&exfe_data.She_email, &exfe_data.Leo_email},
		Config: c,
		Posts: []*exfe_model.Post{&exfe_data.Post1, &exfe_data.Post2},
	}
	l := log.New(os.Stderr, "test", log.LstdFlags)
	arg.Diff(l)

	ce := exfe_service.NewCrossEmail(c)

	html, ics, err := ce.GetBody(&arg, "cross_invitation.html")
	fmt.Println("create cross invitation error:", err)
	f, _ := os.Create("cross_invitation.html")
	f.WriteString(html)
	f, _ = os.Create("cross.ics")
	f.WriteString(ics)

	html, ics, err = ce.GetBody(&arg, "cross_update.html")
	fmt.Println("create cross update error:", err)
	f, _ = os.Create("cross_update.html")
	f.WriteString(html)
}
