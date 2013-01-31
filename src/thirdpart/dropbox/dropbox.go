package dropbox

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-aws/s3"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"model"
	"net/http"
	"strings"
	"time"
)

type Dropbox struct {
	consumer *oauth.Consumer
	bucket   *s3.Bucket
}

func New(config *model.Config) (*Dropbox, error) {
	provider := oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.dropbox.com/1/oauth/request_token",
		AuthorizeTokenUrl: "https://www.dropbox.com/1/oauth/authorize",
		AccessTokenUrl:    "https://api.dropbox.com/1/oauth/access_token",
	}
	consumer := oauth.NewConsumer(config.Thirdpart.Dropbox.Key, config.Thirdpart.Dropbox.Secret, provider)
	aws := s3.New(config.AWS.S3.Domain, config.AWS.S3.Key, config.AWS.S3.Secret)
	bucket, err := aws.GetBucket(fmt.Sprintf("%s-3rdpart-dropbox", config.Thirdpart.Dropbox.BucketPrefix))
	if err != nil {
		return nil, err
	}
	return &Dropbox{
		consumer: consumer,
		bucket:   bucket,
	}, nil
}

func (d *Dropbox) Provider() string {
	return "dropbox"
}

func (d *Dropbox) Grab(to model.Recipient, albumID string) ([]model.Photo, error) {
	var data model.OAuthToken
	err := json.Unmarshal([]byte(to.AuthData), &data)
	if err != nil {
		return nil, fmt.Errorf("can't parse %s auth data(%s): %s", to, to.AuthData, err)
	}
	token := oauth.AccessToken{
		Token:  data.Token,
		Secret: data.Secret,
	}
	resp, err := d.consumer.Get(fmt.Sprintf("https://api.dropbox.com/1/metadata/dropbox%s", albumID), nil, &token)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}
	decoder := json.NewDecoder(resp.Body)
	var list folderList
	err = decoder.Decode(&list)
	if err != nil {
		return nil, err
	}
	ret := make([]model.Photo, 0)
	for _, c := range list.Contents {
		if !c.ThumbExists {
			continue
		}
		caption := c.Path[strings.LastIndex(c.Path, "/"):]
		modified, err := time.Parse(time.RFC1123Z, c.Modified)
		if err != nil {
			modified = time.Now()
		}
		photo := model.Photo{
			Caption: caption,
			By: model.Identity{
				ID: to.IdentityID,
			},
			CreatedAt:       modified.Format("2006-01-02 15:04:05"),
			UpdatedAt:       modified.Format("2006-01-02 15:04:05"),
			Provider:        "dropbox",
			ExternalAlbumID: albumID,
			ExternalID:      c.Rev,
		}

		thumb, big, err := d.savePic(c, to, &token)
		if err != nil {
			continue
		}
		photo.Images.Thumbnail.Url = thumb
		photo.Images.Thumbnail.Height = 480
		photo.Images.Thumbnail.Width = 640
		photo.Images.Fullsize.Url = big
		photo.Images.Thumbnail.Height = 768
		photo.Images.Thumbnail.Width = 1024
	}
	return ret, nil
}

func (d *Dropbox) savePic(c content, to model.Recipient, token *oauth.AccessToken) (string, string, error) {
	thumb, err := d.consumer.Get(fmt.Sprintf("https://api-content.dropbox.com/1/thumbnails/dropbox%s", c.Path), map[string]string{"size": "l"}, token)
	if err != nil {
		return "", "", err
	}
	defer thumb.Body.Close()
	if thumb.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(thumb.Body)
		if err != nil {
			return "", "", err
		}
		return "", "", fmt.Errorf("%s: %s", thumb.Status, body)
	}

	big, err := d.consumer.Get(fmt.Sprintf("https://api-content.dropbox.com/1/thumbnails/dropbox%s", c.Path), map[string]string{"size": "xl"}, token)
	if err != nil {
		return "", "", err
	}
	defer big.Body.Close()
	if big.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(big.Body)
		if err != nil {
			return "", "", err
		}
		return "", "", fmt.Errorf("%s: %s", big.Status, body)
	}

	extIndex := strings.LastIndex(c.Path, ".")
	pathIndex := strings.LastIndex(c.Path, "/")
	thumbName := c.Path
	if extIndex < pathIndex {
		thumbName += "-thumb"
	} else {
		thumbName = fmt.Sprintf("%s-thumb%s", c.Path[:extIndex], c.Path[extIndex:])
	}
	thumbObj, err := d.bucket.CreateObject(fmt.Sprintf("i%d/%s", to.IdentityID, thumbName), c.MimeType)
	if err != nil {
		return "", "", err
	}

	bigObj, err := d.bucket.CreateObject(fmt.Sprintf("i%d/%s", to.IdentityID, c.Path), c.MimeType)
	if err != nil {
		return "", "", err
	}

	thumbObj.Save(thumb.Body)
	bigObj.Save(big.Body)

	return thumbObj.URL(), bigObj.URL(), nil
}

type content struct {
	Rev         string `json:"rev"`
	ThumbExists bool   `json:"thumb_exists"`
	Bytes       int    `json:"bytes"`
	Modified    string `json:"modified"`
	Path        string `json:"path"`
	IsDir       bool   `json:"is_dir"`
	Root        string `json:"root"`
	MimeType    string `json:"mime_type"`
}

type folderList struct {
	Contents []content `json:"contents"`
}
