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
	"strconv"
)

type Cross struct {
	rest.Service `prefix:"/v3/notifier/cross"`

	digest           rest.SimpleNode `route:"/digest" method:"POST"`
	remind           rest.SimpleNode `route:"/remind" method:"POST"`
	invitation       rest.SimpleNode `route:"/invitation" method:"POST"`
	join             rest.SimpleNode `route:"/join" method:"POST"`
	preview          rest.SimpleNode `route:"/preview" method:"POST"`
	update           rest.SimpleNode `route:"/update" method:"POST"`
	updateInvitation rest.SimpleNode `route:"/update_invitation" method:"POST"`
	conversation     rest.SimpleNode `route:"/conversation" method:"POST"`

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

func (c Cross) Digest(ctx rest.Context, requests []model.CrossDigestRequest) {
	if len(requests) == 0 {
		ctx.Return(http.StatusBadRequest, "len(requests) == 0")
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
		ctx.Return(http.StatusBadRequest, err)
		return
	}

	conversations, err := c.platform.GetConversation(cross.Exfee.ID, updatedAt, false, "newer", -1)
	if err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	foldedConversation := 0
	if l := len(conversations); l > 3 {
		foldedConversation = l
		conversations = conversations[0:3]
	}

	weatherIcon := ""
	if t, err := cross.Time.BeginAt.UTCTime("2006-01-02 15:04:05"); err == nil && cross.Place != nil {
		lat, err := strconv.ParseFloat(cross.Place.Lat, 64)
		if err == nil {
			lng, err := strconv.ParseFloat(cross.Place.Lng, 64)
			if err == nil {
				weatherIcon = c.platform.GetWeatherIcon(lat, lng, t)
			}
		}
	}

	arg := map[string]interface{}{
		"To":                 to,
		"Cross":              cross,
		"Config":             c.config,
		"FoldedConversation": foldedConversation,
		"Conversations":      conversations,
		"WeatherIcon":        weatherIcon,
	}

	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_digest", c.domain+"/v3/notifier/cross/digest", &failArg)
	ctx.Return(http.StatusAccepted)
}

func (c Cross) Remind(ctx rest.Context, requests []model.CrossDigestRequest) {
	if len(requests) == 0 {
		ctx.Return(http.StatusBadRequest, "len(requests) == 0")
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
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	cross.Updated = nil

	weatherIcon := ""
	if t, err := cross.Time.BeginAt.UTCTime("2006-01-02 15:04:05"); err == nil && cross.Place != nil {
		lat, err := strconv.ParseFloat(cross.Place.Lat, 64)
		if err == nil {
			lng, err := strconv.ParseFloat(cross.Place.Lng, 64)
			if err == nil {
				weatherIcon = c.platform.GetWeatherIcon(lat, lng, t)
			}
		}
	}

	arg := map[string]interface{}{
		"To":          to,
		"Cross":       cross,
		"Config":      c.config,
		"WeatherIcon": weatherIcon,
	}
	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_remind", c.domain+"/v3/notifier/cross/remind", &failArg)
	ctx.Return(http.StatusAccepted)
}

type InvitationArg struct {
	To      model.Recipient `json:"to"`
	By      model.Identity  `json:"by"`
	CrossId int64           `json:"cross_id"`
	Cross   model.Cross     `json:"cross"`

	Config      *model.Config `json:"-"`
	WeatherIcon string        `json:"-"`
}

func (a InvitationArg) String() string {
	return fmt.Sprintf("{to:%s by:%s cross:%d}", a.To, a.By, a.Cross.ID)
}

func (a *InvitationArg) Parse(config *model.Config, platform *broker.Platform) (err error) {
	a.Config = config

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", a.To.UserID))
	cross, err := platform.FindCross(a.CrossId, query)
	if err != nil {
		return err
	}
	a.Cross = cross

	if t, err := cross.Time.BeginAt.UTCTime("2006-01-02 15:04:05"); err == nil && cross.Place != nil {
		lat, err := strconv.ParseFloat(cross.Place.Lat, 64)
		if err == nil {
			lng, err := strconv.ParseFloat(cross.Place.Lng, 64)
			if err == nil {
				a.WeatherIcon = platform.GetWeatherIcon(lat, lng, t)
			}
		}
	}

	return nil
}

func (a InvitationArg) ToIn(invitations []model.Invitation) bool {
	for _, i := range invitations {
		if a.To.SameUser(&i.Identity) {
			return true
		}
	}
	return false
}

func (a InvitationArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a InvitationArg) PublicLink() string {
	return fmt.Sprintf("%s/#!%d/%s", a.Config.SiteUrl, a.Cross.ID, a.To.Token[1:5])
}

func (a InvitationArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
}

func (a InvitationArg) SendToBy() bool {
	return a.To.SameUser(&a.By)
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
		ret += "..."
	}
	return ret
}

func (c Cross) Invitation(ctx rest.Context, invitation InvitationArg) {
	if invitation.SendToBy() {
		ctx.Return(http.StatusBadRequest, "not send to self")
		return
	}
	if err := invitation.Parse(c.config, c.platform); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	to := &invitation.To

	go SendAndSave(c.localTemplate, c.platform, to, invitation, "cross_invitation", c.domain+"/v3/notifier/cross/invitation", &invitation)
	ctx.Return(http.StatusAccepted)
}

type JoinArg struct {
	To      model.Recipient `json:"to"`
	Invitee model.Identity  `json:"invitee"`
	By      model.Identity  `json:"by"`
	CrossId int64           `json:"cross_id"`

	Cross  model.Cross   `json:"-"`
	Config *model.Config `json:"-"`
}

func (a *JoinArg) Parse(config *model.Config, platform *broker.Platform) (err error) {
	if a.SendToInvitee() {
		return fmt.Errorf("not send to invitee")
	}
	a.Config = config

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", a.To.UserID))
	cross, err := platform.FindCross(a.CrossId, query)
	if err != nil {
		return err
	}
	a.Cross = cross

	return nil
}

func (a JoinArg) SendToInvitee() bool {
	return a.To.SameUser(&a.Invitee)
}

func (a JoinArg) SendToBy() bool {
	return a.To.SameUser(&a.By)
}

func (c Cross) Join(ctx rest.Context, arg JoinArg) {
	if err := arg.Parse(c.config, c.platform); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	to := &arg.To

	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_join", c.domain+"/v3/notifier/cross/arg", &arg)
	ctx.Return(http.StatusAccepted)
}

func (c Cross) Preview(ctx rest.Context, invitation InvitationArg) {
	if err := invitation.Parse(c.config, c.platform); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	to := &invitation.To

	go SendAndSave(c.localTemplate, c.platform, to, invitation, "cross_preview", c.domain+"/v3/notifier/cross/preview", &invitation)
	ctx.Return(http.StatusAccepted)
}

func (c Cross) Update(ctx rest.Context, updates []model.CrossUpdate) {
	if len(updates) == 0 {
		ctx.Return(http.StatusBadRequest, "len(updates) == 0")
		return
	}

	to := &updates[0].To
	if to.SameUser(&updates[0].By) {
		logger.DEBUG("not send with all self updates: %s", to)
		return
	}

	arg, err := updateFromUpdates(updates, c.config)
	if err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}

	if !arg.IsChanged() {
		ctx.Return(http.StatusBadRequest, "not changed")
		return
	}

	failArg := updates[len(updates)-1:]
	failArg[0].OldCross = updates[0].OldCross
	to = &failArg[0].To

	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_update", c.domain+"/v3/notifier/cross/update", &failArg)
	ctx.Return(http.StatusAccepted)
}

