package exfe_service

import (
	"log"
	"github.com/googollee/godis"
	"gobus"
	"exfe/model"
	"fmt"
	"time"
	"os"
)

type UpdateCrossArg struct {
	Cross exfe_model.Cross
	To_identities []exfe_model.Identity
	By_identity exfe_model.Identity

	Old_cross *exfe_model.Cross
	Post *exfe_model.Post
}

type OneIdentityUpdateArg struct {
	Cross exfe_model.Cross
	To_identity exfe_model.Identity
	By_identity exfe_model.Identity

	Old_cross *exfe_model.Cross
	Post *exfe_model.Post
}

type Cross struct {
	queues map[string]*gobus.TailDelayQueue
	config *Config
	log *log.Logger
}

func NewCross(config *Config) *Cross {
	arg := []OneIdentityUpdateArg{}
	redis := godis.New(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)
	log := log.New(os.Stderr, "exfe.cross", log.LstdFlags)
	queues := make(map[string]*gobus.TailDelayQueue)
	for _, p := range [...]string{"twitter", "push", "email"} {
		queue, err := gobus.NewTailDelayQueue(getProviderQueueName(p), config.Cross.Delay[p], arg, redis)
		if err != nil {
			panic(err)
		}
		queues[p] = queue
	}
	return &Cross{
		queues: queues,
		config: config,
		log: log,
	}
}

func (s *Cross) Update(args []*UpdateCrossArg) error {
	for _, arg := range args {
		for _, to := range arg.To_identities {
			update := OneIdentityUpdateArg{
				Cross: arg.Cross,
				To_identity: to,
				By_identity: arg.By_identity,
				Old_cross: arg.Old_cross,
				Post: arg.Post,
			}
			s.dispatch(&update)
		}
	}
	return nil
}

func (s *Cross) getUserIdentityMap(cross *exfe_model.Cross) (identityMap map[uint64]*exfe_model.Identity, userMap map[uint64]*exfe_model.Identity) {
	identityMap = make(map[uint64]*exfe_model.Identity)
	userMap = make(map[uint64]*exfe_model.Identity)

	for _, invitation := range cross.Exfee.Invitations {
		identityMap[invitation.Identity.Id] = &invitation.Identity
		userMap[invitation.Identity.Connected_user_id] = &invitation.Identity
	}
	return
}

func (s *Cross) dispatch(arg *OneIdentityUpdateArg) {
	id := fmt.Sprintf("%d-%d", arg.Cross.Id, arg.To_identity.Id)

	queue, ok := s.queues[arg.To_identity.Provider]
	if !ok {
		if arg.To_identity.Provider == "iOSAPN" || arg.To_identity.Provider == "Android" {
			queue, ok = s.queues["push"]
		}
	}
	if !ok {
		log.Printf("Not support provider: %s", arg.To_identity.Provider)
		return
	}
	if arg.To_identity.Provider != "email" {
		if arg.Post != nil {
			log.Printf("provider %s can't handle post now", arg.To_identity.Provider)
			return
		}
	}
	queue.Push(id, arg)
}

func getProviderQueueName(provider string) string{
	return fmt.Sprintf("exfe:queue:cross:%s", provider)
}

type CrossProviderHandler interface{
	Handle(to_identity *exfe_model.Identity, old_cross, cross *exfe_model.Cross)
}

type CrossProviderBase struct {
	log *log.Logger
	queue *gobus.TailDelayQueue
	config *Config
	client *gobus.Client
	handler CrossProviderHandler
}

func NewCrossProviderBase(provider string, config *Config) (ret CrossProviderBase) {
	ret.log = log.New(os.Stderr, fmt.Sprintf("exfe.cross.%s", provider), log.LstdFlags)

	arg := []OneIdentityUpdateArg{}
	redis := godis.New(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)
	var err error
	ret.queue, err = gobus.NewTailDelayQueue(getProviderQueueName(provider), config.Cross.Delay[provider], arg, redis)
	if err != nil {
		panic(err)
	}

	ret.config = config
	ret.client = gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, provider)
	return
}

