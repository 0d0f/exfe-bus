package main

type MessageType int

const (
	JoinMessage MessageType = 100000
	ChatMessage             = 1
)

type BaseRequest struct {
	Uin      uint64
	Sid      string
	Skey     string
	DeviceID string
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
	Count       int              `json:"Count,omitempty"`
	List        []ContactRequest `json:"List,omitempty"`
	SyncKey     *SyncKey         `json:"SyncKey,omitempty"`
	Msg         *Message         `json:"Msg,omitempty"`
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
