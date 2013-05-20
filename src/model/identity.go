package model

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
)

type Identity struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	UserID   int64  `json:"connected_user_id,omitempty"`
	Avatar   string `json:"avatar_filename,omitempty"`

	Provider         string `json:"provider,omitempty"`
	ExternalID       string `json:"external_id,omitempty"`
	ExternalUsername string `json:"external_username,omitempty"`
	OAuthToken       string `json:"oauth_token,omitempty"`
}

func (i Identity) GetAvatar(x, y int) string {
	resp, err := http.Get(i.Avatar)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	resized, err := innerResize(resp.Body, x, y)
	if err != nil {
		return ""
	}

	buf := bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, resized, &jpeg.Options{70})
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (i Identity) Equal(other Identity) bool {
	if i.ID == other.ID {
		return true
	}
	if i.Provider == other.Provider {
		if i.ExternalID == other.ExternalID {
			return true
		}
		if i.ExternalUsername == other.ExternalUsername {
			return true
		}
	}
	return false
}

func (i Identity) SameUser(other Identity) bool {
	if i.Equal(other) {
		return true
	}
	return i.UserID == other.UserID
}

func (i Identity) String() string {
	return fmt.Sprintf("Identity:(i%d/u%d)", i.ID, i.UserID)
}

func (i Identity) ScreenId() string {
	switch i.Provider {
	case "email":
		return i.ExternalUsername
	case "phone":
		return i.ExternalUsername
	case "twitter":
		return "@" + i.ExternalUsername
	}
	return i.ExternalUsername + "@" + i.Provider
}

type RsvpType string

const (
	RsvpNoresponse   RsvpType = "NORESPONSE"
	RsvpAccepted              = "ACCEPTED"
	RsvpInterested            = "INTERESTED"
	RsvpDeclined              = "DECLINED"
	RsvpRemoved               = "REMOVED"
	RsvpNotification          = "NOTIFICATION"
)

type Invitation struct {
	ID         uint64   `json:"id,omitempty"`
	Host       bool     `json:"host,omitempty"`
	Mates      uint64   `json:"mates,omitempty"`
	Identity   Identity `json:"identity,omitempty"`
	RsvpStatus RsvpType `json:"rsvp_status,omitempty"`
	By         Identity `json:"by_identity,omitempty"`
	UpdatedBy  Identity `json:"updated_by,omitempty"`
	Via        string   `json:"via,omitempty"`
}

func (i *Invitation) String() string {
	return i.Identity.Name
}

func (i Invitation) IsAccepted() bool {
	return i.RsvpStatus == RsvpAccepted
}

func (i Invitation) IsDeclined() bool {
	return i.RsvpStatus == RsvpDeclined
}

func (i Invitation) IsPending() bool {
	return !i.IsAccepted() && !i.IsDeclined()
}

func (i Invitation) IsUpdatedBy(userId int64) bool {
	return i.UpdatedBy.UserID == userId
}

func innerResize(r io.Reader, x int, y int) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	tryX := x
	tryY := tryX * img.Bounds().Dy() / img.Bounds().Dx()
	var offset image.Point
	if tryY < y {
		tryY = y
		tryX = tryY * img.Bounds().Dx() / img.Bounds().Dy()
		offset = image.Pt((tryX-x)/2, 0)
	} else {
		offset = image.Pt(0, (tryY-y)/2)
	}

	img = resize.Resize(uint(tryX), uint(tryY), img, resize.Lanczos3)
	ret := image.NewRGBA(image.Rect(0, 0, x, y))
	draw.Draw(ret, ret.Bounds(), img, offset, draw.Src)

	return ret, nil
}

type OAuthToken struct {
	Token  string `json:"oauth_token"`
	Secret string `json:"oauth_token_secret"`
}
