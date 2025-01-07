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
	weather := services.GetCurrentWeather(coordinates)

	fmt.Println(coordinates.Lat, coordinates.Lon)
	fmt.Println("===============")
	fmt.Printf("Current Temp: %.1f°\n", weather.Current.Temp)
	fmt.Printf("Feels Like: %.1f°\n", weather.Current.FeelsLike)
	fmt.Printf("Wind Speed: %.2f mph\n", weather.Current.WindSpeed)
	fmt.Println("Wind Direction:", services.GetCompassDirection(weather.Current.WindDirection))
}
