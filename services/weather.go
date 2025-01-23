package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kyle-kettler/go-running/models"
)

const (
	weatherBaseURL = "https://api.open-meteo.com/v1/forecast"
)

type WeatherAPIConfig struct {
	CurrentParams []string
	HourlyParams  []string
	Units         map[string]string
	Timezone      string
	ForecastDays  int
}

var defaultWeatherConfig = WeatherAPIConfig{
	CurrentParams: []string{
		"temperature_2m",
		"apparent_temperature",
		"precipitation",
		"rain",
		"showers",
		"snowfall",
		"wind_speed_10m",
		"wind_direction_10m",
	},
	HourlyParams: []string{
		"temperature_2m",
		"apparent_temperature",
		"precipitation_probability",
		"rain",
		"showers",
		"snowfall",
		"snow_depth",
		"cloud_cover",
		"wind_speed_10m",
		"wind_direction_10m",
		"soil_temperature_0cm",
	},
	Units: map[string]string{
		"temperature":   "fahrenheit",
		"wind_speed":    "mph",
		"precipitation": "inch",
	},
	ForecastDays: 2,
}

func buildWeatherURL(coordinates models.Coordinates, config WeatherAPIConfig, timezone *string, forecast bool) string {
	params := url.Values{}

	params.Add("latitude", coordinates.Lat)
	params.Add("longitude", coordinates.Lon)
	if !forecast {
		params.Add("current", strings.Join(config.CurrentParams, ","))
	} else {
		params.Add("hourly", strings.Join(config.HourlyParams, ","))
	}

	if timezone != nil {
		params.Add("timezone", *timezone)
	}

	for unit, value := range config.Units {
		params.Add(unit+"_unit", value)
	}

	if forecast {
		params.Add("forecast_days", strconv.Itoa(config.ForecastDays))
	}

	return weatherBaseURL + "?" + params.Encode()
}

func GetCurrentWeather(coordinates models.Coordinates) models.Weather {
	baseUrl := buildWeatherURL(coordinates, defaultWeatherConfig, nil, false)

	res, err := http.Get(baseUrl)
	if err != nil {
		log.Fatal("error getting weather", err, baseUrl)
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

func GetForecast(coordinates models.Coordinates, timezone string) models.OneDayForecast {
	baseUrl := buildWeatherURL(coordinates, defaultWeatherConfig, &timezone, true)

	res, err := http.Get(baseUrl)
	if err != nil {
		log.Fatal("error getting weather")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("error reading weather body")
	}

	var weather models.OneDayForecast
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

func BuildForecast(forcast models.Hourly) {
	currentTime := time.Now()
	location := currentTime.Location()
	endTime := currentTime.Add(6 * time.Hour)
	count := 0

	for i, temp := range forcast.Temperature {
		forecastTime, err := time.ParseInLocation("2006-01-02T15:04", forcast.Time[i], location)
		if err != nil {
			log.Printf("Error parsing time: %v", err)
			continue
		}

		if forecastTime.After(currentTime) && forecastTime.Before(endTime) {
			formattedTime := forecastTime.Format(time.Kitchen)
			formattedDate := forecastTime.Format("01/02")
			fmt.Printf("%s %s \r\n", formattedDate, formattedTime)
			fmt.Printf("Temp: %.1f° | Feels Like: %.1f° |\n", temp, forcast.FeelsLike[i])
			fmt.Println("=======================================")
			count++
		}
	}
}
