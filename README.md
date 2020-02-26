# unlike

A little command line utility that unlikes all (or some of) the tweets you
previously liked.

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
TWITTER_TOKEN
TWITTER_TOKEN_SECRET
```

```
Usage of ./unlike:
  -dump
        Dump all likes to stdout in json format
  -keep-following
        Don't unlike any tweets from people you follow
  -key string
        Twitter API Consumer Key
  -key-secret string
        Twitter API Secret Key
  -token string
        Twitter API Access Token
  -token-secret string
        Twitter API Access Token Secret
```

### Dump all likes to stdout in json format

```
$ unlike -dump
```

### Unlike only the tweets from people you don't follow

```
$ unlike -keep-following
```

### Unlike __ALL__ tweets

Please note the lack of any other arguments. Running the command with the
required keys and without any other arguments will unlike __all__ tweets. Don't
worry, you'll be asked for confirmation.

```
$ unlike
```