package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
)

func PrintHeaders(resp *http.Response) {
	for headerKey, headerValue := range resp.Header {
		fmt.Printf("%s: %s\n", headerKey, headerValue)
	}
}

func PrintBody(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(body))
}

func GenerateParameterString(parameters map[string]string, sorted bool) string {
	params := ""

	if sorted {
		keys := make([]string, len(params))
		for k := range parameters {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			params += fmt.Sprintf("%s=%s&", k, parameters[k])
		}
	} else {
		for k, v := range parameters {
			params += fmt.Sprintf("%s=%s&", k, v)
		}
	}

	return params[:len(params)-1]
}

func GenerateNonce() (string, error) {
	n := make([]byte, 32)
	_, err := rand.Read(n)
	if err != nil {
		return "", err
	}
	nonce := base64.StdEncoding.EncodeToString(n)

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}

	nonce = reg.ReplaceAllString(nonce, "")

	return nonce, nil
}
