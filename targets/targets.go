package targets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mBot/data"
	"mBot/env"
	"mBot/requests"
	"net/http"
	"os"
)

func CheckTargets(target string, header []string) {

	statuscode, body := requests.DoGetRequest(target, header)

	defer body.Close()
	if statuscode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			log.Println(err)
			return
		}

		var newTargets data.TargetData
		err = json.Unmarshal(bodyBytes, &newTargets)
		bodyString := string(bodyBytes)

		if len(bodyString) > 3 { // Check to see if there is a JSON response
			for i := 0; i < len(newTargets); i++ { // Iterate over the number of unregistered targets
				if newTargets[i].Category.ID >= 3 { // Check to see if new target is a web/host target
				} else {
					if !newTargets[i].Registered {
						fmt.Println("Onboarding to", newTargets[i].Codename, "-", newTargets[i].Slug)
						OnboardTarget("https://platform.synack.com/api/targets/"+newTargets[i].Slug+"/signup", header)
					}
				}
			}
		}
	} else if statuscode == http.StatusUnauthorized {
		log.Println("Bearer token expired.")
		os.Exit(0)
	}
}

func OnboardTarget(target string, header []string) {
	var jsonStr = []byte(`{"ResearcherListing":{"terms":1}}`)
	requests.DoPostRequest(target, header, jsonStr)
}

func CheckForQR(target string, header []string) {
	statuscode, body := requests.DoGetRequest(target, header)

	defer body.Close()
	if statuscode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			log.Println(err)
			return
		}

		var newTargets data.TargetData
		err = json.Unmarshal(bodyBytes, &newTargets)
		bodyString := string(bodyBytes)

		if len(bodyString) > 3 { // Check to see if there is a JSON response
			for i := 0; i < len(newTargets); i++ { // Iterate over the number of unregistered targets
				if newTargets[i].Category.ID >= 3 { // Check to see if new target is a web/host target
				} else {
					log.Printf(env.SuccessColor, "[!] Target "+newTargets[i].Codename+" is in QR!")
				}
			}
		}
	} else if statuscode == http.StatusUnauthorized {
		log.Println("Bearer token expired.")
		os.Exit(0)
	}
}
