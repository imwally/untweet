package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/imwally/untweet/tapi"
)

var (
	key         string
	keySecret   string
	token       string
	tokenSecret string

	keepFollowing bool
	older         string
)

func init() {
	flag.StringVar(&key, "key", "", "Twitter API Consumer Key")
	flag.StringVar(&keySecret, "key-secret", "", "Twitter API Secret Key")
	flag.StringVar(&token, "token", "", "Twitter API Access Token")
	flag.StringVar(&tokenSecret, "token-secret", "", "Twitter API Access Token Secret")
	flag.Parse()
}

func main() {
	tweetsCmd := flag.NewFlagSet("tweets", flag.ExitOnError)
	tweetsOlder := tweetsCmd.Duration("older", time.Second*0, "Destroy tweets older than this time (30m, 24h, 48h, etc..)")

	likesCmd := flag.NewFlagSet("likes", flag.ExitOnError)
	keepFollowing := likesCmd.Bool("keep-following", false, "Don't destroy likes of tweets from people you follow")

	dumpCmd := flag.NewFlagSet("dump", flag.ExitOnError)
	dumpLikes := dumpCmd.Bool("likes", false, "Dump likes")
	dumpTweets := dumpCmd.Bool("tweets", false, "Dump tweets")

	if len(os.Args) < 2 || os.Args[1] == "help" {
		fmt.Println("USAGE:")
		fmt.Printf("    untweet command [command options]\n\n")

		fmt.Println("COMMAND:")
		fmt.Printf("    dump \t Dump likes or tweets\n")
		fmt.Printf("    tweets \t Destroy tweets\n")
		fmt.Printf("    likes \t Destroy likes\n\n")

		fmt.Println("OPTIONS:")
		fmt.Println("    Use -h on each command to view options")

		return
	}

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

	switch os.Args[1] {
	case "dump":
		dumpCmd.Parse(os.Args[2:])

		if *dumpLikes {
			likes, err := ta.GetLikes()
			if err != nil {
				log.Println(err)
			}

			output, err := json.Marshal(likes)
			if err != nil {
				log.Println(err)
			}

			fmt.Println(string(output))

		}

		if *dumpTweets {
			tweets, err := ta.GetTweets()
			if err != nil {
				log.Println(err)
			}

			output, err := json.Marshal(tweets)
			if err != nil {
				log.Println(err)
			}

			fmt.Println(string(output))
		}

	case "tweets":
		tweetsCmd.Parse(os.Args[2:])

		fmt.Printf("Destroy all of your tweets older than %s? [y/N]: ", *tweetsOlder)

		var proceed string
		fmt.Scanln(&proceed)

		if proceed != "y" && proceed != "Y" {
			return
		}

		for max, next := 0, 1; next > 0; {
			batch, err := ta.GetBatchedTweets(max)
			if err != nil {
				log.Println(err)
			}

			for _, tweet := range batch {
				tweetTime, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.CreatedAt)
				if err != nil {
					log.Println(err)
				}

				// Destroy tweets older than specified time
				if time.Since(tweetTime) > *tweetsOlder {
					err := ta.DestroyTweet(tweet.Id)
					if err != nil {
						log.Println(err)
					}
				}
			}

			blen := len(batch)
			if blen > 0 {
				max = batch[blen-1].Id - 1
			}

			next = blen
		}

	case "likes":
		likesCmd.Parse(os.Args[2:])

		if *keepFollowing {
			fmt.Printf("Unlike tweets from people you don't follow? [y/N]: ")
		} else {
			fmt.Printf("Unlike ALL tweets? [y/N]: ")
		}

		var proceed string
		fmt.Scanln(&proceed)

		if proceed != "y" && proceed != "Y" {
			return
		}

		for max, next := 0, 1; next > 0; {
			batch, err := ta.GetBatchedLikes(max)
			if err != nil {
				log.Println(err)
			}

			for _, like := range batch {
				if *keepFollowing && like.Following {
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
}
