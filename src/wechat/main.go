package main

import (
	"broker"
	"daemon"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/googollee/go-aws/s3"
	"logger"
	"model"
	"os"
	"strconv"
	"time"
)

func main() {
	var config model.Config
	_, quit := daemon.Init("exfe.json", &config)

	workType := os.Args[len(os.Args)-1]
	work, ok := config.Wechat[workType]
	if !ok {
		logger.ERROR("unknow work type %s", workType)
		return
	}

	aws := s3.New(config.AWS.S3.Domain, config.AWS.S3.Key, config.AWS.S3.Secret)
	aws.SetACL(s3.ACLPublicRead)
	aws.SetLocationConstraint(s3.LC_AP_SINGAPORE)
	bucket, err := aws.GetBucket(fmt.Sprintf("%s-3rdpart-photos", config.AWS.S3.BucketPrefix))
	if err != nil {
		logger.ERROR("can't create bucket: %s", err)
		return
	}

	platform, err := broker.NewPlatform(&config)
	if err != nil {
		logger.ERROR("can't create platform: %s", err)
		return
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&autocommit=true",
		config.DB.Username, config.DB.Password, config.DB.Addr, config.DB.Port, config.DB.DbName))
	if err != nil {
		logger.ERROR("mysql error:", err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.ERROR("mysql error:", err)
		return
	}
	kvSaver := broker.NewKVSaver(db)

	wc, err := New(work.Username, work.Password, work.PingId, &config)
	if err != nil {
		logger.ERROR("can't create wechat: %s", err)
		return
	}
	defer func() {
		logger.NOTICE("quit")
	}()

	go runServer(work.Addr, work.Port, wc, kvSaver)

	logger.NOTICE("login as %s", wc.userName)

	for {
		select {
		case <-quit:
			return
		default:
		}
		ret, err := wc.Check()
		if err != nil {
			logger.ERROR("can't get last message: %s", err)
			return
		}
		if ret == "0" {
			err = wc.Ping(time.Minute * 30)
			if err != nil {
				logger.ERROR("ping error: %s", err)
				return
			}
			continue
		}
		resp, err := wc.GetLast()
		if err != nil {
			logger.ERROR("can't get last message: %s", err)
			return
		}
		for _, msg := range resp.AddMsgList {
			if msg.MsgType != JoinMessage {
				continue
			}
			uin, cross, err := wc.ConvertCross(bucket, &msg)
			if err != nil {
				logger.ERROR("can't convert to cross: %s", err)
				continue
			}
			uinStr := fmt.Sprintf("%d", uin)
			idStr, exist, err := kvSaver.Check([]string{uinStr})
			if err != nil {
				logger.ERROR("can't check uin %s: %s", uinStr, err)
				continue
			}
			if exist {
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					goto CREATE
				}
				oldCross, err := platform.FindCross(id, nil)
				if err != nil {
					goto CREATE
				}
				exfee := make(map[string]bool)
				host := cross.Exfee.Invitations[0].Identity
				for _, invitation := range cross.Exfee.Invitations {
					exfee[invitation.Identity.Id()] = true
					if invitation.Host {
						host = invitation.Identity
					}
				}
				for _, invitation := range oldCross.Exfee.Invitations {
					if exfee[invitation.Identity.Id()] {
						continue
					}
					invitation.Response = model.Removed
					cross.Exfee.Invitations = append(cross.Exfee.Invitations, invitation)
				}
				err = platform.BotCrossUpdate("cross_id", idStr, cross, host)
				if err != nil {
					logger.ERROR("can't update cross %s: %s", idStr, err)
					goto CREATE
				}
			}
		CREATE:
			cross, err = platform.BotCrossGather(cross)
			if err != nil {
				logger.ERROR("can't gather cross: %s", err)
				continue
			}
			err = kvSaver.Save([]string{fmt.Sprintf("%d", uin)}, fmt.Sprintf("%d", cross.ID))
			if err != nil {
				logger.ERROR("can't save cross id: %s", err)
			}
			err = kvSaver.Save([]string{fmt.Sprintf("e%d@exfe", cross.Exfee.ID)}, fmt.Sprintf("%d@chatroom", uin))
			if err != nil {
				logger.ERROR("can't save exfee id: %s", err)
			}
			logger.INFO("wechat_gather", msg.FromUserName, uin, cross.ID, cross.Exfee.ID, err)
			smith, err := cross.Exfee.FindInvitedUser(model.Identity{
				ExternalID: fmt.Sprintf("%d", wc.baseRequest.Uin),
				Provider:   "wechat",
			})
			if err != nil {
				logger.ERROR("can't find Smith Exfer in cross %d: %s", cross.ID, err)
				continue
			}
			chatroom := fmt.Sprintf("%d@chatroom", uin)
			u := fmt.Sprintf("%s/#!token=%s/routex/", config.SiteUrl, smith.Token)
			err = wc.SendMessage(chatroom, u)
			logger.NOTICE("send %s to %s", u, chatroom)
			if err != nil {
				logger.ERROR("can't send %s to %s", u, chatroom)
				continue
			}
		}
	}
}
