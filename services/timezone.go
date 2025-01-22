package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/kyle-kettler/go-running/models"
)

func GetTimezone(coordinates models.Coordinates) string {
	baseURL := "https://timeapi.io/api/time/current/coordinate"
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal("Error parsing URL")
	}

	q := url.Values{}
	q.Add("latitude", coordinates.Lat)
	q.Add("longitude", coordinates.Lon)

	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal("Timeapi response error")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading response body")
	}
	defer res.Body.Close()

	var tzResponse models.TimezoneResponse
	err = json.Unmarshal(body, &tzResponse)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	return tzResponse.TimeZone
}
