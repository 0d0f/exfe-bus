package exfe_service

import (
	"log/syslog"
	"github.com/simonz05/godis"
	"gobus"
	"exfe/model"
	"fmt"
)

type UpdateCrossArg struct {
	Cross exfe_model.Cross
	Old_cross *exfe_model.Cross
	To_identities []exfe_model.Identity
	By_identity exfe_model.Identity
}

type OneIdentityUpdateArg struct {
	Cross exfe_model.Cross
	Old_cross *exfe_model.Cross
	To_identity exfe_model.Identity
	By_identity exfe_model.Identity
}

type Cross struct {
	twitterQueue *gobus.TailDelayQueue
	config *Config
	log *syslog.Writer
}

func NewCross(config *Config) *Cross {
	arg := []OneIdentityUpdateArg{}
	redis := godis.New(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)
	log, err := syslog.New(syslog.LOG_DEBUG, "exfe.cross")
	if err != nil {
		panic(err)
	}
	queue, err := gobus.NewTailDelayQueue(crossTwitterQueueName, config.Cross.Twitter_delay, arg, redis)
	if err != nil {
		panic(err)
	}
	return &Cross{
		twitterQueue: queue,
		config: config,
		log: log,
	}
}

func (s *Cross) Update(args []*UpdateCrossArg) error {
	for _, arg := range args {
		for _, to := range arg.To_identities {
			update := OneIdentityUpdateArg{
				Cross: arg.Cross,
				Old_cross: arg.Old_cross,
				To_identity: to,
				By_identity: arg.By_identity,
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

	switch arg.To_identity.Provider {
	case "twitter":
		s.twitterQueue.Push(id, arg)
	default:
		s.log.Err(fmt.Sprintf("Not support provider: %s", arg.To_identity.Provider))
	}
}
