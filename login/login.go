package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	//"net/url"
	"time"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"golang.org/x/net/html"
	//"github.com/pquerna/otp/totp"
	//"github.com/pquerna/otp"
	"golang.org/x/net/publicsuffix"
	"net/http/cookiejar"
)

type AuthResponse struct {
	AuthToken     string `json:"auth_token"`
	ProgressToken string `json:"progress_token"`
}

type Config struct {
	Password string `json:"password"`
	OtpAuth  string `json:"otpauth"`
}

var client *http.Client

func init() {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	proxyURL, err := url2.Parse("http://127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = &http.Client{
		Jar:       jar,
		Transport: transport,
	}
}

func readConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := homeDir + "/.synack.conf"
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}


func getRequest(url string) (string, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", err
    }

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        body, _ := ioutil.ReadAll(resp.Body)
        return "", fmt.Errorf("GET request failed with status %d: %s", resp.StatusCode, body)
    }

    doc, err := html.Parse(resp.Body)
    if err != nil {
        return "", err
    }

    var csrfToken string
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "meta" {
            isCSRFToken := false
            for _, a := range n.Attr {
                if a.Key == "name" && a.Val == "csrf-token" {
                    isCSRFToken = true
                }
                if isCSRFToken && a.Key == "content" {
                    csrfToken = a.Val
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)

    if csrfToken == "" {
        return "", fmt.Errorf("No csrf-token found in HTML content")
    }

    return csrfToken, nil
}

func generateTOTP(secret string, digits int, period uint) string {
	// Decode the secret
	key, _ := base32.StdEncoding.DecodeString(secret)

	// Calculate the number of time steps
	epochSeconds := uint64(time.Now().Unix())
	timeSteps := epochSeconds / uint64(period)

	// Convert time steps to byte array
	timeStepsBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeStepsBytes, timeSteps)

	// Calculate HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(timeStepsBytes)
	hash := hmacSha1.Sum(nil)

	// Use the last nibble (4 bits) of the hash to get the index for truncation
	offset := hash[19] & 0xF

	// Extract a 4-byte dynamic binary code from the HMAC result at the specified index
	code := binary.BigEndian.Uint32(hash[offset : offset+4])

	// Ensure the code is only the specified number of digits
	code %= uint32(math.Pow10(digits))

	// Convert the code to the desired number of digits as a string
	return fmt.Sprintf(fmt.Sprintf("%%0%dd", digits), code)
}

func postRequest(url string, data map[string]string, csrfToken string) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Csrf-Token", csrfToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 422 {
		log.Println("Error:", string(responseBody))
	}

	return responseBody, nil
}

func main() {
    csrfToken, err := getRequest("https://login.synack.com/")
    if err != nil {
        log.Fatal("Error making GET request:", err)
    }

	config, err := readConfig()
	if err != nil {
		log.Fatal("Error reading config:", err)
	}

	// Using the provided generateTOTP function instead of generateTOTPCode
	code := generateTOTP(config.OtpAuth, 7, 10)

	data1 := map[string]string{
		"email":    "kluosaha@gmail.com",
		"password": config.Password,
	}
	respData1, err := postRequest("https://login.synack.com/api/authenticate", data1, csrfToken)
	if err != nil {
		log.Fatal("Error in first POST request:", err)
	}

	var respObj1 AuthResponse
	err = json.Unmarshal(respData1, &respObj1)
	if err != nil {
		log.Fatal("Error unmarshalling response:", err)
	}

	data2 := map[string]string{
		"authy_token":    code,
		"progress_token": respObj1.ProgressToken,
	}

	respData2, err := postRequest("https://login.synack.com/api/authenticate", data2, csrfToken)
	if err != nil {
		log.Fatal("Error in second POST request:", err)
	}

	var respObj2 AuthResponse
	err = json.Unmarshal(respData2, &respObj2)
	if err != nil {
		log.Fatal("Error unmarshalling response:", err)
	}

	if respObj2.AuthToken != "" {
		fmt.Println("Login successful!")
		fmt.Println("Auth Bearer Token:", respObj2.AuthToken)
	} else {
		log.Fatal("Authentication failed.")
	}
}
