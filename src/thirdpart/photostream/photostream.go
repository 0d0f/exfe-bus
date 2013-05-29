package photostream

import (
	"broker"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-aws/s3"
	"logger"
	"model"
	"strconv"
	"time"
)

type Derivative struct {
	Height   string `json:"height"`
	Width    string `json:"width"`
	Checksum string `json:"checksum"`
}

func (d Derivative) URL(urls UrlList) (string, error) {
	urlMeta, ok := urls.Items[d.Checksum]
	if !ok {
		return "", fmt.Errorf("can't find checksum: %s", d.Checksum)
	}
	location := urls.Locations[urlMeta.UrlLocation]
	return fmt.Sprintf("%s://%s%s", location.Scheme, location.Hosts[0], urlMeta.UrlPath), nil
}

type Derivatives map[string]Derivative

func (d Derivatives) GetMinAndMax() (Derivative, Derivative, error) {
	var big, small int
	for k := range d {
		id, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		if big == 0 {
			big, small = id, id
			continue
		}
		if id > big {
			big = id
		}
		if id < small {
			small = id
		}
	}
	if big == 0 || small == 0 {
		return Derivative{}, Derivative{}, fmt.Errorf("can't find derivatie")
	}
	return d[fmt.Sprintf("%d", small)], d[fmt.Sprintf("%d", big)], nil
}

type PhotoMeta struct {
	PhotoGuid   string      `json:"photoGuid"`
	DateCreated string      `json:"dateCreated"`
	Caption     string      `json:"caption"`
	Derivatives Derivatives `json:"derivatives"`
}

type StreamingList struct {
	Photos []PhotoMeta `json:"photos"`
}

type UrlRequest struct {
	PhotoGuids []string `json:"photoGuids"`
}

type UrlMeta struct {
	UrlExpiry   string `json:"url_expiry"`
	UrlPath     string `json:"url_path"`
	UrlLocation string `json:"url_location"`
}

type LocationMeta struct {
	Scheme string   `json:"scheme"`
	Hosts  []string `json:"hosts"`
}

type UrlList struct {
	Locations map[string]LocationMeta `json:"locations"`
	Items     map[string]UrlMeta      `json:"items"`
}

type Photostream struct {
	domain string
	bucket *s3.Bucket
}

func New(config *model.Config) (*Photostream, error) {
	aws := s3.New(config.AWS.S3.Domain, config.AWS.S3.Key, config.AWS.S3.Secret)
	aws.SetACL(s3.ACLPublicRead)
	aws.SetLocationConstraint(s3.LC_AP_SINGAPORE)
	bucket, err := aws.GetBucket(fmt.Sprintf("%s-3rdpart-photos", config.AWS.S3.BucketPrefix))
	if err != nil {
		return nil, err
	}
	return &Photostream{
		domain: config.Thirdpart.Photostream.Domain,
		bucket: bucket,
	}, nil
}

func (p *Photostream) Provider() string {
	return "photostream"
}

func (p *Photostream) Grab(to model.Recipient, albumID string) ([]model.Photo, error) {
	list, err := p.getList(albumID)
	if err != nil {
		return nil, fmt.Errorf("get streaming failed: %s", err)
	}
	guids := make([]string, len(list.Photos))
	for i, photo := range list.Photos {
		guids[i] = photo.PhotoGuid
	}
	urls, err := p.getUrls(albumID, guids)
	if err != nil {
		return nil, fmt.Errorf("get urls failed: %s", err)
	}

	ret := make([]model.Photo, 0)
	for _, photo := range list.Photos {
		preview, fullsize, err := photo.Derivatives.GetMinAndMax()
		if err != nil {
			logger.ERROR("can't find derivative of %s", photo.PhotoGuid)
			continue
		}
		t, err := time.Parse(photo.DateCreated, "2006-01-02T15:04:05Z")
		if err != nil {
			t = time.Now()
		}
		url, err := preview.URL(urls)
		if err != nil {
			continue
		}
		resp, err := broker.Http("GET", url, "application/json", nil)
		reader, err := broker.HttpResponse(resp, err)
		if err != nil {
			logger.ERROR("can't grab preview of %s from %s: %s", photo.PhotoGuid, url, err)
			continue
		}
		defer reader.Close()
		length, err := strconv.Atoi(resp.Header.Get("Content-Length"))
		if err != nil {
			logger.ERROR("can't parse %s length(%s): %s", photo.PhotoGuid, resp.Header.Get("Content-Length"), err)
			continue
		}

		object, err := p.bucket.CreateObject(fmt.Sprintf("/i%d/photostream/%s.jpg", to.IdentityID, photo.PhotoGuid), "image/jpeg")
		if err != nil {
			logger.ERROR("can't save %s to s3: %s", photo.PhotoGuid, err)
			continue
		}
		object.SetDate(time.Now())
		err = object.SaveReader(reader, int64(length))
		if err != nil {
			logger.ERROR("save %s to s3 failed: %s", photo.PhotoGuid, err)
			continue
		}

		photo_ := model.Photo{
			Caption: photo.Caption,
			By: model.Identity{
				ID: to.IdentityID,
			},
			CreatedAt:       t.Format("2006-01-02 15:04:05"),
			UpdatedAt:       t.Format("2006-01-02 15:04:05"),
			Provider:        p.Provider(),
			ExternalAlbumID: albumID,
			ExternalID:      photo.PhotoGuid,
		}
		photo_.Images.Fullsize.Height, _ = strconv.Atoi(fullsize.Height)
		photo_.Images.Fullsize.Width, _ = strconv.Atoi(fullsize.Width)
		photo_.Images.Fullsize.Url = fmt.Sprintf("photostream://%s/%s/%s", albumID, photo.PhotoGuid, fullsize.Checksum)
		photo_.Images.Preview.Height, _ = strconv.Atoi(preview.Height)
		photo_.Images.Preview.Width, _ = strconv.Atoi(preview.Width)
		photo_.Images.Preview.Url = object.URL()

		ret = append(ret, photo_)
	}
	return ret, nil
}

func (p *Photostream) Get(to model.Recipient, pictureIDs []string) ([]string, error) {
	return nil, fmt.Errorf("not support photostream.")
}

func (p *Photostream) getList(albumID string) (StreamingList, error) {
	url := fmt.Sprintf("https://%s/%s/sharedstreams/webstream", p.domain, albumID)
	buf := bytes.NewBufferString(`{"streamCtag":null}`)
	var list StreamingList
	resp, err := broker.HttpResponse(broker.Http("POST", url, "application/json", buf.Bytes()))
	if err != nil {
		return list, err
	}
	defer resp.Close()

	decoder := json.NewDecoder(resp)
	err = decoder.Decode(&list)
	return list, nil
}

func (p *Photostream) getUrls(albumID string, guids []string) (UrlList, error) {
	req := UrlRequest{
		PhotoGuids: guids,
	}
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(req)
	var list UrlList
	url := fmt.Sprintf("https://%s/%s/sharedstreams/webasseturls", p.domain, albumID)
	resp, err := broker.HttpResponse(broker.Http("POST", url, "application/json", buf.Bytes()))
	if err != nil {
		return list, err
	}
	defer resp.Close()

	decoder := json.NewDecoder(resp)
	err = decoder.Decode(&list)
	return list, nil
}
