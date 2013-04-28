package model

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
)

type UpdateInfo struct {
	UpdatedAt string `json:"updated_at"`
}

type Cross struct {
	ID          uint64                   `json:"id,omitempty"`
	By          Identity                 `json:"by_identity,omitempty"`
	Title       string                   `json:"title,omitempty"`
	Description string                   `json:"description,omitempty"`
	Time        *CrossTime               `json:"time,omitempty"`
	Place       *Place                   `json:"place,omitempty"`
	Exfee       Exfee                    `json:"exfee,omitempty"`
	Updated     map[string]UpdateInfo    `json:"updated,omitempty"`
	Widgets     []map[string]interface{} `json:"widget"`
}

func (c Cross) Equal(other *Cross) bool {
	return c.ID == other.ID
}

func (c Cross) String() string {
	return fmt.Sprintf("Cross:%d", c.ID)
}

func (c Cross) Ics(config *Config, to Recipient) string {
	url := fmt.Sprintf("%s/v2/ics/crosses/%d?token=%s", config.SiteApi, c.ID, to.Token)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (c Cross) TitleBackground(config *Config) (string, error) {
	bgUrl := c.findBackground(config)
	if bgUrl == "" {
		return "", nil
	}
	bg, err := http.Get(bgUrl)
	if err != nil {
		return "", nil
	}
	defer bg.Body.Close()
	if bg.StatusCode != 200 {
		return "", nil
	}

	buf := bytes.NewBuffer(nil)
	err = MakeTitle(buf, bg.Body, config.Pin, 640, 150, 199, 60, 2, c.Place.Lat, c.Place.Lng, 80)
	if err != nil {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (c Cross) Timezone(to Recipient) string {
	if to.Timezone != "" {
		return to.Timezone
	}
	return c.Time.BeginAt.Timezone
}

func (c Cross) Link(to Recipient, config *Config) string {
	return fmt.Sprintf("%s/#!token=%s", config.SiteUrl, to.Token)
}

func (c Cross) PublicLink(to Recipient, config *Config) string {
	token := to.Token
	if len(token) > 5 {
		token = token[1:5]
	}
	return fmt.Sprintf("%s/#!%d/%s", config.SiteUrl, c.ID, token)
}

func (c Cross) findBackground(config *Config) string {
	for _, w := range c.Widgets {
		if t, ok := w["type"].(string); !ok || t != "Background" {
			continue
		}
		if img, ok := w["image"]; ok && img != "" {
			return fmt.Sprintf("%s/static/img/xbg/%s", config.SiteUrl, img)
		}
	}
	return fmt.Sprintf("%s/static/img/xbg/default.jpg", config.SiteUrl)
}

func MakeTitle(w io.Writer, bg io.Reader, pin io.Reader, width, height, offsetY int, columnX, columnWidth int, latitude, longitude string, mapWidth int) error {
	img, _, err := image.Decode(bg)
	if err != nil {
		return err
	}
	if y := img.Bounds().Dy() - offsetY; height > y {
		height = y
	}

	rect := image.Rect(0, 0, img.Bounds().Dx(), height*img.Bounds().Dx()/width)
	rect = img.Bounds().Sub(image.Pt(0, offsetY)).Intersect(rect)
	var out draw.Image = image.NewRGBA(rect)
	draw.Draw(out, rect, img, image.Pt(0, offsetY), draw.Src)
	img = resize.Resize(uint(width), uint(height), out, resize.Lanczos3)
	out, ok := img.(draw.Image)
	if !ok {
		out = image.NewRGBA(img.Bounds())
		draw.Draw(out, out.Bounds(), img, image.Pt(0, 0), draw.Src)
	}

	draw.Draw(out, out.Bounds(), image.NewUniform(color.RGBA{0, 0, 0, 51}), image.Pt(0, 0), draw.Over)
	rect = image.Rect(columnX, 0, columnX+columnWidth, height)
	draw.Draw(out, rect, image.NewUniform(color.RGBA{0x40, 0x40, 0x40, 0xff / 4}), image.Pt(0, 0), draw.Over)

	if latitude != "" && longitude != "" {
		pinImage, _, err := image.Decode(pin)
		if err != nil {
			return err
		}

		pinTopMargin := 10
		pinPoint := image.Pt(mapWidth/2, pinImage.Bounds().Dy()+pinTopMargin)
		img, rect, err = GetMap(latitude, longitude, mapWidth, height, pinPoint)
		if err != nil {
			return err
		}
		draw.Draw(out, image.Rect(width-mapWidth, 0, width, height), img, rect.Min, draw.Src)

		pinPoint.X, pinPoint.Y = width-mapWidth/2-pinImage.Bounds().Dx()/2, pinTopMargin
		draw.Draw(out, image.Rect(pinPoint.X, pinPoint.Y, width, height), pinImage, image.Pt(0, 0), draw.Over)
	}

	err = jpeg.Encode(w, out, &jpeg.Options{70})
	if err != nil {
		return err
	}
	return nil
}

func GetMap(latitude, longitude string, width, height int, pin image.Point) (image.Image, image.Rectangle, error) {
	if pin.X >= width || pin.Y >= height {
		return nil, image.Rectangle{}, fmt.Errorf("pin's point not in the image")
	}
	w, h, offsetX, offsetY := centerize(width, height, pin)
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?center=%s,%s&zoom=13&size=%dx%d&maptype=road&sensor=false",
		latitude, longitude, w, h)
	resp, err := http.Get(url)
	if err != nil {
		return nil, image.Rectangle{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, image.Rectangle{}, fmt.Errorf("%s", resp.Status)
	}
	ret, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, image.Rectangle{}, err
	}
	rect := ret.Bounds().Intersect(image.Rect(offsetX, offsetY, width+offsetX, height+offsetY))

	return ret, rect, nil
}

func centerize(width, height int, pin image.Point) (w int, h int, offsetX int, offsetY int) {
	w, h = pin.X*2, pin.Y*2
	if right := width - pin.X; right > pin.X {
		w = right * 2
		offsetX = right - pin.X
	}
	if bottom := height - pin.Y; bottom > pin.Y {
		h = bottom * 2
		offsetY = bottom - pin.Y
	}
	return
}

type CrossDigestRequest struct {
	To        Recipient `json:"to"`
	CrossId   int64     `json:"cross_id"`
	UpdatedAt string    `json:"updated_at"`
}

type CrossUpdate struct {
	To       Recipient `json:"to"`
	OldCross Cross     `json:"old_cross"`
	Cross    Cross     `json:"cross"`
	By       Identity  `json:"by"`
}

type CrossUpdates []CrossUpdate

func (u CrossUpdates) String() string {
	if len(u) == 0 {
		return "{updates:0}"
	}
	return fmt.Sprintf("{to:%s with:%s updates:%d}", u[0].To, u[0].Cross, len(u))
}

type CrossInvitation struct {
	To    Recipient `json:"to"`
	Cross Cross     `json:"cross"`
	By    Identity  `json:"by"`

	Config *Config `json:"-"`
}

func (a CrossInvitation) String() string {
	return fmt.Sprintf("{to:%s cross:%d}", a.To, a.Cross.ID)
}

func (a *CrossInvitation) Parse(config *Config) (err error) {
	a.Config = config
	return nil
}

func (a CrossInvitation) ToIn(invitations []Invitation) bool {
	for _, i := range invitations {
		if a.To.SameUser(&i.Identity) {
			return true
		}
	}
	return false
}

func (a CrossInvitation) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a CrossInvitation) PublicLink() string {
	return fmt.Sprintf("%s/#!%d/%s", a.Config.SiteUrl, a.Cross.ID, a.To.Token[1:5])
}

func (a CrossInvitation) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
}

func (a CrossInvitation) IsCreator() bool {
	return a.To.SameUser(&a.Cross.By)
}

func (a CrossInvitation) LongDescription() bool {
	if len(a.Cross.Description) > 200 {
		return true
	}
	return false
}

func (a CrossInvitation) ListInvitations() string {
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