func (b *CrossProviderBase) Serve() {
	for {
		t, err := b.queue.NextWakeup()
		if err != nil {
			log.Printf("next wakeup error: %s", err)
			break
		}
		time.Sleep(t)
		args, err := b.queue.Pop()
		if err != nil {
			log.Printf("pop from delay queue failed: %s", err)
			continue
		}
		if args != nil {
			updates := args.([]OneIdentityUpdateArg)
			old_cross := updates[0].Old_cross
			cross := &updates[len(updates)-1].Cross
			to_identity := &updates[0].To_identity

			b.handler.Handle(to_identity, old_cross, cross)
		}
	}
}

func findToken(to *exfe_model.Identity, cross *exfe_model.Cross) (ret *string) {
	for _, invitation := range cross.Exfee.Invitations {
		if invitation.Identity.Connected_user_id == to.Connected_user_id {
			ret = &invitation.Token
			break
		}
	}
	return
}

func diffExfee(log *log.Logger, old, new_ *exfe_model.Exfee) (accepted map[uint64]*exfe_model.Identity, declined map[uint64]*exfe_model.Identity, newlyInvited map[uint64]*exfe_model.Invitation, removed map[uint64]*exfe_model.Identity) {
	oldId := make(map[uint64]*exfe_model.Invitation)
	newId := make(map[uint64]*exfe_model.Invitation)

	accepted = make(map[uint64]*exfe_model.Identity)
	declined = make(map[uint64]*exfe_model.Identity)
	newlyInvited = make(map[uint64]*exfe_model.Invitation)
	removed = make(map[uint64]*exfe_model.Identity)

	for i, v := range old.Invitations {
		if v.Rsvp_status == "NOTIFICATION" {
			continue
		}
		if _, ok := oldId[v.Identity.Connected_user_id]; ok {
			log.Printf("more than one non-notification status in exfee %d, user id %d", old.Id, v.Identity.Connected_user_id)
		}
		oldId[v.Identity.Connected_user_id] = &old.Invitations[i]
	}
	for i, v := range new_.Invitations {
		if v.Rsvp_status == "NOTIFICATION" {
			continue
		}
		if _, ok := newId[v.Identity.Connected_user_id]; ok {
			log.Printf("more than one non-notification status in exfee %d, user id %d", old.Id, v.Identity.Connected_user_id)
		}
		newId[v.Identity.Connected_user_id] = &new_.Invitations[i]
	}

	for k, v := range newId {
		switch v.Rsvp_status {
		case "ACCEPTED":
			if inv, ok := oldId[k]; !ok || inv.Rsvp_status != v.Rsvp_status {
				accepted[k] = &v.Identity
			}
		case "DECLINED":
			if inv, ok := oldId[k]; !ok || inv.Rsvp_status != v.Rsvp_status {
				declined[k] = &v.Identity
			}
		}
		if _, ok := oldId[k]; !ok {
			newlyInvited[k] = v
		}
	}
	for k, v := range oldId {
		if _, ok := newId[k]; !ok {
			removed[k] = &v.Identity
		}
	}
	return
}

type NewInvitationData struct {
	ToUserName    string
	IsHost        bool
	Title         string
	Time          string
	Place         string
	SiteUrl       string
	Token         string
}

func newInvitationData(log *log.Logger, siteUrl string, to *exfe_model.Identity, cross *exfe_model.Cross) *NewInvitationData {
	t, err := cross.Time.StringInZone(to.Timezone)
	if err != nil {
		log.Printf("Time parse error: %s", err)
		return nil
	}
	isHost := cross.By_identity.Connected_user_id == to.Connected_user_id
	return &NewInvitationData{
		ToUserName:    to.External_username,
		IsHost:        isHost,
		Title:         cross.Title,
		Time:          t,
		Place:         cross.Place.String(),
		SiteUrl:       siteUrl,
		Token:         *findToken(to, cross),
	}
}