type UpdateInvitationArg struct {
	To       model.Recipient  `json:"to"`
	Invitees []model.Identity `json:"invitees"`
	By       model.Identity   `json:"by"`
	CrossId  int64            `json:"cross_id"`

	Cross  model.Cross   `json:"-"`
	Config *model.Config `json:"-"`
}

func (a UpdateInvitationArg) HasMany() bool {
	return len(a.Invitees) > 1
}

func (a UpdateInvitationArg) SendToInvitee() bool {
	for _, invitee := range a.Invitees {
		if a.To.SameUser(&invitee) {
			return true
		}
	}
	return false
}

func (c Cross) UpdateInvitation(ctx rest.Context, arg UpdateInvitationArg) {
	if len(arg.Invitees) == 0 {
		ctx.Return(http.StatusBadRequest, "len(invitees) == 0")
		return
	}
	if arg.SendToInvitee() {
		ctx.Return(http.StatusNoContent, "not send to invitee")
		return
	}

	arg.Config = c.config

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", arg.To.UserID))
	cross, err := c.platform.FindCross(arg.CrossId, query)
	if err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	arg.Cross = cross

	to := &arg.To

	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_update_invitation", c.domain+"/v3/notifier/cross/update_invitation", &arg)
	ctx.Return(http.StatusAccepted)
}

func (c Cross) Conversation(ctx rest.Context, updates []model.ConversationUpdate) {
	failArg := updates
	arg, err := ArgFromConversations(updates, c.config, c.platform)
	if err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	needSend := false
	to := &failArg[0].To
	for _, update := range updates {
		if !to.SameUser(&update.Post.By) {
			needSend = true
			break
		}
	}
	if !needSend {
		logger.DEBUG("not send with all self updates: %s", to)
		return
	}

	oldPosts, err := c.platform.GetConversation(arg.Cross.Exfee.ID, arg.Posts[0].CreatedAt, false, "older", 2)
	if err != nil {
		logger.ERROR("get conversation error: %s", err)
	} else {
		for i := len(oldPosts) - 1; i >= 0; i-- {
			arg.OldPosts = append(arg.OldPosts, oldPosts[i])
		}
	}

	go SendAndSave(c.localTemplate, c.platform, to, arg, "cross_conversation", c.domain+"/v3/notifier/cross/conversation", &failArg)
	ctx.Return(http.StatusAccepted)
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

func (a ConversationArg) HasMany() bool {
	return len(a.Posts) > 1
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
		if in(&i, ret.OldCross.Exfee.Invitations) && !in(&i, ret.OldCross.Exfee.Pending) {
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
	peopleChanged += len(a.NewPending)
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

func (a *UpdateArg) IsCrossComboChanged() bool {
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

func (a *UpdateArg) IsComboChanged() bool {
	changedNumber := 0
	if a.IsTitleChanged() {
		changedNumber++
	}
	if a.IsTimeChanged() {
		changedNumber++
	}
	if a.IsPlaceChanged() {
		changedNumber++
	}
	if a.IsDescriptionChanged() {
		changedNumber++
	}
	if a.IsExfeeChanged() {
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
