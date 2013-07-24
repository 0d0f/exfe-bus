package main

import (
	"fmt"
	"model"
)

type MessageType int

const (
	JoinMessage   MessageType = 10000
	FriendRequest             = 37
	SystemMessage             = 51
	ChatMessage               = 1
)

type BaseRequest struct {
	Uin      uint64 `json:",omitempty"`
	Sid      string `json:",omitempty"`
	Skey     string `json:",omitempty"`
	DeviceID string `json:",omitempty"`
}

type SyncKey struct {
	Count int
	List  []map[string]int
}

type ContactRequest struct {
	UserName   string
	ChatRoomId uint64
}

type Request struct {
	BaseRequest BaseRequest
	Count       int              `json:",omitempty"`
	List        []ContactRequest `json:",omitempty"`
	SyncKey     *SyncKey         `json:",omitempty"`
	Msg         *Message         `json:",omitempty"`
}

type VerifyUser struct {
	Value            string
	VerifyUserTicket string
}

type VerifyRequest struct {
	BaseRequest        BaseRequest
	Opcode             int
	SceneList          []int
	SceneListCount     int
	VerifyContent      string
	VerifyUserList     []VerifyUser
	VerifyUserListSize int
}

type Member struct {
	AttrStatus      int
	DisplayName     string
	MemberStatus    int
	NickName        string
	PYInitial       string
	PYQuanPin       string
	RemarkPYInitial string
	RemarkPYQuanPin string
	Uin             uint64
	UserName        string
}

type Contact struct {
	Alias            string
	AppAccountFlag   int
	AttrStatus       int
	City             string
	ContactFlag      int
	HeadImgUrl       string
	HideInputBarFlag int
	MemberCount      int
	MemberList       []Member
	NickName         string
	OwnerUin         uint64
	PYInitial        string
	PYQuanPin        string
	Province         string
	RemarkName       string
	RemarkPYInitial  string
	RemarkPYQuanPin  string
	Sex              int
	Signature        string
	SnsFlag          int
	StarFriend       int
	Statues          int
	Uin              uint64
	UniFriend        int
	UserName         string
	VerifyFlag       int
}

func (c Contact) ToIdentity(headerUrl string) model.Identity {
	return model.Identity{
		ExternalID:       fmt.Sprintf("%d", c.Uin),
		ExternalUsername: c.UserName,
		Provider:         "wechat",
		Name:             c.NickName,
		Avatar:           headerUrl,
		Locale:           "zh_cn",
		Timezone:         "Asia/Shanghai",
	}
}

type Message struct {
	AppInfo struct {
		AddID string `json:",omitempty"`
		Type  int    `json:",omitempty"`
	} `json:",omitempty"`
	AppMsgType           int         `json:",omitempty"`
	Content              string      `json:",omitempty"`
	CreateTime           int64       `json:",omitempty"`
	ClientMsgId          int64       `json:",omitempty"`
	FileName             string      `json:",omitempty"`
	FileSize             string      `json:",omitempty"`
	ForwardFlag          int         `json:",omitempty"`
	FromUserName         string      `json:",omitempty"`
	ImgStatus            int         `json:",omitempty"`
	LocalID              int64       `json:",omitempty"`
	MediaId              string      `json:",omitempty"`
	MsgId                int64       `json:",omitempty"`
	MsgType              MessageType `json:",omitempty"`
	PlayLength           int         `json:",omitempty"`
	Status               int         `json:",omitempty"`
	StatusNotifyCode     int         `json:",omitempty"`
	StatusNotifyUserName string      `json:",omitempty"`
	ToUserName           string      `json:",omitempty"`
	Type                 int         `json:",omitempty"`
	Url                  string      `json:",omitempty"`
	VoiceLength          int         `json:",omitempty"`
	RecommendInfo        struct {
		UserName   string `json:",omitempty"`
		Ticket     string `json:",omitempty"`
		VerifyFlag int    `json:",omitempty"`
	} `json:",omitempty"`
}

type Response struct {
	BaseResponse struct {
		ErrMsg string
		Ret    int
	}
	User Contact

	AddMsgCount  int
	AddMsgList   []Message
	ContinueFlag int

	Count       int
	ContactList []Contact

	Skey    string
	SyncKey SyncKey
}
