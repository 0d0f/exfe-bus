package notifier

import (
	"broker"
	"fmt"
	"formatter"
	"github.com/googollee/go-rest"
	"logger"
	"model"
	"net/http"
	"net/url"
)

type Cross struct {
	rest.Service `prefix:"/v3/notifier/cross"`

	Digest       rest.Processor `path:"/digest" method:"POST"`
	Remind       rest.Processor `path:"/remind" method:"POST"`
	Invitation   rest.Processor `path:"/invitation" method:"POST"`
	Preview      rest.Processor `path:"/preview" method:"POST"`
	Update       rest.Processor `path:"/update" method:"POST"`
	Conversation rest.Processor `path:"/conversation" method:"POST"`

	localTemplate *formatter.LocalTemplate
	config        *model.Config
	platform      *broker.Platform
	domain        string
}

func NewCross(localTemplate *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) *Cross {
	return &Cross{
		localTemplate: localTemplate,
		config:        config,
		platform:      platform,
		domain:        fmt.Sprintf("http://%s:%d", config.ExfeService.Addr, config.ExfeService.Port),
	}
}

func (c Cross) HandleDigest(requests []model.CrossDigestRequest) {
	if len(requests) == 0 {
		c.Error(http.StatusBadRequest, fmt.Errorf("len(requests) == 0"))
		return
	}
	crossId := requests[0].CrossId
	updatedAt := requests[0].UpdatedAt

	failArg := requests[len(requests)-1:]
	failArg[0].CrossId = crossId
	failArg[0].UpdatedAt = updatedAt

	to := &failArg[0].To

	query := make(url.Values)
	query.Set("updated_at", updatedAt)
	query.Set("user_id", fmt.Sprintf("%d", to.UserID))
	cross, err := c.platform.FindCross(crossId, query)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	arg := map[string]interface{}{
		"To":     to,
		"Cross":  cross,
		"Config": c.config,
	}

	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_digest", c.domain+"/v3/notifier/cross/digest", &failArg)
	c.WriteHeader(http.StatusAccepted)
}

func (c Cross) HandleRemind(requests []model.CrossDigestRequest) {
	if len(requests) == 0 {
		c.Error(http.StatusBadRequest, fmt.Errorf("len(requests) == 0"))
		return
	}
	crossId := requests[0].CrossId
	failArg := requests[len(requests)-1:]
	failArg[0].CrossId = crossId
	to := &failArg[0].To

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", to.UserID))
	cross, err := c.platform.FindCross(crossId, query)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}
	cross.Updated = nil

	arg := map[string]interface{}{
		"To":     to,
		"Cross":  cross,
		"Config": c.config,
	}
	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_remind", c.domain+"/v3/notifier/cross/remind", &failArg)
	c.WriteHeader(http.StatusAccepted)
}

func (c Cross) HandleInvitation(invitation model.CrossInvitation) {
	invitation.Config = c.config
	to := &invitation.To

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", to.UserID))
	cross, err := c.platform.FindCross(invitation.CrossId, query)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}
	invitation.Cross = cross

	go SendAndSave(c.localTemplate, c.platform, to, invitation, "cross_invitation", c.domain+"/v3/notifier/cross/invitation", &invitation)
	c.WriteHeader(http.StatusAccepted)
}

func (c Cross) HandlePreview(invitation model.CrossInvitation) {
	invitation.Config = c.config
	to := &invitation.To

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", to.UserID))
	cross, err := c.platform.FindCross(invitation.CrossId, query)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}
	invitation.Cross = cross

	go SendAndSave(c.localTemplate, c.platform, to, invitation, "cross_preview", c.domain+"/v3/notifier/cross/preview", &invitation)
	c.WriteHeader(http.StatusAccepted)
}

func (c Cross) HandleUpdate(updates []model.CrossUpdate) {
	if len(updates) == 0 {
		c.Error(http.StatusBadRequest, fmt.Errorf("len(updates) == 0"))
		return
	}

	to := &updates[0].To
	if to.SameUser(&updates[0].By) {
		logger.DEBUG("not send with all self updates: %s", to)
		return
	}

	arg, err := updateFromUpdates(updates, c.config)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	if !arg.IsChanged() {
		c.Error(http.StatusBadRequest, fmt.Errorf("not changed"))
		return
	}

	failArg := updates[len(updates)-1:]
	failArg[0].OldCross = updates[0].OldCross
	to = &failArg[0].To

	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_update", c.domain+"/v3/notifier/cross/update", &failArg)
	c.WriteHeader(http.StatusAccepted)
}

