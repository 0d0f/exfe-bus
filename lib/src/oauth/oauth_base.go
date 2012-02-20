package oauth

import (
	"github.com/garyburd/go-oauth"
	"net/http"
	"net/url"
	"encoding/json"
	"io"
	"fmt"
	"strings"
//	"io/ioutil"
)

type OAuthBase struct {
	ClientToken oauth.Credentials

	client *http.Client
}

type OAuthRequest struct {
	OAuthBase

	RequestTokenUri string
	AuthorizationUri string
	AccessTokenUri string
}

type OAuthClient struct {
	OAuthBase
	AccessToken oauth.Credentials
	ApiBaseUri string
	Headers url.Values
}


func SetHttpClient(oauth *OAuthBase, client *http.Client) {
	oauth.client = client
}


func (o *OAuthBase) init() {
	if (o.client == nil) {
		SetHttpClient(o, http.DefaultClient)
	}
}


func CreateOAuthRequest(token, secret, requestTokenUri, authorzationUri, accessTokenUri string) *OAuthRequest {
	return &OAuthRequest{
		OAuthBase: OAuthBase{
			ClientToken: oauth.Credentials{
				Token: token,
				Secret: secret,
			},
		},
		RequestTokenUri: requestTokenUri,
		AuthorizationUri: authorzationUri,
		AccessTokenUri: accessTokenUri,
	}
}


func (o *OAuthRequest) getClient() (client oauth.Client) {
	client = oauth.Client{
		Credentials: o.ClientToken,
		TemporaryCredentialRequestURI: o.RequestTokenUri,
		ResourceOwnerAuthorizationURI: o.AuthorizationUri,
		TokenRequestURI: o.AccessTokenUri,
	}
	return
}

func (o *OAuthRequest) AuthorizeUrl(callback string, request_params url.Values, auth_params url.Values) (token *oauth.Credentials, url string, err error) {
	o.init()
	client := o.getClient()
	token, err = client.RequestTemporaryCredentials(o.client, callback, request_params)
	if (err != nil) {
		return
	}

	url = client.AuthorizationURL(token, auth_params)
	return
}

func (o *OAuthRequest) AccessClient(token *oauth.Credentials, verifier string, apiBaseUri string) (client *OAuthClient, err error) {
	o.init()
	c := o.getClient()
	t, _, err := c.RequestToken(o.client, token, verifier)
	if (err != nil) {
		return
	}

	client = &OAuthClient{
		OAuthBase: o.OAuthBase,
		AccessToken: oauth.Credentials{
			Token: t.Token,
			Secret: t.Secret,
		},
		ApiBaseUri: apiBaseUri,
	}

	return
}


func isExistKey(m map[string]interface{}, key string) (ok bool) {
	_, ok = m[key]
	return
}

func CreateWithJson(r io.Reader) (*OAuthClient, error) {
	decoder := json.NewDecoder(r)
	var t map[string]interface{}
	decoder.Decode(&t)

	if !isExistKey(t, "OAuthBase") { return nil, fmt.Errorf("Can't find OAuthBase") }
	if !isExistKey(t, "AccessToken") { return nil, fmt.Errorf("Can't find AccessToken") }
	if !isExistKey(t, "ApiBaseUri") { return nil, fmt.Errorf("Can't find ApiBaseUri") }
	headers, err := url.ParseQuery(t["Headers"].(string))
	if (err != nil) {
		return nil, err
	}

	return &OAuthClient{
		OAuthBase: OAuthBase{
			ClientToken: oauth.Credentials{
				Token: t["OAuthBase"].(map[string]interface{})["ClientToken"].(map[string]interface{})["Token"].(string),
				Secret: t["OAuthBase"].(map[string]interface{})["ClientToken"].(map[string]interface{})["Secret"].(string),
			},
		},
		AccessToken: oauth.Credentials{
			Token: t["AccessToken"].(map[string]interface{})["Token"].(string),
			Secret: t["AccessToken"].(map[string]interface{})["Secret"].(string),
		},
		ApiBaseUri: t["ApiBaseUri"].(string),
		Headers: headers,
	}, nil
}

func (o *OAuthClient) Dump(w io.Writer) error {
	encoder := json.NewEncoder(w)
	dump := make(map[string]interface{})
	dump["OAuthBase"] = o.OAuthBase
	dump["AccessToken"] = o.AccessToken
	dump["ApiBaseUri"] = o.ApiBaseUri
	dump["Headers"] = o.Headers.Encode()
	return encoder.Encode(dump)
}

func (o *OAuthClient) getClient() *oauth.Client {
	return &oauth.Client{
		Credentials: o.ClientToken,
	}
}

func (o *OAuthClient) GetRequest(method, path string, params url.Values) (*http.Request, error) {
	header := http.Header{}
	for k, _ := range o.Headers {
		header.Add(k, o.Headers.Get(k))
	}
	uri := strings.TrimRight(o.ApiBaseUri, "/") + "/" + strings.TrimLeft(path, "/")
	header.Add("Authorization", o.getClient().AuthorizationHeader(&o.AccessToken, method, uri, params))

	u, err := url.Parse(uri + "?" + params.Encode())
	if (err != nil) {
		return nil, err
	}
	return &http.Request{
		Method: strings.ToUpper(method),
		URL: u,
		Header: header,
	}, nil
}

func (o *OAuthClient) SendRequest(request *http.Request) (io.ReadCloser, error) {
	o.init()
	resp, err := o.client.Do(request)
	if (err != nil) {
		return nil, err
	}

	if resp.StatusCode / 100 != 2 {
		return resp.Body, fmt.Errorf(resp.Status)
	}

	return resp.Body, nil
}
