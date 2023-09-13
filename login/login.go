package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
)

// Structs for our expected response and config
type AuthResponse struct {
	AuthToken     string `json:"auth_token"`
	ProgressToken string `json:"progress_token"`
}

type Config struct {
	Password  string `json:"password"`
	OTPAUTHURL string `json:"otpauth_url"`
}

func readConfig() (*Config, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Construct the path to the config file
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

func postRequest(url string, data map[string]string) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Setup proxy
	proxyURL, err := url2.Parse("http://127.0.0.1:8080")
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
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

	// Check for the 422 status code after reading the response body
	if resp.StatusCode == 422 {
		log.Println("Error:", string(responseBody))
	}

	return responseBody, nil
}

func GenerateTOTPCode(otpURL string) (string, error) {
	urlObj, err := url2.Parse(otpURL)
	if err != nil {
		return "", err
	}

	if urlObj.Scheme != "otpauth" || urlObj.Host != "totp" {
		return "", fmt.Errorf("Only otpauth://totp/ URLs are supported currently.")
	}

	secretQuery := urlObj.Query().Get("secret")
	if secretQuery == "" {
		return "", fmt.Errorf("No secret parameter found in the URL.")
	}

	return totp.GenerateCode(secretQuery, time.Now())
}

func main() {
	// Read the config
	config, err := readConfig()
	if err != nil {
		log.Fatal("Error reading config:", err)
	}

	// Generate TOTP token
	authyToken, err := GenerateTOTPCode(config.OTPAUTHURL)
	if err != nil {
		log.Fatal("Error generating TOTP token:", err)
	}

	// First POST request
	data1 := map[string]string{
		"email":    "kluosaha@gmail.com",
		"password": config.Password,
	}
	respData1, err := postRequest("https://login.synack.com/api/authenticate", data1)
	if err != nil {
		log.Fatal("Error in first POST request:", err)
	}

	var respObj1 AuthResponse
	err = json.Unmarshal(respData1, &respObj1)
	if err != nil {
		log.Fatal("Error unmarshalling response:", err)
	}

	// Second POST request
	data2 := map[string]string{
		"authy_token":    authyToken,
		"progress_token": respObj1.ProgressToken,
	}
	respData2, err := postRequest("https://login.synack.com/api/authenticate", data2)
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

