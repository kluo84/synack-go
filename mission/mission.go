package mission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mBot/data"
	"mBot/env"
	"mBot/requests"
	"net/http"
	"os"
	"sync"
)

func CheckClaimed(target string, header []string, missionurl string) {
	statuscode, body := requests.DoGetRequest(target, header)
	defer body.Close()

	if statuscode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			log.Println(err)
			return
		}
		var claimedMissions data.MissionData
		err = json.Unmarshal(bodyBytes, &claimedMissions)
		if err != nil {
			fmt.Println("Error unmarshalling data:", err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		lenResp := len(claimedMissions)
		if len(bodyString) > 3 { // Check to see if there is a JSON response
			fmt.Println("Starting loop to attempt to grab missions...")
			for i := 0; i < lenResp; i++ {
				missionID := claimedMissions[i].ID
				missionEditURL := missionurl + claimedMissions[i].Organization.ID + "/listings/" + claimedMissions[i].Listing.ID + "/campaigns/" + claimedMissions[i].Campaign.ID + "/tasks/" + missionID
				fmt.Println("Version:", claimedMissions[i].Version)
				EditMission(missionID, missionEditURL, header, claimedMissions[i].Version)
			}
		}
	} else if statuscode == http.StatusUnauthorized {
		log.Printf(env.WarningColor, "Bearer token expired.")
		os.Exit(0)
	}
}

func CheckMissions(target string, header []string) {
	statuscode, body := requests.DoGetRequest(target, header)
	defer body.Close()

	if statuscode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			log.Println(err)
			return
		}

		var newMissions data.MissionData
		err = json.Unmarshal(bodyBytes, &newMissions)
		bodyString := string(bodyBytes)
		lenResp := len(newMissions)
		if len(bodyString) > 3 { // Check to see if there is a JSON response
			var wg sync.WaitGroup
			wg.Add(lenResp)

			fmt.Println("Starting loop to attempt to grab missions...")
			for i := 0; i < lenResp; i++ {
				go func(i int) {
					defer wg.Done()
					log.Println("Saw", newMissions[i].Payout.Amount, newMissions[i].Listing.Title, "mission -", newMissions[i].Title)
					GrabMission("https://platform.synack.com/api/tasks/v1/organizations/"+newMissions[i].Organization.ID+"/listings/"+newMissions[i].Listing.ID+"/campaigns/"+newMissions[i].Campaign.ID+"/tasks/"+newMissions[i].ID+"/transitions", header, newMissions[i].Payout.Amount, newMissions[i].Listing.Title, newMissions[i].Title, 10)
					log.Println("Attempted to grab", newMissions[i].Payout.Amount, newMissions[i].Listing.Title, "mission -", newMissions[i].Title)
				}(i)
			}
			wg.Wait()
			fmt.Println("Finished attempt to grab new missions...")
		}
	} else if statuscode == http.StatusUnauthorized {
		log.Printf(env.WarningColor, "Bearer token expired.")
		os.Exit(0)
	}
}

func EditMission(missionID string, url string, header []string, version int) {
	payload, err := json.Marshal(map[string]interface{}{
		data.SuspectedVulnerability[0]: missionID,
		data.SuspectedVulnerability[1]: "This is a response to test submission via my bot",
		data.SuspectedVulnerability[2]: "not_exploitable",
		data.SuspectedVulnerability[3]: version,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2.
	client := &http.Client{}

	// 3.
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	requests.SetHeaders(req, header)
	if err != nil {
		log.Fatal(err)
	}

	// 4.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// 5.
	defer resp.Body.Close()

	// 6.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}

func GrabMission(target string, header []string, pay int, project string, title string, n int) {
	var jsonStr = []byte(`{"type": "CLAIM"}`)
	statuscode, body := requests.DoPostRequest(target, header, jsonStr)
	defer body.Close()
	if statuscode == http.StatusCreated {
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			log.Println(err)
			return
		}

		bodyString := string(bodyBytes)
		if len(bodyString) > 0 { // Check to see if there is a JSON response
			log.Printf(env.SuccessColor, "Grabbed $ "+fmt.Sprint(pay)+" "+project+" mission - "+title)
		} else {
			log.Printf(env.WarningColor, "Error grabbing $"+fmt.Sprint(pay)+" "+project+" mission - "+title)
		}
	} else if statuscode == http.StatusUnauthorized {
		log.Printf(env.WarningColor, "Bearer token expired.")
		os.Exit(0)
	} else if statuscode == http.StatusInternalServerError {
		if n > 0 {
			log.Printf(env.WarningColor, "500 Error. Retrying...")
			n = n - 1
			GrabMission(target, header, pay, project, title, n)
		} else {
			log.Printf(env.ErrorColor, "500 Error. 10 attempts have been made. Giving up!")
		}
	} else if statuscode == http.StatusForbidden {
		log.Printf(env.ErrorColor, "Too slow! Someone grabbed it.")
	} else if statuscode == http.StatusPreconditionFailed {
		log.Printf(env.ErrorColor, "You already have the maximum number of missions for this target.")
	} else {
		log.Printf(env.WarningColor, "Status Code: "+fmt.Sprint(statuscode))
	}
	body.Close()
}
