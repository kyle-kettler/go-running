package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/kyle-kettler/go-running/models"
)

func GetWeather(coordinates models.Coordinates) models.Weather {
    baseUrl := "https://api.open-meteo.com/v1/forecast?latitude=" + coordinates.Lat + "&longitude=" + coordinates.Lon + "&current=temperature_2m,apparent_temperature,precipitation,rain,showers,snowfall,wind_speed_10m,wind_direction_10m&temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch&forecast_days=1"

    res, err := http.Get(baseUrl)
    if err != nil {
        log.Fatal("error getting weather")
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        log.Fatal("error reading weather body")
    }

    var weather models.Weather
    err = json.Unmarshal(body, &weather)
    if err != nil {
        log.Fatal("Error unmarshalling JSON:", err)
    }

    return weather
}

func GetCompassDirection(degrees float64) string {
    switch {
    case degrees >= 337.5 || degrees < 22.5:
        return "N"
    case degrees >= 22.5 && degrees < 67.5:
        return "NE"
    case degrees >= 67.5 && degrees < 112.5:
        return "E"
    case degrees >= 112.5 && degrees < 157.5:
        return "SE"
    case degrees >= 157.5 && degrees < 202.5:
        return "S"
    case degrees >= 202.5 && degrees < 247.5:
        return "SW"
    case degrees >= 247.5 && degrees < 292.5:
        return "W"
    case degrees >= 292.5 && degrees < 337.5:
        return "NW"
    default:
        return "Invalid degree value"
    }
}
