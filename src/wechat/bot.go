package main

import (
	"broker"
	"fmt"
	"github.com/googollee/go-aws/s3"
	"logger"
	"math/rand"
	"model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Bot struct {
	platform *broker.Platform
	config   *model.Config
	bucket   *s3.Bucket
	kvSaver  *broker.KVSaver
	wc       *WeChat
}

func (b *Bot) JoinRequest(msg Message) {
	username, ticket, verifyFlag := msg.RecommendInfo.UserName, msg.RecommendInfo.Ticket, msg.RecommendInfo.VerifyFlag
	if verifyFlag == 0 {
		err := b.wc.Verify(username, ticket)
		if err != nil {
			logger.ERROR("can't add user %s: %s", username, err)
			return
		}
	}
	b.GreetNewFriend(username)
}

func (b *Bot) Join(msg Message) {
	if strings.HasSuffix(msg.FromUserName, "@chatroom") {
		chatroomId, cross, err := b.ConvertCross(msg)
		if err != nil {
			logger.ERROR("can't convert to cross: %s", err)
			return
		}
		crossIdStr, exist, err := b.kvSaver.Check([]string{chatroomId})
		if err != nil {
			logger.ERROR("can't check uin %s: %s", chatroomId, err)
			return
		}
		if exist {
			if err := b.UpdateCross(crossIdStr, cross); err == nil {
				return
			}
		}
		b.GatherCross(chatroomId, cross)
	} else {
		b.GreetNewFriend(msg.FromUserName)
	}
}

func (b *Bot) GreetNewFriend(username string) {
	req := []ContactRequest{
		ContactRequest{
			UserName: username,
		},
	}
	contacts, err := b.wc.GetContact(req)
	if err != nil {
		logger.ERROR("can't get contact for %s: %s", username, err)
		return
	}
	if len(contacts) == 0 {
		logger.ERROR("can't get contact for %s: no return", username)
		return
	}
	contact := contacts[0]
	// add user contact to exfe
	err = b.wc.SendMessage(contact.UserName, "要为您的微信群画张“活点地图”，把我加到群里就行啦。在群里打开过活点地图的人，能互相看到方位。")
	if err != nil {
		logger.ERROR("can't send to %s greet: %s", username, err)
		return
	}

	go func() {
		headerUrl := "https://wx.qq.com" + contact.HeadImgUrl
		resp, err := b.wc.request("GET", headerUrl, nil, nil)
		if err == nil {
			headerUrl, err = b.SaveHeader(contact.Uin, resp)
			if err != nil {
				logger.ERROR("save header failed: %s", err)
			}
		} else {
			logger.ERROR("can't get header %s: %s", headerUrl, err)
		}
		user, _, err := b.platform.GetUserByIdentity(contact.ToIdentity(headerUrl))
		if err != nil {
			return
		}
		logger.INFO("wechat_newuser", "uin", contact.Uin, "user", user.Id)
		password := fmt.Sprintf("%04d", rand.Intn(1e5))
		err = b.platform.SetPassword(user.Id, password)
		if err != nil {
			return
		}
		b.wc.SendMessage(contact.UserName, fmt.Sprintf("活点地图”是 水滴·X· 群组工具的功能之一。为了避免您的微信账号被他人误领，您的·X·默认密码为： %s。", password))
		if !user.Password {
			b.SetPassword(user.Id)
		}
	}()
}

func (b *Bot) SetPassword(userId int64) {
}

func (b *Bot) SaveHeader(uin uint64, resp *http.Response) (string, error) {
	defer resp.Body.Close()

	headerPath := fmt.Sprintf("/thirdpart/weichat/%d.jpg", uin)
	obj, err := b.bucket.CreateObject(headerPath, resp.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	obj.SetDate(time.Now())
	length, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return "", err
	}
	err = obj.SaveReader(resp.Body, length)
	if err != nil {
		return "", err
	}
	return obj.URL(), nil
}

