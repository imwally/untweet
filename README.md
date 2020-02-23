# unlike

A little command line utility that unlikes all (or some of) the tweets you
previously liked.

## Requirements

To communicate with Twitter's API you will need the following:

- Consumer Key
- Consumer Secret
- Access Token
- Access Token Secret

Which means you will need to create an app through [Twitter's developer
portal](https://developer.twitter.com/en/docs/basics/getting-started).

## How to Use

The 4 keys from above are required for every call to `unlike`. 

```
Usage of unlike:
  -access-token string
        Twitter API Access Token
  -access-token-secret string
        Twitter API Access Token Secret
  -consumer string
        Twitter API Consumer Key
  -dump
        Dump all likes to stdout in json format
  -keep-following
        Keep liked tweets from people you follow
  -secret string
        Twitter API Secret Key
```

### Dump all likes to stdout in json format

```
$ unlike -consumer "xxxxxxxxx" \
        -secret "xxxxxxxxx" \
        -access-token "xxxxxxxxx" \
        -access-token-secret "xxxxxxxxx" \
        -dump
```

### Unlike only the tweets from people you don't follow

```
$ unlike -consumer "xxxxxxxxx" \
        -secret "xxxxxxxxx" \
        -access-token "xxxxxxxxx" \
        -access-token-secret "xxxxxxxxx" \
        -keep-following
```

### Unlike __ALL__ tweets

Please note the lack of any other arguments. Running the command with the
required keys and without any other arguments will unlike __all__ tweets. Don't
worry, you'll be asked for confirmation.

```
$ unlike -consumer "xxxxxxxxx" \
        -secret "xxxxxxxxx" \
        -access-token "xxxxxxxxx" \
        -access-token-secret "xxxxxxxxx"
```
