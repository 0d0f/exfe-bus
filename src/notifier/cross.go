package notifier

import (
	"broker"
	"fmt"
	"formatter"
	"logger"
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
	query.Set("user_id", fmt.Sprintf("%d", to.UserID))
	cross, err := c.platform.FindCross(crossId, query)
	if err != nil {
		return err
	}

	arg := map[string]interface{}{
		"To":     to,
		"Cross":  cross,
		"Config": c.config,
	}
	text, err := GenerateContent(c.localTemplate, "cross_digest", to.Provider, to.Language, arg)
	if err != nil {
		return err
	}
	_, err = c.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

func (c Cross) V3Remind(requests []model.CrossDigestRequest) error {
	if len(requests) == 0 {
		return fmt.Errorf("len(requests) == 0")
	}
	to := requests[len(requests)-1].To
	crossId := requests[0].CrossId

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", to.UserID))
	cross, err := c.platform.FindCross(crossId, query)
	if err != nil {
		return err
	}
	cross.Updated = nil

	arg := map[string]interface{}{
		"To":     to,
		"Cross":  cross,
		"Config": c.config,
	}
	text, err := GenerateContent(c.localTemplate, "cross_remind", to.Provider, to.Language, arg)
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

	text, err := GenerateContent(c.localTemplate, "cross_invitation", to.Provider, to.Language, invitation)
	if err != nil {
		return err
	}
	_, err = c.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

func (c Cross) V3Update(updates []model.CrossUpdate) error {
	logger.DEBUG("enter update")
	if len(updates) == 0 {
		return fmt.Errorf("len(updates) == 0")
	}

	to := updates[0].To
	if to.SameUser(&updates[0].By) {
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
	text, err := GenerateContent(c.localTemplate, "cross_update", to.Provider, to.Language, arg)
	if err != nil {
		return err
	}
	_, err = c.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

func (c Cross) V3Conversation(updates []model.ConversationUpdate) error {
	arg, err := ArgFromConversations(updates, c.config, c.platform)
	if err != nil {
		return err
	}
	needSend := false
	to := arg.To
	for _, update := range updates {
		if !to.SameUser(&update.Post.By) {
			needSend = true
		}
	}
	if !needSend {
		c.config.Log.Debug("not send with all self updates: %s", to)
		return nil
	}

	oldPosts, err := c.platform.GetConversation(arg.Cross.Exfee.ID, to.Token, arg.Posts[0].CreatedAt, false, "older", 2)
	if err != nil {
		logger.ERROR("get conversation error: %s", err)
	} else {
		for i := len(oldPosts) - 1; i >= 0; i-- {
			arg.OldPosts = append(arg.OldPosts, oldPosts[i])
		}
	}

	text, err := GenerateContent(c.localTemplate, "cross_conversation", to.Provider, to.Language, arg)
	if err != nil {
		return err
	}
	_, err = c.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

type ConversationArg struct {
	model.ThirdpartTo
	Cross    model.Cross
	OldPosts []model.Post
	Posts    []*model.Post
}

func ArgFromConversations(updates []model.ConversationUpdate, config *model.Config, platform *broker.Platform) (*ConversationArg, error) {
	if updates == nil && len(updates) == 0 {
		return nil, fmt.Errorf("no update info")
	}

	to := updates[0].To
	posts := make([]*model.Post, len(updates))

	for i, update := range updates {
		if !to.Equal(&update.To) {
			return nil, fmt.Errorf("updates not send to same recipient: %s, %s", to, update.To)
		}
		posts[i] = &updates[i].Post
	}

	crossId := updates[0].CrossId

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", to.UserID))
	cross, err := platform.FindCross(crossId, query)
	if err != nil {
		return nil, err
	}

	ret := &ConversationArg{
		Cross: cross,
		Posts: posts,
	}
	ret.To = to
	err = ret.Parse(config)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a ConversationArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a ConversationArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return "+00:00"
}

func (a ConversationArg) Bys() []*model.Identity {
	var ret []*model.Identity
	for _, post := range a.Posts {
		isSame := false
		for _, i := range ret {
			if i.SameUser(post.By) {
				isSame = true
				break
			}
		}
		if !isSame {
			ret = append(ret, &post.By)
		}
	}
	return ret
}

func in(id *model.Invitation, ids []model.Invitation) bool {
	for _, i := range ids {
		if id.Identity.SameUser(i.Identity) {
			return true
		}
	}
	return false
}

type SummaryArg struct {
	model.ThirdpartTo
	OldCross *model.Cross     `json:"-"`
	Cross    model.Cross      `json:"-"`
	Bys      []model.Identity `json:"-"`

	NewInvited  []model.Identity `json:"-"`
	Removed     []model.Identity `json:"-"`
	NewAccepted []model.Identity `json:"-"`
	OldAccepted []model.Identity `json:"-"`
	NewDeclined []model.Identity `json:"-"`
	NewPending  []model.Identity `json:"-"`
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

		NewInvited:  make([]model.Identity, 0),
		Removed:     make([]model.Identity, 0),
		NewAccepted: make([]model.Identity, 0),
		OldAccepted: make([]model.Identity, 0),
		NewDeclined: make([]model.Identity, 0),
		NewPending:  make([]model.Identity, 0),
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
			ret.NewAccepted = append(ret.NewAccepted, i.Identity)
		} else {
			ret.OldAccepted = append(ret.OldAccepted, i.Identity)
		}
	}
	for _, i := range ret.Cross.Exfee.Declined {
		if !in(&i, ret.OldCross.Exfee.Declined) {
			ret.NewDeclined = append(ret.NewDeclined, i.Identity)
		}
	}
	for _, i := range ret.Cross.Exfee.Pending {
		if !in(&i, ret.OldCross.Exfee.Pending) {
			ret.NewDeclined = append(ret.NewPending, i.Identity)
		}
	}
	for _, i := range ret.Cross.Exfee.Invitations {
		if !in(&i, ret.OldCross.Exfee.Invitations) {
			ret.NewInvited = append(ret.NewInvited, i.Identity)
		}
	}
	for _, i := range ret.OldCross.Exfee.Invitations {
		if !in(&i, ret.Cross.Exfee.Invitations) {
			ret.Removed = append(ret.Removed, i.Identity)
		}
	}
	return ret, nil
}

func (a *SummaryArg) NeedShowBy() bool {
	return true
}

func (a *SummaryArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
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

	return false
}

func (a *SummaryArg) IsExfeeChanged() bool {
	peopleChanged := len(a.NewAccepted)
	peopleChanged += len(a.NewDeclined)
	peopleChanged += len(a.NewInvited)
	peopleChanged += len(a.Removed)
	if peopleChanged > 0 {
		return true
	}
	return false
}

func (a *SummaryArg) IsRsvpComboChanged() []model.Identity {
	count := 0
	var ret []model.Identity
	if len(a.NewAccepted) > 0 {
		count++
		ret = append(ret, a.NewAccepted...)
	}
	if len(a.NewDeclined) > 0 {
		count++
		ret = append(ret, a.NewDeclined...)
	}
	if len(a.NewPending) > 0 {
		count++
		ret = append(ret, a.NewPending...)
	}
	if len(a.Removed) > 0 {
		count++
		ret = append(ret, a.Removed...)
	}
	if count > 1 {
		return ret
	}
	return nil
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