func (c Cross) HandleConversation(updates []model.ConversationUpdate) {
	failArg := updates
	arg, err := ArgFromConversations(updates, c.config, c.platform)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}
	needSend := false
	to := &failArg[0].To
	for _, update := range updates {
		if !to.SameUser(&update.Post.By) {
			needSend = true
		}
	}
	if !needSend {
		logger.DEBUG("not send with all self updates: %s", to)
		return
	}

	oldPosts, err := c.platform.GetConversation(arg.Cross.Exfee.ID, to.Token, arg.Posts[0].CreatedAt, false, "older", 2)
	if err != nil {
		logger.ERROR("get conversation error: %s", err)
	} else {
		for i := len(oldPosts) - 1; i >= 0; i-- {
			arg.OldPosts = append(arg.OldPosts, oldPosts[i])
		}
	}

	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_conversation", c.domain+"/v3/notifier/cross/conversation", &failArg)
	c.WriteHeader(http.StatusAccepted)
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

type UpdateArg struct {
	model.ThirdpartTo
	OldCross *model.Cross     `json:"-"`
	Cross    model.Cross      `json:"-"`
	Bys      []model.Identity `json:"-"`

	TitleChangedBy       *model.Identity
	DescriptionChangedBy *model.Identity
	TimeChangedBy        *model.Identity
	PlaceChangedBy       *model.Identity

	NewInvited  []model.Identity `json:"-"`
	Removed     []model.Identity `json:"-"`
	NewAccepted []model.Identity `json:"-"`
	OldAccepted []model.Identity `json:"-"`
	NewDeclined []model.Identity `json:"-"`
	NewPending  []model.Identity `json:"-"`
}

func updateFromUpdates(updates []model.CrossUpdate, config *model.Config) (*UpdateArg, error) {
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

	ret := &UpdateArg{
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

	crossTime, _ := ret.Cross.Time.StringInZone(ret.To.Timezone)
	for i := len(updates) - 1; i >= 0; i-- {
		c := updates[i].OldCross
		if ret.TitleChangedBy == nil && c.Title != ret.Cross.Title {
			ret.TitleChangedBy = &c.By
		}
		if ret.DescriptionChangedBy == nil && c.Description != ret.Cross.Description {
			ret.DescriptionChangedBy = &c.By
		}
		t, _ := c.Time.StringInZone(ret.To.Timezone)
		if ret.TimeChangedBy == nil && t != crossTime {
			ret.TimeChangedBy = &c.By
		}
		if ret.PlaceChangedBy == nil && !c.Place.Same(ret.Cross.Place) {
			ret.PlaceChangedBy = &c.By
		}
	}

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
			ret.NewPending = append(ret.NewPending, i.Identity)
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

func (a *UpdateArg) NeedShowBy() bool {
	return true
}

func (a *UpdateArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
}

func (a *UpdateArg) IsChanged() bool {
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

func (a *UpdateArg) IsExfeeChanged() bool {
	peopleChanged := len(a.NewAccepted)
	peopleChanged += len(a.NewDeclined)
	peopleChanged += len(a.NewInvited)
	peopleChanged += len(a.Removed)
	if peopleChanged > 0 {
		return true
	}
	return false
}

func (a *UpdateArg) IsResponseComboChanged() []model.Identity {
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

func (a *UpdateArg) IsTimeChanged() bool {
	return a.TimeChangedBy != nil
}

func (a *UpdateArg) IsTitleChanged() bool {
	return a.TitleChangedBy != nil
}

func (a *UpdateArg) IsPlaceChanged() bool {
	return a.PlaceChangedBy != nil
}

func (a *UpdateArg) IsPlaceTitleChanged() bool {
	return a.Cross.Place.Title != a.OldCross.Place.Title
}

func (a *UpdateArg) IsPlaceDescChanged() bool {
	if a.Cross.Place.Description != a.OldCross.Place.Description {
		return true
	}
	if a.Cross.Place.Lng != a.OldCross.Place.Lng {
		return true
	}
	if a.Cross.Place.Lat != a.OldCross.Place.Lat {
		return true
	}
	return false
}

func (a *UpdateArg) IsDescriptionChanged() bool {
	return a.DescriptionChangedBy != nil
}

func (a *UpdateArg) IsComboChanged() bool {
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

func (a UpdateArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a UpdateArg) PublicLink() string {
	return fmt.Sprintf("%s/#!%d/%s", a.Config.SiteUrl, a.Cross.ID, a.To.Token[1:5])
}
