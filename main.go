package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Coordinates struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type CurrentWeather struct {
	Time          string  `json:"time"`
	Temp          float64 `json:"temperature_2m"`
	FeelsLike     float64 `json:"apparent_temperature"`
	Precipitation float64 `json:"precipitation"`
	Rain          float64 `json:"rain"`
	Showers       float64 `json:"showers"`
	Snow          float64 `json:"snow"`
	WindSpeed     float64 `json:"wind_speed_10m"`
	WindDirection float64 `json:"wind_direction_10m"`
}

type Weather struct {
	Current CurrentWeather
}

func getCoordinates(city, state, country string) Coordinates {
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

	var coordinates []Coordinates
	err = json.Unmarshal(body, &coordinates)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}

	if len(coordinates) == 0 {
		log.Fatal("No coordinates found")
	}

	return coordinates[0]
}

func getWeather(coordinates Coordinates) Weather {
	baseUrl := "https://api.open-meteo.com/v1/forecast?latitude=" + coordinates.Lat + "&longitude=" + coordinates.Lon + "&current=temperature_2m,apparent_temperature,precipitation,rain,showers,snowfall,wind_speed_10m,wind_direction_10m&temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch&forecast_days=1"

	res, err := http.Get(baseUrl)
	if err != nil {
		log.Fatal("error getting weather")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("error reading weather body")
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}

	return weather
}

func getCompassDirection(degrees float64) string {
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can not load env")
	}
	city := os.Getenv("CITY")
	state := os.Getenv("STATE")
	country := os.Getenv("COUNTRY")

	coordinates := getCoordinates(city, state, country)
	weather := getWeather(coordinates)

	fmt.Println("Current Temp:", weather.Current.Temp)
	fmt.Println("Feels Like:", weather.Current.FeelsLike)
	fmt.Println("Wind Speed:", weather.Current.WindSpeed)
	fmt.Println("Wind Speed:", getCompassDirection(weather.Current.WindDirection))

}
