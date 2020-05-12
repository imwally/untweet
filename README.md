# untweet

A little command line utility that destroys likes and tweets.

## Install

### Binaries

Check out the [latest
release](https://github.com/imwally/untweet/releases/latest) for macOS
and Linux binaries.

### Go

`go get -u github.com/imwally/untweet/cmd/untweet`

## Requirements

To communicate with Twitter's API you will need the following:

- Consumer Key
- Consumer Secret
- Access Token
- Access Token Secret

Which means you will need to create an app through [Twitter's developer
portal](https://developer.twitter.com/en/docs/basics/getting-started).

## How to Use

The 4 keys from above are required for every call to
`untweet`. Environment variables are also supported:

```
TWITTER_API_KEY
TWITTER_API_KEY_SECRET
TWITTER_API_TOKEN
TWITTER_API_TOKEN_SECRET
```

or use command line flags:

```
  -key string
    	Twitter API Consumer Key
  -key-secret string
    	Twitter API Secret Key
  -token string
    	Twitter API Access Token
  -token-secret string
    	Twitter API Access Token Secret
```

### Usage

```
USAGE:
    untweet command [command options]

COMMAND:
    dump         Dump likes or tweets
    tweets       Destroy tweets
    likes        Destroy likes

OPTIONS:
    Use -h on each command to view options
```

## Backup Before You Destroy

### Dump all likes or tweets to stdout in json format

```
$ untweet -likes
```

```
$ untweet -tweets
```

## Other Useful Information

### Keep likes of tweets from people you follow

```
$ untweet likes -keep-following
```
