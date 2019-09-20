package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var KeyConsumer string
var KeySecret string
var AccessToken string
var AccessTokenSecret string

var KeepFollowing bool
var DumpLikes bool

func init() {
	flag.StringVar(&KeyConsumer, "consumer", "", "Twitter API Consumer Key")
	flag.StringVar(&KeySecret, "secret", "", "Twitter API Secret Key")
	flag.StringVar(&AccessToken, "accesstoken", "", "Twitter API Access Token")
	flag.StringVar(&AccessTokenSecret, "accesstokensecret", "", "Twitter API Access Token Secret")

	flag.BoolVar(&KeepFollowing, "keepfollowing", false, "Keep liked tweets from people you follow")
	flag.BoolVar(&DumpLikes, "dump", false, "Dump all likes to stdout in json format")
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

	if DumpLikes {
		likes, err := ta.GetLikes("imwally")
		if err != nil {
			fmt.Println("error:", err)
		}

		output, err := json.Marshal(likes)
		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Println(string(output))
		return
	}
}
