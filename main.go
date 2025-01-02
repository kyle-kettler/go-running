package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kyle-kettler/go-running/services"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Can not load env")
    }

    city := os.Getenv("CITY")
    state := os.Getenv("STATE")
    country := os.Getenv("COUNTRY")

    coordinates := services.GetCoordinates(city, state, country)
    weather := services.GetWeather(coordinates)

    fmt.Println(coordinates.Lat, coordinates.Lon)
    fmt.Println("Current Temp:", weather.Current.Temp)
    fmt.Println("Feels Like:", weather.Current.FeelsLike)
    fmt.Println("Wind Speed:", weather.Current.WindSpeed)
    fmt.Println("Wind Direction:", services.GetCompassDirection(weather.Current.WindDirection))
}
