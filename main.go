package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"mBot/env"
	"mBot/mission"
	"mBot/targets"
)

func getBearerToken(clientID, clientSecret string) (string, error) {
	resp, err := http.PostForm("https://platform.synack.com/oauth/token",
		url.Values{
			"client_id":     {clientID},
			"client_secret": {clientSecret},
			"grant_type":    {"client_credentials"},
		})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}

func main() {
	clientID := flag.String("client_id", "", "Client ID for API authentication")
	clientSecret := flag.String("client_secret", "", "Client Secret for API authentication")
	delay := flag.Uint("d", 60, "Time (in seconds) between requests")
	flag.Parse()

	token, err := getBearerToken(*clientID, *clientSecret)
	if err != nil {
		log.Fatalf("Failed to retrieve bearer token: %v", err)
	}

	headers := []string{
		"mBot/1.0", "Bearer " + token, "same-origin", "cors",
		"https://platform.synack.com/tasks/user/available",
		"xxxx", "application/json", "close",
	}

	urls := []string{
                // urls[0] = all unregistered targets
                "https://platform.synack.com/api/targets?filter%5Bprimary%5D=unregistered&filter%5Bsecondary%5D=all&filter%5Bcategory%5D=all&filter%5Bindustry%5D=all&sorting%5Bfield%5D=dateUpdated&sorting%5Bdirection%5D=desc",
                // urls[1] = available missions sorted by price
                "https://platform.synack.com/api/tasks/v1/tasks?sortBy=price-sort-desc&withHasBeenViewedInfo=true&status=PUBLISHED&page=0&pageSize=20",
                // urls[2] = QR window
                "https://platform.synack.com/api/targets?filter%5Bprimary%5D=all&filter%5Bsecondary%5D%5B%5D=a&filter%5Bsecondary%5D%5B%5D=l&filter%5Bsecondary%5D%5B%5D=l&filter%5Bsecondary%5D%5B%5D=quality_period&filter%5Bcategory%5D=all&filter%5Bindustry%5D=all&sorting%5Bfield%5D=dateUpdated&sorting%5Bdirection%5D=desc",
                // urls[3] = claimed missions
                "https://platform.synack.com/api/tasks/v1/tasks?withHasBeenViewedInfo=true&status=CLAIMED&page=0&pageSize=20",
                // urls[4] = beginning of URL to edit missions
                "https://platform.synack.com/api/tasks/v1/organizations/",
        }

	for {
		log.Printf(env.InfoColor, "Checking in...")
		targets.CheckTargets(urls[0], headers)
		mission.CheckMissions(urls[1], headers)
		targets.CheckForQR(urls[2], headers)

		secs := time.Duration(*delay) * time.Second
		time.Sleep(secs)
	}
}
