package main

import (
	"fmt"
	"net/http"
)

func main() {
	url := "https://login.synack.com/"

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	cookieHeader := constructCookieHeader(response.Cookies())
	fmt.Println("Cookie:", cookieHeader)
}

func constructCookieHeader(cookies []*http.Cookie) string {
	var cookieStrings []string
	for _, cookie := range cookies {
		cookieStrings = append(cookieStrings, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	return join(cookieStrings, "; ")
}

func join(strs []string, sep string) string {
	var result string
	for i, str := range strs {
		result += str
		if i < len(strs)-1 {
			result += sep
		}
	}
	return result
}

