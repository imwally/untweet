# untweet

A little command line utility that destroys likes and tweets.

## Install

### Binaries

Check out the [latest
release](https://github.com/imwally/unlike/releases/latest) for macOS
and Linux binaries.

### Go

`go get -u github.com/imwally/unlike/cmd/unlike`

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
`unlike`. Environment variables are also supported:

```
TWITTER_API_KEY
TWITTER_API_KEY_SECRET
TWITTER_API_TOKEN
TWITTER_API_TOKEN_SECRET
```

```
Usage of untweet:
  -destroy-likes
    	Destroy your likes
  -destroy-tweets
    	Destroy your tweets
  -dump-likes
    	Dump all likes to stdout in json format
  -dump-tweets
    	Dump all of your tweets to stdout in json format
  -keep-following
    	Don't destroy likes of tweets from people you follow
  -key string
    	Twitter API Consumer Key
  -key-secret string
    	Twitter API Secret Key
  -token string
    	Twitter API Access Token
  -token-secret string
    	Twitter API Access Token Secret
```

## Backup Before You Destroy

### Dump all likes or tweets to stdout in json format

```
$ untweet -dump-likes
```

```
$ untweet -dump-tweets
```

### Keep likes of people you follow

```
$ untweet -destroy-likes -keep-following
```
