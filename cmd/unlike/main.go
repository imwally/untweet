package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/imwally/unlike/tapi"
)

var (
	keyConsumer       string
	keySecret         string
	accessToken       string
	accessTokenSecret string

	keepFollowing bool
	dumpLikes     bool
)

func init() {
	flag.StringVar(&keyConsumer, "consumer", "", "Twitter API Consumer Key")
	flag.StringVar(&keySecret, "secret", "", "Twitter API Secret Key")
	flag.StringVar(&accessToken, "accesstoken", "", "Twitter API Access Token")
	flag.StringVar(&accessTokenSecret, "accesstokensecret", "", "Twitter API Access Token Secret")

	flag.BoolVar(&keepFollowing, "keepfollowing", false, "Don't unlike any tweets from people you follow")
	flag.BoolVar(&dumpLikes, "dump", false, "Dump all likes to stdout in json format")
	flag.Parse()
}

func main() {
	if keySecret == "" {
		fmt.Fprintf(os.Stderr, "error: no secret key set\n")
		os.Exit(2)
	}

	if keyConsumer == "" {
		fmt.Fprintf(os.Stderr, "error: no consumer key set\n")
		os.Exit(2)
	}

	if accessToken == "" {
		fmt.Fprintf(os.Stderr, "error: no access token set\n")
		os.Exit(2)
	}

	if accessTokenSecret == "" {
		fmt.Fprintf(os.Stderr, "error: no access token secret set\n")
		os.Exit(2)
	}

	ta := &tapi.TwitterAPI{
		KeyConsumer:       keyConsumer,
		KeySecret:         keySecret,
		AccessToken:       accessToken,
		AccessTokenSecret: accessTokenSecret,
	}

	// If dump is specified then ONLY dump likes and disregard other flags
	if dumpLikes {
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
	if keepFollowing {
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
			if keepFollowing && like.Following {
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
