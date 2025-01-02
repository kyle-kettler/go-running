package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/kyle-kettler/go-running/models"
)

func GetCoordinates(city, state, country string) models.Coordinates {
    baseURL := "https://nominatim.openstreetmap.org/search"
    u, err := url.Parse(baseURL)
    if err != nil {
        log.Fatal("Error parsing URL")
    }

    q := url.Values{}
    q.Add("q", fmt.Sprintf("%s %s %s", city, state, country))
    q.Add("format", "jsonv2")

    u.RawQuery = q.Encode()

    res, err := http.Get(u.String())
    if err != nil {
        log.Fatal("Nominatim response error")
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        log.Fatal("Error reading response body")
    }
    defer res.Body.Close()

    var coordinates []models.Coordinates
    err = json.Unmarshal(body, &coordinates)
    if err != nil {
        log.Fatal("Error unmarshalling JSON:", err)
    }

    if len(coordinates) == 0 {
        log.Fatal("No coordinates found")
    }

    return coordinates[0]
}
