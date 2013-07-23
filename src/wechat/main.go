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
)

func main() {
	var config model.Config
	quit := daemon.Init("exfe.json", &config)

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

	bot := &Bot{
		platform: platform,
		config:   &config,
		bucket:   bucket,
		kvSaver:  kvSaver,
		wc:       wc,
	}

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
			// err = wc.Ping(time.Minute * 30)
			// if err != nil {
			// 	logger.ERROR("ping error: %s", err)
			// 	return
			// }
			continue
		}
		resp, err := wc.GetLast()
		if err != nil {
			logger.ERROR("can't get last message: %s", err)
			return
		}
		for _, msg := range resp.AddMsgList {
			switch msg.MsgType {
			case FriendRequest:
				bot.JoinRequest(msg)
			case JoinMessage:
				bot.Join(msg)
			case ChatMessage:
			default:
				logger.DEBUG("msg type: %d, content: %#v", msg.MsgType, msg)
			}
		}
	}
}
