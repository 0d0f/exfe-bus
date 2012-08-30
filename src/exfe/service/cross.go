package exfe_service

import (
	"exfe/model"
	"fmt"
	"github.com/googollee/godis"
	"gobus"
	"log"
	"os"
	"time"
)

type UpdateCrossArg struct {
	Cross         exfe_model.Cross
	To_identities []exfe_model.Identity
	By_identity   exfe_model.Identity

	Old_cross *exfe_model.Cross
	Post      *exfe_model.Post
}

type OneIdentityUpdateArg struct {
	Cross       exfe_model.Cross
	To_identity exfe_model.Identity
	By_identity exfe_model.Identity

	Old_cross *exfe_model.Cross
	Post      *exfe_model.Post
}

type Cross struct {
	config *Config
	log    *log.Logger
	post   *CrossPost
}

func NewCross(config *Config) *Cross {
	log := log.New(os.Stderr, "exfe.cross", log.LstdFlags)
	return &Cross{
		config: config,
		log:    log,
		post:   NewCrossPost(config),
	}
}

func (s *Cross) Update(args []*UpdateCrossArg) error {
	s.log.Printf("received %d updates", len(args))
	for _, arg := range args {
		for _, to := range arg.To_identities {
			update := OneIdentityUpdateArg{
				Cross:       arg.Cross,
				To_identity: to,
				By_identity: arg.By_identity,
				Old_cross:   arg.Old_cross,
				Post:        arg.Post,
			}
			s.dispatch(&update)
		}
	}
	return nil
}

func (s *Cross) getUserIdentityMap(cross *exfe_model.Cross) (identityMap map[int64]*exfe_model.Identity, userMap map[int64]*exfe_model.Identity) {
	identityMap = make(map[int64]*exfe_model.Identity)
	userMap = make(map[int64]*exfe_model.Identity)

	for _, invitation := range cross.Exfee.Invitations {
		identityMap[invitation.Identity.Id] = &invitation.Identity
		userMap[invitation.Identity.Connected_user_id] = &invitation.Identity
	}
	return
}

func (s *Cross) dispatch(arg *OneIdentityUpdateArg) {
	id := fmt.Sprintf("%d-%d", arg.Cross.Id, arg.To_identity.Id)

	queueName := arg.To_identity.Provider
	switch arg.To_identity.Provider {
	case "iOS":
		fallthrough
	case "Android":
		queueName = "push"
	case "facebook":
		queueName = "email"
		arg.To_identity.External_id = fmt.Sprintf("%s@facebook.com", arg.To_identity.External_username)
	}
	if queueName != "push" && queueName != "email" && queueName != "twitter" {
		log.Printf("Not support provider: %s", arg.To_identity.Provider)
		return
	}

	if arg.To_identity.Provider != "email" {
		if arg.Post != nil {
			log.Printf("push post to %s@%s", arg.To_identity.External_id, arg.To_identity.Provider)
			s.post.SendPost(arg)
			return
		}
	}

	redis := godis.New(s.config.Redis.Netaddr, s.config.Redis.Db, s.config.Redis.Password)
	queue, err := gobus.NewTailDelayQueue(getProviderQueueName(queueName), s.config.Cross.Delay[queueName], []OneIdentityUpdateArg{}, redis)
	if err != nil {
		log.Printf("can't connect to redis: %s", err)
		return
	}
	err = queue.Push(id, arg)
	if err != nil {
		log.Printf("can't push to queue(%s): %s", queue, err)
	}
	s.log.Printf("dispatch %s to %s", id, queue)
}

func getProviderQueueName(provider string) string {
	return fmt.Sprintf("exfe:queue:cross:%s", provider)
}

type CrossProviderHandler interface {
	Handle(arg *ProviderArg)
}

type CrossProviderBase struct {
	log     *log.Logger
	queue   *gobus.TailDelayQueue
	config  *Config
	client  *gobus.Client
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
			arg := &ProviderArg{
				Old_cross:   updates[0].Old_cross,
				Cross:       &updates[len(updates)-1].Cross,
				To_identity: &updates[0].To_identity,
			}

			log.Printf("Got %d updates to %s", len(updates), updates[0].To_identity.ExternalId())
			b.handler.Handle(arg)
		}
	}
}