func (b *Bot) ConvertCross(msg Message) (string, model.Cross, error) {
	if !strings.HasSuffix(msg.FromUserName, "@chatroom") {
		return "", model.Cross{}, fmt.Errorf("%s", "not join chat room")
	}
	chatroomReq := []ContactRequest{
		ContactRequest{
			UserName: msg.FromUserName,
		},
	}
	fmt.Println("get contact list")
	chatrooms, err := b.wc.GetContact(chatroomReq)
	if err != nil {
		return "", model.Cross{}, err
	}
	var chatroom Contact
	for _, c := range chatrooms {
		if c.UserName == msg.FromUserName {
			chatroom = c
			break
		}
	}
	if chatroom.UserName != msg.FromUserName {
		return "", model.Cross{}, fmt.Errorf("can't find chatroom %s", msg.FromUserName)
	}

	ret := model.Cross{}
	ret.Title = "Cross with "
	ret.Exfee.Invitations = make([]model.Invitation, len(chatroom.MemberList))
	var host *model.Identity
	for i, member := range chatroom.MemberList {
		if i < 3 && member.Uin != b.wc.baseRequest.Uin {
			ret.Title += member.NickName + ", "
		}
		var headerUrl string
		resp, err := b.wc.GetChatroomHeader(member.UserName, msg.FromUserName)
		if err == nil {
			headerUrl, err = b.SaveHeader(member.Uin, resp)
			if err != nil {
				logger.ERROR("can't save header: %s", err)
			}
		} else {
			logger.ERROR("can't get user %s header at %s: %s", member.UserName, msg.FromUserName, err)
		}
		ret.Exfee.Invitations[i].Identity = model.Identity{
			ExternalID:       fmt.Sprintf("%d", member.Uin),
			ExternalUsername: member.UserName,
			Provider:         "wechat",
			Nickname:         member.NickName,
			Avatar:           headerUrl,
			Locale:           "zh_cn",
			Timezone:         "Asia/Shanghai",
		}
		if member.Uin == chatroom.OwnerUin {
			ret.Exfee.Invitations[i].Host = true
			host = &ret.Exfee.Invitations[i].Identity
		}
	}
	ret.Title = ret.Title[:len(ret.Title)-2]
	ret.By = *host
	for i := range ret.Exfee.Invitations {
		ret.Exfee.Invitations[i].By = *host
		ret.Exfee.Invitations[i].UpdatedBy = *host
	}
	return chatroom.UserName, ret, nil
}

func (b *Bot) GatherCross(chatroomId string, cross model.Cross) {
	cross, err := b.platform.BotCrossGather(cross)
	if err != nil {
		logger.ERROR("can't gather cross: %s", err)
		return
	}
	err = b.kvSaver.Save([]string{chatroomId}, fmt.Sprintf("%d", cross.ID))
	if err != nil {
		logger.ERROR("can't save cross id: %s", err)
	}
	err = b.kvSaver.Save([]string{fmt.Sprintf("e%d@exfe", cross.Exfee.ID)}, chatroomId)
	if err != nil {
		logger.ERROR("can't save exfee id: %s", err)
	}
	logger.INFO("wechat_gather", chatroomId, "cross", cross.ID, "exfee", cross.Exfee.ID)
	smith, err := cross.Exfee.FindInvitedUser(model.Identity{
		ExternalID: fmt.Sprintf("%d", b.wc.baseRequest.Uin),
		Provider:   "wechat",
	})
	if err != nil {
		logger.ERROR("can't find Smith Exfer in cross %d: %s", cross.ID, err)
		return
	}
	u := fmt.Sprintf("%s/#!token=%s/routex/", b.config.SiteUrl, smith.Token)
	err = b.wc.SendMessage(chatroomId, u)
	logger.NOTICE("send %s to %s", u, chatroomId)
	if err != nil {
		logger.ERROR("can't send %s to %s", u, chatroomId)
		return
	}
}

func (b *Bot) UpdateCross(crossIdStr string, cross model.Cross) error {
	crossId, err := strconv.ParseInt(crossIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid cross id %s: %s", crossIdStr, err)
	}
	oldCross, err := b.platform.FindCross(crossId, nil)
	if err != nil {
		return fmt.Errorf("can't find cross by id %d: %s", crossId, err)
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
		if invitation.Identity.Provider != "wechat" || exfee[invitation.Identity.Id()] {
			continue
		}
		invitation.Response = model.Removed
		cross.Exfee.Invitations = append(cross.Exfee.Invitations, invitation)
	}
	err = b.platform.BotCrossUpdate("cross_id", crossIdStr, cross, host)
	if err != nil {
		logger.ERROR("can't update cross %s: %s", crossIdStr, err)
	}
	logger.INFO("wechat_update", "cross", crossIdStr)
	return nil
}
