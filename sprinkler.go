package main

// CannySprinkler is the interface that defines the functionality of a sprinkler system
type CannySprinkler interface {
	// GetLatLon returns the latitude and longitude of the location
	GetLatLon() (float64, float64)

	// GetBarrelVolume returns the volume of the water barrel
	GetBarrelVolume() *float64

	// GetPumpOutput returns the output of the pump
	GetPumpOutput() *float64

	// GetSoilMoisture returns the soil moisture
	GetSoilMoisture() *float64

	// RainsToday checks if it's raining today
	RainsToday(weatherData map[string]interface{}) bool

	// RainsTomorrow checks if it will rain tomorrow
	RainsTomorrow(weatherData map[string]interface{}) bool

	// RainedYesterday checks if it rained yesterday
	RainedYesterday(weatherData map[string]interface{}) bool

	// DaysToNextRain returns the number of days until the next rain
	DaysToNextRain(weatherData map[string]interface{}) int

	// SprinkleNow determines if sprinkling should be done now
	SprinkleNow() bool

	// GetSprinkleTime calculates the time to sprinkle in seconds
	GetSprinkleTime() float64
}
