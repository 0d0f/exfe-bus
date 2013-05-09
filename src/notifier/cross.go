package notifier

import (
	"broker"
	"bytes"
	"fmt"
	"formatter"
	"model"
	"net/url"
)

type Cross struct {
	localTemplate *formatter.LocalTemplate
	config        *model.Config
	platform      *broker.Platform
}

func NewCross(localTemplate *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) *Cross {
	return &Cross{
		localTemplate: localTemplate,
		config:        config,
		platform:      platform,
	}
}

func (c Cross) V3Digest(requests []model.CrossDigestRequest) error {
	if len(requests) == 0 {
		return fmt.Errorf("len(requests) == 0")
	}
	to := requests[len(requests)-1].To
	crossId := requests[0].CrossId
	updatedAt := requests[0].UpdatedAt

	query := make(url.Values)
	query.Set("updated_at", updatedAt)
	query.Set("user_id", fmt.Sprint("%d", to.UserID))
	cross, err := c.platform.FindCross(crossId, query)
	if err != nil {
		return err
	}

	arg := map[string]interface{}{
		"To":     to,
		"Cross":  cross,
		"Config": c.config,
	}
	text, err := GenerateContent(c.localTemplate, "v3_cross_digest", to.Provider, to.Language, arg)
	if err != nil {
		return err
	}
	_, err = c.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

func (c Cross) V3Invitation(invitation model.CrossInvitation) error {
	invitation.Config = c.config
	to := invitation.To

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", to.UserID))
	cross, err := c.platform.FindCross(invitation.CrossId, query)
	if err != nil {
		return err
	}
	invitation.Cross = cross

	text, err := GenerateContent(c.localTemplate, "v3_cross_invitation", to.Provider, to.Language, invitation)
	if err != nil {
		return err
	}
	_, err = c.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

func (c Cross) V3Summary(updates []model.CrossUpdate) error {
	if len(updates) == 0 {
		return fmt.Errorf("len(updates) == 0")
	}

	to := updates[0].To
	selfUpdates := true
	for _, update := range updates {
		if !to.SameUser(&update.By) {
			selfUpdates = false
		}
	}
	if selfUpdates {
		c.config.Log.Debug("not send with all self updates: %s", to)
		return nil
	}

	arg, err := SummaryFromUpdates(updates, c.config)
	if err != nil {
		return err
	}

	if !arg.IsChanged() {
		return nil
	}

	to = arg.To
	text, err := GenerateContent(c.localTemplate, "v3_cross_summary", to.Provider, to.Language, arg)
	if err != nil {
		return err
	}
	_, err = c.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

func in(id *model.Invitation, ids []model.Invitation) bool {
	for _, i := range ids {
		if id.Identity.SameUser(i.Identity) {
			return true
		}
	}
	return false
}

func (c *Cross) Summary(updates model.CrossUpdates) error {
	if len(updates) == 0 {
		return fmt.Errorf("len(updates) == 0")
	}

	to := updates[0].To
	if to.Provider == "twitter" {
		c.config.Log.Debug("not send to twitter: %s", to)
		return nil
	}
	selfUpdates := true
	for _, update := range updates {
		if !to.SameUser(&update.By) {
			selfUpdates = false
		}
	}
	if selfUpdates {
		c.config.Log.Debug("not send with all self updates: %s", to)
		return nil
	}

	text, err := c.getSummaryContent(updates)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}

	if text == "" {
		return nil
	}

	_, err = c.platform.Send(to, text)
	if err != nil {
		return fmt.Errorf("send error: %s", err)
	}
	return nil
}

func (c *Cross) Invite(invitation model.CrossInvitation) error {
	text, err := c.getInvitationContent(invitation)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}

	_, err = c.platform.Send(invitation.To, text)
	if err != nil {
		return fmt.Errorf("send error: %s", err)
	}

	return nil
}

func (c *Cross) getSummaryContent(updates []model.CrossUpdate) (string, error) {
	arg, err := SummaryFromUpdates(updates, c.config)
	if err != nil {
		return "", err
	}

	if !arg.IsChanged() {
		return "", nil
	}

	text, err := GetContent(c.localTemplate, "cross_summary", arg.To, arg)
	if err != nil {
		return "", fmt.Errorf("can't get content: %s", err)
	}

	return text, nil
}

func (c *Cross) getInvitationContent(arg model.CrossInvitation) (string, error) {
	err := arg.Parse(c.config)
	if err != nil {
		return "", err
	}
	arg.Cross.Exfee.Parse()

	text, err := GetContent(c.localTemplate, "cross_invitation", arg.To, arg)
	if err != nil {
		return "", fmt.Errorf("can't get content: %s", err)
	}

	return text, nil
}

type SummaryArg struct {
	model.ThirdpartTo
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
	err := ret.Parse(config)
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
	// for _, i := range ret.Cross.Exfee.Interested {
	// 	if !in(&i, ret.OldCross.Exfee.Interested) {
	// 		ret.NewInterested = append(ret.NewInterested, i)
	// 	}
	// }
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

func (a *SummaryArg) NeedEmail() bool {
	if a.IsTitleChanged() {
		return true
	}
	if a.IsTimeChanged() {
		return true
	}
	if a.IsPlaceChanged() {
		return true
	}
	if a.IsDescriptionChanged() {
		return true
	}
	peopleChanged := len(a.NewInvited)
	peopleChanged += len(a.Removed)
	peopleChanged += len(a.NewAccepted)
	peopleChanged += len(a.NewDeclined)
	if peopleChanged > 0 {
		return true
	}
	return false
}

func (a *SummaryArg) IsChanged() bool {
	if a.IsTitleChanged() {
		return true
	}
	if a.IsTimeChanged() {
		return true
	}
	if a.IsPlaceChanged() {
		return true
	}
	if a.IsDescriptionChanged() {
		return true
	}
	if a.IsExfeeChanged() {
		return true
	}
	peopleChanged := len(a.NewInvited)
	peopleChanged += len(a.Removed)
	if peopleChanged > 0 {
		return true
	}
	return false
}

func (a *SummaryArg) IsExfeeChanged() bool {
	peopleChanged := len(a.NewAccepted)
	peopleChanged += len(a.NewDeclined)
	peopleChanged += len(a.NewInterested)
	peopleChanged += len(a.NewPending)
	if peopleChanged > 0 {
		return true
	}
	return false
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
	return !a.Cross.Place.Same(a.OldCross.Place)
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

func (a SummaryArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a SummaryArg) PublicLink() string {
	return fmt.Sprintf("%s/#!%d/%s", a.Config.SiteUrl, a.Cross.ID, a.To.Token[1:5])
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

func (a *SummaryArg) ShowBy(ids []model.Invitation) bool {
	if len(a.Bys) != 1 {
		return false
	}
	if len(ids) != 1 {
		return false
	}
	if a.Bys[0].SameUser(ids[0].Identity) {
		return false
	}
	return true
}
