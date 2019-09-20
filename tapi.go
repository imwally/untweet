package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const apiURL string = "https://api.twitter.com/1.1/"

type User struct {
	Id          int    `json:"id"`
	ScreenName  string `json:"screen_name"`
	DisplayName string `json:"name"`
	Following   bool   `json:"following"`
}

type Tweet struct {
	CreatedAt string `json:"created_at"`
	Id        int    `json:"id"`
	User      `json:"user"`
	URL       string
}

type TwitterAPI struct {
	KeyConsumer       string
	KeySecret         string
	AccessToken       string
	AccessTokenSecret string
	BearerToken       string
}

type TwitterAPIRequest struct {
	Parameters map[string]string
	Headers    http.Header
	EndPoint   string
	Method     string
	Body       string
	Auth       string
}

func (ta *TwitterAPI) GenerateOauthSignature(tar *TwitterAPIRequest, nonce string, ts string) string {
	params := make(map[string]string)
	for k, v := range tar.Parameters {
		params[k] = v
	}

	params["oauth_consumer_key"] = url.QueryEscape(ta.KeyConsumer)
	params["oauth_nonce"] = url.QueryEscape(nonce)
	params["oauth_signature_method"] = url.QueryEscape("HMAC-SHA1")
	params["oauth_timestamp"] = url.QueryEscape(ts)
	params["oauth_token"] = url.QueryEscape(ta.AccessToken)
	params["oauth_version"] = url.QueryEscape("1.0")

	baseURL, _ := url.Parse(tar.EndPoint)
	baseURL.RawQuery = ""

	baseString := GenerateParameterString(params, true)
	baseString = tar.Method + "&" + url.QueryEscape(baseURL.String()) + "&" + url.QueryEscape(baseString)

	key := url.QueryEscape(ta.KeySecret) + "&" + url.QueryEscape(ta.AccessTokenSecret)
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(baseString))

	sig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return sig
}

func (ta *TwitterAPI) Request(tar *TwitterAPIRequest) ([]byte, error) {
	if tar == nil {
		return nil, errors.New("error: unsupported resource")
	}

	client := &http.Client{}
	req, err := http.NewRequest(tar.Method, tar.EndPoint, strings.NewReader(tar.Body))
	if err != nil {
		return nil, err
	}

	if tar.Headers != nil {
		req.Header = tar.Headers
	}

	if tar.Auth == "oauth" {
		nonce, err := GenerateNonce()
		if err != nil {
			return nil, err
		}

		ts := strconv.FormatInt(time.Now().Unix(), 10)
		sig := ta.GenerateOauthSignature(tar, nonce, ts)

		header := "OAuth "
		header += fmt.Sprintf("%s=\"%s\", ", "oauth_consumer_key", ta.KeyConsumer)
		header += fmt.Sprintf("%s=\"%s\", ", "oauth_nonce", nonce)
		header += fmt.Sprintf("%s=\"%s\", ", "oauth_signature", url.QueryEscape(sig))
		header += fmt.Sprintf("%s=\"%s\", ", "oauth_signature_method", "HMAC-SHA1")
		header += fmt.Sprintf("%s=\"%s\", ", "oauth_timestamp", ts)
		header += fmt.Sprintf("%s=\"%s\", ", "oauth_token", ta.AccessToken)
		header += fmt.Sprintf("%s=\"%s\"", "oauth_version", "1.0")
		req.Header.Add("authorization", header)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Header.Get("X-Rate-Limit-Remaining") == "0" {
		resetHeader := resp.Header.Get("X-Rate-Limit-Reset")
		unixTime, err := strconv.ParseInt(resetHeader, 0, 64)
		if err != nil {
			log.Println(err)
		}

		resetTime := time.Unix(unixTime, 0)
		if err != nil {
			log.Println(err)
		}
		until := time.Until(resetTime)

		log.Println("hit rate limit, waiting until", resetTime, "to proceed")
		time.Sleep(until)

		return ta.Request(tar)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (ta *TwitterAPI) GetBatchedLikes(sn string, maxId int) ([]Tweet, error) {
	params := make(map[string]string)
	params["include_entities"] = "true"
	params["screen_name"] = sn
	params["count"] = "200"

	if maxId > 0 {
		params["max_id"] = strconv.Itoa(maxId)
	}

	req := NewRequest("favorites/list", params)
	resp, err := ta.Request(req)
	if err != nil {
		return nil, err
	}

	var tweets []Tweet
	json.Unmarshal(resp, &tweets)

	return tweets, nil
}

func (ta *TwitterAPI) GetLikes(sn string) ([]Tweet, error) {
	var likes []Tweet
	for max, next := 0, 1; next > 0; {
		batch, err := ta.GetBatchedLikes(sn, max)
		if err != nil {
			return nil, err
		}

		for _, like := range batch {
			like.URL = fmt.Sprintf("https://twitter.com/%s/status/%d", like.ScreenName, like.Id)
			likes = append(likes, like)
		}

		log.Println("gathered", len(likes), "likes")

		blen := len(batch)
		if blen > 0 {
			max = batch[blen-1].Id - 1
		}

		next = blen
	}

	return likes, nil
}

func (ta *TwitterAPI) DestroyLike(id int) error {
	params := make(map[string]string)
	params["id"] = strconv.Itoa(id)

	log.Printf("destroying tweet %d\n", id)
	req := NewRequest("favorites/destroy", params)

	_, err := ta.Request(req)
	if err != nil {
		return err
	}

	return nil
}

func NewRequest(resource string, parameters map[string]string) *TwitterAPIRequest {
	switch resource {
	case "favorites/list":
		return &TwitterAPIRequest{
			Parameters: parameters,
			Method:     http.MethodGet,
			EndPoint:   apiURL + resource + ".json?" + GenerateParameterString(parameters, false),
			Auth:       "oauth",
		}
	case "favorites/destroy":
		return &TwitterAPIRequest{
			Parameters: parameters,
			Method:     http.MethodPost,
			EndPoint:   apiURL + resource + ".json?" + GenerateParameterString(parameters, false),
			Auth:       "oauth",
		}
	}

	return nil
}
