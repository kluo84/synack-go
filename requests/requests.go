package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"
)

func SetHeaders(req *http.Request, header []string) {
	req.Header.Set("User-Agent", header[0])
	req.Header.Set("Authorization", header[1])
	req.Header.Set("Sec-Fetch-Site", header[2])
	req.Header.Set("Sec-Fetch-Mode", header[3])
	req.Header.Set("Referer", header[4])
	req.Header.Set("X-CSRF-Token", header[5])
	req.Header.Set("Content-Type", header[6])
	req.Header.Set("Connection", header[7])
}

func DoGetRequest(target string, header []string) (int, io.ReadCloser) {
	ctx, cancel := context.WithTimeout(context.Background(), 30000*time.Millisecond)
	defer cancel()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	SetHeaders(req, header)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return resp.StatusCode, resp.Body
}

func DoPostRequest(target string, header []string, jsonStr []byte) (int, io.ReadCloser) {
	ctx, cancel := context.WithTimeout(context.Background(), 30000*time.Millisecond)
	defer cancel()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("POST", target, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
	}

	SetHeaders(req, header)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Println(err)
	}
	return resp.StatusCode, resp.Body
}
