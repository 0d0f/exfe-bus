package notifier

import (
	"bytes"
	"fmt"
	"formatter"
	"gobus"
	"model"
	"service/args"
	"thirdpart"
)

type SummaryArg struct {
	ArgBase
	OldCross *model.Cross     `json:"-"`
	Cross    model.Cross      `json:"-"`
	Bys      []model.Identity `json:"-"`

	NewInvited    []model.Invitation `json:"-"`
	Removed       []model.Invitation `json:"-"`
	NewAccepted   []model.Invitation `json:"-"`
	OldAccepted   []model.Invitation `json:"-"`
	NewDeclined   []model.Invitation `json:"-"`
	NewInterested []model.Invitation `json:"-"`
	NewPending    []model.Invitation `json:"-"`
}

func SummaryFromUpdates(updates []model.CrossUpdate, config *model.Config) (*SummaryArg, error) {
	if updates == nil && len(updates) == 0 {
		return nil, fmt.Errorf("no update info")
	}

	to := updates[0].To
	bys := make([]model.Identity, 0)

Bys:
	for _, update := range updates {
		if !to.Equal(&update.To) {
			return nil, fmt.Errorf("updates not send to same recipient: %s, %s", to, update.To)
		}
		for _, i := range bys {
			if update.By.SameUser(i) {
				continue Bys
			}
		}
		bys = append(bys, update.By)
	}

	ret := &SummaryArg{
		Bys:      bys,
		OldCross: &updates[0].OldCross,
		Cross:    updates[len(updates)-1].Cross,

		NewInvited:    make([]model.Invitation, 0),
		Removed:       make([]model.Invitation, 0),
		NewAccepted:   make([]model.Invitation, 0),
		OldAccepted:   make([]model.Invitation, 0),
		NewDeclined:   make([]model.Invitation, 0),
		NewInterested: make([]model.Invitation, 0),
		NewPending:    make([]model.Invitation, 0),
	}
	ret.To = to
	err := ret.ArgBase.Parse(config)
	if err != nil {
		return nil, err
	}

	ret.Cross.Exfee.Parse()
	ret.OldCross.Exfee.Parse()

	for _, i := range ret.Cross.Exfee.Accepted {
		if !in(&i, ret.OldCross.Exfee.Accepted) {
			ret.NewAccepted = append(ret.NewAccepted, i)
		} else {
			ret.OldAccepted = append(ret.OldAccepted, i)
		}
	}
	for _, i := range ret.Cross.Exfee.Declined {
		if !in(&i, ret.OldCross.Exfee.Declined) {
			ret.NewDeclined = append(ret.NewDeclined, i)
		}
	}
	for _, i := range ret.Cross.Exfee.Interested {
		if !in(&i, ret.OldCross.Exfee.Interested) {
			ret.NewInterested = append(ret.NewInterested, i)
		}
	}
	for _, i := range ret.Cross.Exfee.Pending {
		if !in(&i, ret.OldCross.Exfee.Pending) {
			ret.NewPending = append(ret.NewPending, i)
		}
	}
	for _, i := range ret.Cross.Exfee.Invitations {
		if !in(&i, ret.OldCross.Exfee.Invitations) {
			ret.NewInvited = append(ret.NewInvited, i)
		}
	}
	for _, i := range ret.OldCross.Exfee.Invitations {
		if !in(&i, ret.Cross.Exfee.Invitations) {
			ret.Removed = append(ret.Removed, i)
		}
	}
	return ret, nil
}

func (a *SummaryArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
}

func (a *SummaryArg) TotalOldAccepted() int {
	ret := 0
	for _, e := range a.OldAccepted {
		ret += 1 + int(e.Mates)
	}
	return ret
}

func (a *SummaryArg) IsTimeChanged() bool {
	oldtime, _ := a.OldCross.Time.StringInZone(a.To.Timezone)
	time, _ := a.Cross.Time.StringInZone(a.To.Timezone)
	return oldtime != time
}

func (a *SummaryArg) IsTitleChanged() bool {
	return a.OldCross.Title != a.Cross.Title
}

func (a *SummaryArg) IsPlaceChanged() bool {
	return !a.Cross.Place.Same(&a.OldCross.Place)
}

func (a *SummaryArg) IsPlaceTitleChanged() bool {
	return a.Cross.Place.Title != a.OldCross.Place.Title
}

func (a *SummaryArg) IsPlaceDescChanged() bool {
	return a.Cross.Place.Description != a.OldCross.Place.Description
}

func (a *SummaryArg) IsDescriptionChanged() bool {
	return a.Cross.Description != a.OldCross.Description
}

func (a *SummaryArg) IsComboChanged() bool {
	changedNumber := 0
	if a.IsTimeChanged() {
		changedNumber++
	}
	if a.IsTimeChanged() {
		changedNumber++
	}
	if a.IsPlaceTitleChanged() {
		changedNumber++
	}
	if a.IsPlaceDescChanged() {
		changedNumber++
	}
	if a.IsDescriptionChanged() {
		changedNumber++
	}

	return changedNumber > 1
}

