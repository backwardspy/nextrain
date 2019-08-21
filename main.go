package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/backwardspy/nextrain/transportapi"
)

func getNewCredentials() (appID string, key string) {
	fmt.Println("You will need an app on the 3scale Transport API Developer Portal.")
	fmt.Println()
	fmt.Println("    https://developer.transportapi.com/admin/applications/new")
	fmt.Println()
	fmt.Println("When the app has been created, enter the App ID and Key below.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter App ID: ")
	appID, err := reader.ReadString('\n')
	if err != nil {
		log.Panic("failed to read appID")
	}
	appID = strings.TrimSpace(appID)

	fmt.Print("Enter Key: ")
	key, err = reader.ReadString('\n')
	if err != nil {
		log.Panic("failed to read key")
	}
	key = strings.TrimSpace(key)

	return
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %v FROM TO", os.Args[0])
		return
	}

	station_from := strings.ToUpper(os.Args[1])
	station_to := strings.ToUpper(os.Args[2])

	api := transportapi.New()

	err := api.LoadCredentials()
	if err != nil {
		appID, key := getNewCredentials()
		api.Authenticate(appID, key)
		api.SaveCredentials()
	}

	updates, err := api.TrainUpdatesLive(station_from, station_to)
	if err != nil {
		return
	}

	for _, departure := range updates.Departures.All {
		aim := departure.AimedDepartureTime.Time.Format("15:04")
		expected := departure.ExpectedDepartureTime.Time.Format("15:04")
		parts := []string{
			fmt.Sprintf("%v from %v to %v ", strings.Title(departure.Mode), updates.StationName, departure.DestinationName),
		}

		if departure.Platform != "" {
			parts = append(parts, fmt.Sprintf("on platform %v ", departure.Platform))
		}

		parts = append(parts, fmt.Sprintf("at %v (expected at %v). Leaving in %v minutes.", aim, expected, departure.BestDepartureEstimateMins))

		fmt.Println(strings.Join(parts, ""))
	}
}
