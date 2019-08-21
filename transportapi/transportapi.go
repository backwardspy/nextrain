package transportapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/shibukawa/configdir"
)

var configDir = configdir.New("backwardspy", "nextrain").QueryFolders(configdir.Global)[0]

const authFile = "auth.json"

type auth struct {
	AppID string
	Key   string
}

// TransportAPI exposes methods for calling the 3scale Transport API.
type TransportAPI struct {
	auth    auth
	http    http.Client
	baseURL *url.URL
}

// New initialises a new 3scale Transport API client.
func New() TransportAPI {
	var api TransportAPI
	api.http = http.Client{
		Timeout: time.Second * 3,
	}

	baseURL, err := url.Parse("http://transportapi.com")
	if err != nil {
		log.Panic("failed to parse URL http://transportapi.com")
	}

	api.baseURL = baseURL

	return api
}

// SaveCredentials persists API credentials to user's config directory to avoid asking in future.
func (api *TransportAPI) SaveCredentials() {
	data, err := json.Marshal(api.auth)
	if err != nil {
		log.Panic("failed to marshal api.auth as json")
	}

	err = configDir.WriteFile(authFile, data)
	if err != nil {
		fmt.Println("WARNING: Failed to persist API authentication credentials to config directory:", configDir.Path)
	}
}

// LoadCredentials attempts to load API auth credentials from the user's config directory.
func (api *TransportAPI) LoadCredentials() error {
	data, err := configDir.ReadFile(authFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &api.auth)
	if err != nil {
		fmt.Println(
			"WARNING: Failed to load API authentication credentials from",
			filepath.Join(configDir.Path, authFile),
		)
		fmt.Println("The file appears to be invalid. Correct the issue or delete the file and try again.")
		return err
	}

	return nil
}

// Authenticate sets up the API to use the given authentication credentials.
func (api *TransportAPI) Authenticate(appID string, key string) {
	api.auth = auth{appID, key}
}

func (api *TransportAPI) get(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Failed to create GET", url)
		return nil, err
	}

	// fmt.Println(req.Method, req.URL)

	resp, err := api.http.Do(req)
	if err != nil {
		fmt.Println("Failed to GET", url)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body from GET", url)
		return nil, err
	}

	return body, nil
}

func (api *TransportAPI) makeURL(endpoint string) string {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		log.Panic("failed to parse endpoint", endpoint)
	}

	qs := endpointURL.Query()
	qs.Set("app_id", api.auth.AppID)
	qs.Set("app_key", api.auth.Key)
	endpointURL.RawQuery = qs.Encode()

	return api.baseURL.ResolveReference(endpointURL).String()
}

// TrainUpdatesLive gets live service updates at a given station: departures, arrivals or passes.
func (api *TransportAPI) TrainUpdatesLive(stationCode string, callingAt string) (*LiveTrainUpdate, error) {
	endpoint := fmt.Sprintf("/v3/uk/train/station/%v/live.json?calling_at=%v", stationCode, callingAt)
	url := api.makeURL(endpoint)
	data, err := api.get(url)
	if err != nil {
		fmt.Println("ERROR: Failed to get live train updates:", err)
		return nil, err
	}

	var update LiveTrainUpdate
	err = json.Unmarshal(data, &update)
	if err != nil {
		fmt.Println("ERROR: Failed to parse live train updates:", err)
		return nil, err
	}

	return &update, nil
}