func (a *SummaryArg) ListBy(limit int, join string) string {
	buf := bytes.NewBuffer(nil)
	for i, by := range a.Bys {
		if buf.Len() > 0 {
			buf.WriteString(join)
		}
		if i >= limit {
			buf.WriteString("etc")
			break
		}
		buf.WriteString(by.Name)
	}
	return buf.String()
}

func (a *SummaryArg) NeedShowBy() bool {
	if len(a.Bys) != 1 {
		return false
	}
	if a.To.SameUser(&a.Bys[0]) {
		return false
	}
	return true
}

type InvitationArg struct {
	ArgBase
	Cross model.Cross `json:"cross"`
}

func (a InvitationArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
}

func (a InvitationArg) IsCreator() bool {
	return a.To.SameUser(&a.Cross.By)
}

func (a InvitationArg) LongDescription() bool {
	if len(a.Cross.Description) > 200 {
		return true
	}
	return false
}

func (a InvitationArg) ListInvitations() string {
	l := len(a.Cross.Exfee.Invitations)
	max := 3
	ret := ""
	for i := 0; i < 3 && i < l; i++ {
		if i > 0 {
			ret += ", "
		}
		ret += a.Cross.Exfee.Invitations[i].Identity.Name
	}
	if l > max {
		ret += "…"
	}
	return ret
}

type Cross struct {
	localTemplate *formatter.LocalTemplate
	config        *model.Config
}

func NewCross(localTemplate *formatter.LocalTemplate, config *model.Config) *Cross {
	return &Cross{
		localTemplate: localTemplate,
		config:        config,
	}
}

func (c *Cross) Summary(updates []model.CrossUpdate) error {
	private, public, err := c.getSummaryContent(updates)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}

	url := fmt.Sprintf("http://%s:%d", c.config.ExfeService.Addr, c.config.ExfeService.Port)
	client, err := gobus.NewClient(fmt.Sprintf("%s/%s", url, "Thirdpart"))
	if err != nil {
		return fmt.Errorf("can't create gobus client: %s", err)
	}

	arg := args.SendArg{
		To:             &updates[0].To,
		PrivateMessage: private,
		PublicMessage:  public,
		Info: &thirdpart.InfoData{
			CrossID: updates[0].Cross.ID,
			Type:    thirdpart.CrossUpdate,
		},
	}
	var ids string
	err = client.Do("Send", &arg, &ids)
	if err != nil {
		return fmt.Errorf("send error: %s", err)
	}
	return nil
}

func (c *Cross) Invite(invitations []InvitationArg) error {
	url := fmt.Sprintf("http://%s:%d/Thirdpart", c.config.ExfeService.Addr, c.config.ExfeService.Port)
	client, err := gobus.NewClient(url)
	if err != nil {
		return fmt.Errorf("can't create gobus client: %s", err)
	}

	for _, i := range invitations {
		private, public, err := c.getInvitationContent(i)
		if err != nil {
			return fmt.Errorf("can't get content: %s", err)
		}

		arg := args.SendArg{
			To:             &i.To,
			PrivateMessage: private,
			PublicMessage:  public,
			Info: &thirdpart.InfoData{
				CrossID: i.Cross.ID,
				Type:    thirdpart.CrossInvitation,
			},
		}
		var ids string
		err = client.Do("Send", &arg, &ids)
		if err != nil {
			return fmt.Errorf("send error: %s", err)
		}
	}
	return nil
}

func (c *Cross) getSummaryContent(updates []model.CrossUpdate) (string, string, error) {
	arg, err := SummaryFromUpdates(updates, c.config)
	if err != nil {
		return "", "", err
	}

	private, err := GetContent(c.localTemplate, "cross_summary", arg)
	if err != nil {
		return "", "", fmt.Errorf("can't get content: %s", err)
	}
	public, err := GetContent(c.localTemplate, "cross_summary_public", arg)
	if err != nil {
		return "", "", fmt.Errorf("can't get content: %s", err)
	}

	return private, public, nil
}

func (c *Cross) getInvitationContent(arg InvitationArg) (string, string, error) {
	err := arg.Parse(c.config)
	if err != nil {
		return "", "", err
	}
	arg.Cross.Exfee.Parse()

	private, err := GetContent(c.localTemplate, "cross_invitation", arg)
	if err != nil {
		return "", "", fmt.Errorf("can't get content: %s", err)
	}
	public, err := GetContent(c.localTemplate, "cross_invitation_public", arg)
	if err != nil {
		return "", "", fmt.Errorf("can't get content: %s", err)
	}

	return private, public, nil
}

func in(id *model.Invitation, ids []model.Invitation) bool {
	for _, i := range ids {
		if id.Identity.SameUser(i.Identity) {
			return true
		}
	}
	return false
}
