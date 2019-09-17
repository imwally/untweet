package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	KeyConsumer string
	KeySecret   string
	Token       string
}

type TwitterAPIRequest struct {
	Headers    http.Header
	Parameters map[string]string
	Method     string
	EndPoint   string
	Body       string
	Auth       string
}

func printHeaders(resp *http.Response) {
	for headerKey, headerValue := range resp.Header {
		fmt.Printf("%s: %s\n", headerKey, headerValue)
	}
}

func printBody(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(body))
}

func (ta *TwitterAPI) Request(tar *TwitterAPIRequest) ([]byte, error) {
	if tar == nil {
		return nil, errors.New("error: invalid resource")
	}

	// Setup new HTTP client with proper method, end point, and body
	client := &http.Client{}
	req, err := http.NewRequest(tar.Method, tar.EndPoint, strings.NewReader(tar.Body))
	if err != nil {
		return nil, err
	}

	// Set HTTP headers
	if tar.Headers != nil {
		req.Header = tar.Headers
	}

	// Set HTTP basic auth if needed
	if tar.Auth == "basic" {
		req.SetBasicAuth(ta.KeyConsumer, ta.KeySecret)
	}

        if tar.Auth == "application" {
                req.Header.Add("Authorization", "Bearer "+ta.Token)
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

func UnmarshalToken(b []byte) string {
	var tr TokenResponse
	json.Unmarshal(b, &tr)

	return tr.Token
}

func NewRequest(resource string, parameters map[string]string) *TwitterAPIRequest {

	params := ""
	for k, v := range parameters {
		params += fmt.Sprintf("%s=%s&", k, v)
	}

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
			EndPoint: API_URL + resource + ".json?" + params,
			Auth:     "application",
		}
	}

	return nil
}

func main() {
	ta := &TwitterAPI{
		KeyConsumer: "",
		KeySecret:   "",
	}

	btoken, err := ta.Request(NewRequest("oauth2/token", nil))
	if err != nil {
		log.Println(err)
	}

	ta.Token = UnmarshalToken(btoken)

        fmt.Println(ta)
}
