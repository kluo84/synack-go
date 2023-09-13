package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <otpauth-url>", os.Args[0])
	}

	otpURL := os.Args[1]
	urlObj, err := url.Parse(otpURL)
	if err != nil {
		log.Fatalf("Failed to parse the otpauth URL: %v", err)
	}

	if urlObj.Scheme != "otpauth" || urlObj.Host != "totp" {
		log.Fatalf("Only otpauth://totp/ URLs are supported currently.")
	}

	secretQuery := urlObj.Query().Get("secret")
	if secretQuery == "" {
		log.Fatalf("No secret parameter found in the URL.")
	}

	code, err := totp.GenerateCode(secretQuery, time.Now())
	if err != nil {
		log.Fatalf("Failed to generate TOTP code: %v", err)
	}

	fmt.Println("Generated TOTP code:", code)
}

