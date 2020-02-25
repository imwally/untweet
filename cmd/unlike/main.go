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
	key         string
	keySecret   string
	token       string
	tokenSecret string

	keepFollowing bool
	dumpLikes     bool
)

func init() {
	flag.StringVar(&key, "key", "", "Twitter API Consumer Key")
	flag.StringVar(&keySecret, "key-secret", "", "Twitter API Secret Key")
	flag.StringVar(&token, "token", "", "Twitter API Access Token")
	flag.StringVar(&tokenSecret, "token-secret", "", "Twitter API Access Token Secret")

	flag.BoolVar(&keepFollowing, "keep-following", false, "Don't unlike any tweets from people you follow")
	flag.BoolVar(&dumpLikes, "dump", false, "Dump all likes to stdout in json format")
	flag.Parse()
}

func main() {

	if key == "" {
		if key = os.Getenv("TWITTER_API_KEY"); key == "" {
			fmt.Fprintf(os.Stderr, "error: no api key set\n")
			os.Exit(2)
		}
	}

	if keySecret == "" {
		if keySecret = os.Getenv("TWITTER_API_KEY_SECRET"); keySecret == "" {
			fmt.Fprintf(os.Stderr, "error: no secret key set\n")
			os.Exit(2)
		}
	}

	if token == "" {
		if token = os.Getenv("TWITTER_API_TOKEN"); token == "" {
			fmt.Fprintf(os.Stderr, "error: no access token set\n")
			os.Exit(2)
		}
	}

	if tokenSecret == "" {
		if tokenSecret = os.Getenv("TWITTER_API_TOKEN_SECRET"); tokenSecret == "" {
			fmt.Fprintf(os.Stderr, "error: no access token secret set\n")
			os.Exit(2)
		}
	}

	ta := &tapi.TwitterAPI{
		KeyConsumer:       key,
		KeySecret:         keySecret,
		AccessToken:       token,
		AccessTokenSecret: tokenSecret,
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
