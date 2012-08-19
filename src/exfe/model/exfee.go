package exfe_model

type Exfee struct {
	Id          uint64
	Invitations []Invitation
}

func (e *Exfee) FindInvitation(identity *Identity) *Invitation {
	for _, inv := range e.Invitations {
		if identity.DiffId() == inv.Identity.DiffId() {
			return &inv
		}
	}
	return nil
}
