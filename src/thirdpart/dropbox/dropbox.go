package dropbox

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-aws/s3"
	"github.com/googollee/go-logger"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Dropbox struct {
	consumer *oauth.Consumer
	bucket   *s3.Bucket
	log      *logger.SubLogger
}

func New(config *model.Config) (*Dropbox, error) {
	provider := oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.dropbox.com/1/oauth/request_token",
		AuthorizeTokenUrl: "https://www.dropbox.com/1/oauth/authorize",
		AccessTokenUrl:    "https://api.dropbox.com/1/oauth/access_token",
	}
	consumer := oauth.NewConsumer(config.Thirdpart.Dropbox.Key, config.Thirdpart.Dropbox.Secret, provider)
	aws := s3.New(config.AWS.S3.Domain, config.AWS.S3.Key, config.AWS.S3.Secret)
	aws.SetACL(s3.ACLPublicRead)
	aws.SetLocationConstraint(s3.LC_AP_SINGAPORE)
	bucket, err := aws.GetBucket(fmt.Sprintf("%s-3rdpart-photos", config.AWS.S3.BucketPrefix))
	if err != nil {
		return nil, err
	}
	return &Dropbox{
		consumer: consumer,
		bucket:   bucket,
		log:      config.Log.SubPrefix("dropbox"),
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
	path := escapePath(albumID)
	path = fmt.Sprintf("https://api.dropbox.com/1/metadata/dropbox%s", path)
	resp, err := d.consumer.Get(path, nil, &token)
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
	fmt.Println("get meta")
	ret := make([]model.Photo, 0)
	for _, c := range list.Contents {
		fmt.Println("process", c.Path)
		if !c.ThumbExists {
			d.log.Info("%s %s is not picture.", to, c.Path)
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
			d.log.Info("%s %s can't save: %s", to, c.Path, err)
			continue
		}
		photo.Images.Preview.Url = thumb
		photo.Images.Preview.Height = 480
		photo.Images.Preview.Width = 640
		photo.Images.Fullsize.Url = big
		photo.Images.Fullsize.Height = 768
		photo.Images.Fullsize.Width = 1024
		ret = append(ret, photo)
	}
	return ret, nil
}

func (d *Dropbox) savePic(c content, to model.Recipient, token *oauth.AccessToken) (string, string, error) {
	path := fmt.Sprintf("https://api-content.dropbox.com/1/thumbnails/dropbox%s", escapePath(c.Path))
	thumbPath := fmt.Sprintf("/dropbox/i%d%s", to.IdentityID, getThumbName(c.Path))
	bigPath := fmt.Sprintf("/dropbox/i%d%s", to.IdentityID, c.Path)
	thumb, err := d.saveFile(path, "l", thumbPath, c.MimeType, token)
	if err != nil {
		return "", "", err
	}
	big, err := d.saveFile(path, "xl", bigPath, c.MimeType, token)
	if err != nil {
		return "", "", err
	}

	return thumb, big, nil
}

func (d *Dropbox) saveFile(from, size, to, mime string, token *oauth.AccessToken) (string, error) {
	resp, err := d.consumer.Get(from, map[string]string{"size": size}, token)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("%s: %s", resp.Status, body)
	}
	length, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return "", err
	}

	object, err := d.bucket.CreateObject(to, mime)
	if err != nil {
		return "", err
	}
	object.SetDate(time.Now())
	err = object.SaveReader(resp.Body, int64(length))
	if err != nil {
		return "", err
	}
	return object.URL(), nil
}

func getThumbName(path string) string {
	extIndex := strings.LastIndex(path, ".")
	pathIndex := strings.LastIndex(path, "/")
	thumbName := path
	if extIndex < pathIndex {
		return thumbName + "-thumb"
	}
	return fmt.Sprintf("%s-thumb%s", path[:extIndex], path[extIndex:])
}

func escapePath(path string) string {
	if path[0] == '/' {
		path = "/dropbox" + path
	} else {
		path = "/dropbox/" + path
	}
	path = url.QueryEscape(path)
	path = strings.Replace(path, "%2F", "/", -1)
	path = strings.Replace(path, "+", "%20", -1)
	return path
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
