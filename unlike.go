package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const API string = "https://api.twitter.com/1.1/"
const API_TOKEN string = "https://api.twitter.com/oauth2/token"

type TokenResp struct {
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

func GenerateBearerToken(consumer, secret string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, API_TOKEN, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(consumer, secret)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tr TokenResp
	err = json.Unmarshal(body, &tr)
	if err != nil {
		return "", err
	}

	return tr.Token, nil
}

func GetLikedTweetIds(token string, name string, count int, max int) ([]Tweet, error) {
	endPoint := ""
	if max == 0 {
		endPoint = fmt.Sprintf("%sfavorites/list.json?count=%d&screen_name=%s", API, count, name)
	} else {
		endPoint = fmt.Sprintf("%sfavorites/list.json?count=%d&screen_name=%s&max_id=%d", API, count, name, max)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, endPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	var tweets []Tweet
	err = json.Unmarshal(body, &tweets)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func UnlikeTweet(token string, id int) error {
	endPoint := fmt.Sprintf("%sfavorites/destroy.json?id=%d", API, id)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, endPoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

func main() {
	token, err := GenerateBearerToken(KEY_CONSUMER, KEY_SECRET)
	if err != nil {
		log.Println(err)
	}

	i := 0
	var tweets []Tweet
	tweets, _ = GetLikedTweetIds(token, "imwally", 2, 0)
	for _, tweet := range tweets {
		fmt.Printf("%d: %s\t%d\t%v\n", i, tweet.CreatedAt, tweet.Id, tweet.User)
		i++
	}

}
