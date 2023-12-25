package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

type InviteGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Members     int    `json:"approximate_member_count"`
	Online      int    `json:"approximate_presence_count"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Response struct {
	Code        string      `json:"code"`
	Guild       InviteGuild `json:"guild"`
	Valid       bool        `json:"valid"`
	ExpiresAt   string      `json:"expires_at"`
	Description string      `json:"description"`
}

func FetchServerInfo(inviteCode string) error {
	resp, err := http.Get(fmt.Sprintf("https://discord.com/api/v9/invites/%s?with_counts=true", inviteCode))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var responseJSON Response
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		return err
	}

	Expiry, err := time.Parse(time.RFC3339, responseJSON.ExpiresAt)
	if err != nil {
		return err
	}

	est, err := time.LoadLocation("America/New_York") // EST time
	if err != nil {
		return err
	}

	Expiry = Expiry.In(est)

	// fmt.Printf("Valid: %t\n", responseJSON.Valid) // todo: fix invalid response
	fmt.Printf("\x1b[95m[\x1b[97mServer Name\x1b[95m] \x1b[97m%s\x1b[95m\n", responseJSON.Guild.Name)
	fmt.Printf("\x1b[95m[\x1b[97mInvite Code\x1b[95m] \x1b[97m%s\x1b[95m\n", responseJSON.Code)
	fmt.Printf("\x1b[95m[\x1b[97mInvite Expiry\x1b[95m (\x1b[97mEST\x1b[95m)]\x1b[97m %s\n", Expiry.Format("Monday, January 2, 2006 (1/2/06) at 03:04:05 PM"))
	fmt.Printf("\x1b[95m[\x1b[97mServer ID\x1b[95m]\x1b[95m \x1b[97m%s\x1b[95m\n", responseJSON.Guild.ID)
	// fmt.Printf("\x1b[95m[\x1b[97mServer Members\x1b[95m] \x1b[97m%d\x1b[95m\n", responseJSON.Guild.Members) // todo: fix invalid response
	// fmt.Printf("\x1b[95m[\x1b[97mServer Online\x1b[95m] \x1b[97m%d\x1b[95m\n", responseJSON.Guild.Online)   // todo: fix invalid response
	ReadInput()
	return nil
}

func main() {
	Config := "settings.json"

	settingsData, err := ioutil.ReadFile(Config)
	if err != nil {
		fmt.Println("Error reading settings file:", err)
		os.Exit(1)
	}

	var settings map[string]string
	err = json.Unmarshal(settingsData, &settings)
	if err != nil {
		fmt.Println("Error parsing settings:", err)
		os.Exit(1)
	}

	inviteCode, ok := settings["inviteCode"]
	if !ok || inviteCode == "" {
		fmt.Println("settings.json is empty.")
		os.Exit(1)
	}

	// lazy so we will just remove the link to get the code
	regex := regexp.MustCompile(`(https?:\/\/)?(www\.)?(discord\.gg|discord\.com\/invite)\/`)
	inviteCode = regex.ReplaceAllString(inviteCode, "")

	err = FetchServerInfo(inviteCode)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func ReadInput() string {
	var input string
	fmt.Scanln(&input)
	//Clear()
	return input
}
