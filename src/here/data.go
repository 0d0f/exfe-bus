package here

import (
	"strconv"
	"time"
)

type Identity struct {
	ExternalID       string `json:"external_id"`
	ExternalUsername string `json:"external_username"`
	Provider         string `json:"provider"`
}

type Card struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Avatar     string     `json:"avatar"`
	Bio        string     `json:"bio"`
	Identities []Identity `json:"identities"`
}

type Data struct {
	Token     string   `json:"token"`
	Latitude  string   `json:"latitude"`
	Longitude string   `json:"longitude"`
	Accuracy  string   `json:"accuracy"`
	Traits    []string `json:"traits"`
	Card      Card     `json:"card"`

	latitude  float64 `json:"-"`
	longitude float64 `json:"-"`
	accuracy  float64 `json:"-"`

	UpdatedAt time.Time `json:"-"`
}

func (d Data) HasGPS() bool {
	return d.Latitude != "" && d.Longitude != "" && d.Accuracy != ""
}

func (d *Data) Init() error {
	if d.HasGPS() {
		var err error
		d.latitude, err = strconv.ParseFloat(d.Latitude, 64)
		if err != nil {
			return err
		}
		d.longitude, err = strconv.ParseFloat(d.Longitude, 64)
		if err != nil {
			return err
		}
		d.accuracy, err = strconv.ParseFloat(d.Accuracy, 64)
		if err != nil {
			return err
		}
	} else {
		d.Latitude = ""
		d.Longitude = ""
		d.Accuracy = ""
	}
	return nil
}
