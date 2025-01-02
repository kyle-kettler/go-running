package models

type Coordinates struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type Timezone struct {
	Zone string `json:"timeZone"`
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

type Hourly struct {
	Time                []string  `json:"time"`
	Temperature         []float64 `json:"temperature_2m"`
	FeelsLike           []float64 `json:"apparent_temperature"`
	PrecipitationChance []int     `json:"precipitation_probability"`
	Rain                []int     `json:"rain"`
	Showers             []int     `json:"showers"`
	Snowfall            []int     `json:"snowfall"`
	SnowDepth           []int     `json:"snow_depth"`
	CloudCover          []int     `json:"cloud_cover"`
	WindSpeed           []float64 `json:"wind_speed_10m"`
	WindDirection       []int     `json:"wind_direction_10m"`
	SoilTemperature     []float64 `json:"soil_temperature_0cm"`
}

type OneDayForecast struct {
	Hourly Hourly `json:"hourly"`
}

type Weather struct {
	Current CurrentWeather
}
