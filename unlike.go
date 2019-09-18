package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const API_URL string = "https://api.twitter.com/1.1/"
const API_URL_TOKEN string = "https://api.twitter.com/oauth2/token"

type TokenResponse struct {
	Type  string `json:"token_type"`
	Token string `json:"access_token"`
}

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Following int    `json:"following"`
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

func PrintHeaders(resp *http.Response) {
	for headerKey, headerValue := range resp.Header {
		fmt.Printf("%s: %s\n", headerKey, headerValue)
	}
}

func PrintBody(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(body))
}

func UnmarshalToken(b []byte) string {
	var tr TokenResponse
	json.Unmarshal(b, &tr)

	return tr.Token
}

func GenerateParameterString(parameters map[string]string, sorted bool) string {
	params := ""

	if sorted {
		keys := make([]string, len(params))
		for k := range parameters {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			params += fmt.Sprintf("%s=%s&", k, parameters[k])
		}
	} else {
		for k, v := range parameters {
			params += fmt.Sprintf("%s=%s&", k, v)
		}
	}

	return params[:len(params)-1]
}

func GenerateNonce() (string, error) {
	n := make([]byte, 32)
	_, err := rand.Read(n)
	if err != nil {
		return "", err
	}
	nonce := base64.StdEncoding.EncodeToString(n)

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}

	nonce = reg.ReplaceAllString(nonce, "")

	return nonce, nil
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

func (ta *TwitterAPI) Request(tar *TwitterAPIRequest) ([]byte, error) {
	if tar == nil {
		return nil, errors.New("error: invalid resource")
	}

	client := &http.Client{}
	req, err := http.NewRequest(tar.Method, tar.EndPoint, strings.NewReader(tar.Body))
	if err != nil {
		return nil, err
	}

	if tar.Headers != nil {
		req.Header = tar.Headers
	}

	if tar.Auth == "basic" {
		req.SetBasicAuth(ta.KeyConsumer, ta.KeySecret)
	}

	if tar.Auth == "application" {
		req.Header.Add("Authorization", "Bearer "+ta.BearerToken)
	}

	if tar.Auth == "oauth" {
		nonce, err := GenerateNonce()
		if err != nil {
			log.Println(err)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func NewRequest(resource string, parameters map[string]string) *TwitterAPIRequest {
	switch resource {
	case "oauth2/token":
		tarHeaders := http.Header{}
		tarHeaders.Add("Content-Type", "application/x-www-form-urlencoded")
		return &TwitterAPIRequest{
			Method:   http.MethodPost,
			EndPoint: API_URL_TOKEN,
			Body:     "grant_type=client_credentials",
			Headers:  tarHeaders,
			Auth:     "basic",
		}
	case "favorites/list":
		return &TwitterAPIRequest{
			Method:   http.MethodGet,
			EndPoint: API_URL + resource + ".json?" + GenerateParameterString(parameters, false),
			Auth:     "application",
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

var KeyConsumer string
var KeySecret string
var AccessToken string
var AccessTokenSecret string

func init() {
	// Setup flags
	flag.StringVar(&KeyConsumer, "consumer", "", "Twitter API Consumer Key")
	flag.StringVar(&KeySecret, "secret", "", "Twitter API Secret Key")
	flag.StringVar(&AccessToken, "accesstoken", "", "Twitter API Access Token")
	flag.StringVar(&AccessTokenSecret, "accesstokensecret", "", "Twitter API Access Token Secret")
	flag.Parse()
}

func main() {
	if KeySecret == "" {
		fmt.Println("error: no secret key set")
		os.Exit(2)
	}

	if KeyConsumer == "" {
		fmt.Println("error: no consumer key set")
		os.Exit(2)
	}

	if AccessToken == "" {
		fmt.Println("error: no access token set")
		os.Exit(2)
	}

	if AccessTokenSecret == "" {
		fmt.Println("error: no access token secret set")
		os.Exit(2)
	}

	ta := &TwitterAPI{
		KeyConsumer:       KeyConsumer,
		KeySecret:         KeySecret,
		AccessToken:       AccessToken,
		AccessTokenSecret: AccessTokenSecret,
	}

	fmt.Println(ta)
}
