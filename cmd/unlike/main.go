package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/imwally/unlike/tapi"
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

}

func main() {
	ta := &tapi.TwitterAPI{
		KeyConsumer:       KeyConsumer,
		KeySecret:         KeySecret,
		AccessToken:       AccessToken,
		AccessTokenSecret: AccessTokenSecret,
	}

	// If dump is specified then ONLY dump likes and disregard other flags
	if DumpLikes {
		likes, err := ta.GetLikes()
		if err != nil {
			log.Println(err)
		}

		output, err := json.Marshal(likes)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(string(output))
		return
	}

	// Make sure user knows which likes will be destroyed
	if KeepFollowing {
		fmt.Printf("Unlike tweets from people you don't follow? [y/n]: ")
	} else {
		fmt.Printf("Unlike ALL tweets? [y/n]: ")
	}

	var proceed string
	fmt.Scanln(&proceed)

	if proceed != "y" && proceed != "Y" {
		return
	}

	// Proceed to destroy likes in batches of 200
	for max, next := 0, 1; next > 0; {
		batch, err := ta.GetBatchedLikes(max)
		if err != nil {
			log.Println(err)
		}

		for _, like := range batch {
			if KeepFollowing && like.Following {
				log.Println("keeping", like.Id)
				continue
			}

			err := ta.DestroyLike(like.Id)
			if err != nil {
				log.Println(err)
			}
		}

		blen := len(batch)
		if blen > 0 {
			max = batch[blen-1].Id - 1
		}

		next = blen
	}
}
