package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kyle-kettler/go-running/database"
	"github.com/kyle-kettler/go-running/services"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Error initializing datadase:", err)
	}
	defer db.Close()

	savedAddress, err := database.GetLocation(db)
	if err != nil {
		log.Fatal("Error getting saved address: ", err)
	}

	address := flag.String("address", savedAddress, "Address")

	flag.Parse()

	if flag.NFlag() > 0 {
		if err := database.SaveLocation(db, *address); err != nil {
			log.Fatal("Error saving location:", err)
		}
	}

	if *address == "" {
		log.Fatal("Please provide an address using the -address flag")
	}

	coordinates := services.GetCoordinates(*address)
	timezone := services.GetTimezone(coordinates)
	weather := services.GetCurrentWeather(coordinates)
	forcast := services.GetForecast(coordinates, timezone)

	fmt.Println("Timezone: ", timezone)
	fmt.Printf("Current Temp: %.1f°\n", weather.Current.Temp)
	fmt.Printf("Feels Like: %.1f°\n", weather.Current.FeelsLike)
	fmt.Printf("Wind Speed: %.2f mph\n", weather.Current.WindSpeed)
	fmt.Println("Wind Direction:", services.GetCompassDirection(weather.Current.WindDirection))
	fmt.Println("===============")
	fmt.Println("Hourly Forecast:")
	for i, temp := range forcast.Hourly.Temperature {
		if i < 24 { // Show next 24 hours
			hour := forcast.Hourly.Time[i]
			fmt.Printf("%s: %.1f°\n", hour, temp)
		}
	}
}
