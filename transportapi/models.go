package transportapi

import (
	"encoding/json"
	"fmt"
	"time"
)

// OptionalTime wraps a time.Time pointer to allow JSON unmarshaling.
type OptionalTime struct {
	Time *time.Time
}

// UnmarshalJSON unmarshals a HH:mm time string into a OptionalTime
func (t *OptionalTime) UnmarshalJSON(data []byte) error {
	var timeString string
	err := json.Unmarshal(data, &timeString)
	if err != nil {
		fmt.Println("DEBUG couldn't unmarshal:", string(data))
		return err
	}

	if timeString == "" {
		return nil
	}

	parsed, err := time.Parse("15:04", timeString)
	if err != nil {
		fmt.Println("DEBUG couldn't parse", timeString, "as a time")
		return err
	}

	year, month, day := time.Now().Date()
	parsed = parsed.AddDate(year, int(month)-1, day-1)

	t.Time = &parsed

	return nil
}

// LiveTrainUpdateDeparture represents a single departure in a set of departures.
type LiveTrainUpdateDeparture struct {
	Mode                      string       `json:"mode"`
	Platform                  string       `json:"platform"`
	Operator                  string       `json:"operator"`
	OperatorName              string       `json:"operator_name"`
	AimedDepartureTime        OptionalTime `json:"aimed_departure_time"`
	AimedArrivalTime          OptionalTime `json:"aimed_arrival_time"`
	DestinationName           string       `json:"destination_name"`
	Status                    string       `json:"status"`
	ExpectedArrivalTime       OptionalTime `json:"expected_arrival_time"`
	ExpectedDepartureTime     OptionalTime `json:"expected_departure_time"`
	BestArrivalEstimateMins   int          `json:"best_arrival_estimate_mins"`
	BestDepartureEstimateMins int          `json:"best_departure_estimate_mins"`
}

// LiveTrainUpdateDepartures contains set(s) of departures.
type LiveTrainUpdateDepartures struct {
	All []LiveTrainUpdateDeparture `json:"all"`
}

// LiveTrainUpdate contains a set of departures from a station.
type LiveTrainUpdate struct {
	StationName string                    `json:"station_name"`
	StationCode string                    `json:"station_code"`
	Departures  LiveTrainUpdateDepartures `json:"departures"`
}
