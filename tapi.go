package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const API_URL string = "https://api.twitter.com/1.1/"

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Following bool   `json:"following"`
}

type Tweet struct {
	CreatedAt string `json:"created_at"`
	Id        int    `json:"id"`
	User      `json:"user"`
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

func GenerateOauthSignature(ta *TwitterAPI, tar *TwitterAPIRequest, nonce string, ts string) string {
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

func (ta *TwitterAPI) getLikes(sn string, count int, max int) ([]Tweet, error) {
	params := make(map[string]string)
	params["screen_name"] = sn
	params["count"] = strconv.Itoa(count)

	if max > 0 {
		params["max_id"] = strconv.Itoa(max)
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
	likes, err := ta.getLikes(sn, 30, 0)
	if err != nil {
		return nil, err
	}

	fmt.Println("Got", len(likes), "likes so far ...")
	for next := 1; next > 0; {
		max := likes[len(likes)-1].Id
		batch, err := ta.getLikes(sn, 30, max-1)
		if err != nil {
			return nil, err
		}

		for _, like := range batch {
			likes = append(likes, like)
		}

		next = len(batch)
		fmt.Println("Got", len(likes), "likes so far ...")
	}

	return likes, nil
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
		sig := GenerateOauthSignature(ta, tar, nonce, ts)

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

		reset := resp.Header.Get("X-Rate-Limit-Reset")
		unixTime, _ := strconv.ParseInt(reset, 0, 64)
		resetTime := time.Unix(unixTime, 0)
		if err != nil {
			fmt.Println(err)
		}
		until := time.Until(resetTime)

		fmt.Println("rate limit hit, waiting", until, "...")
		time.Sleep(until)

		return ta.Request(tar)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func NewRequest(resource string, parameters map[string]string) *TwitterAPIRequest {
	switch resource {
	case "favorites/list":
		return &TwitterAPIRequest{
			Parameters: parameters,
			Method:     http.MethodGet,
			EndPoint:   API_URL + resource + ".json?" + GenerateParameterString(parameters, false),
			Auth:       "oauth",
		}
	case "favorites/destroy":
		return &TwitterAPIRequest{
			Parameters: parameters,
			Method:     http.MethodPost,
			EndPoint:   API_URL + resource + ".json?" + GenerateParameterString(parameters, false),
			Auth:       "oauth",
		}
	}

	return nil
}
