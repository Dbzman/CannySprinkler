package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// OpenWeatherSprinkler implements the CannySprinkler interface
type OpenWeatherSprinkler struct {
	apiKey        string
	lat           float64
	lon           float64
	barrelVolume  *float64
	pumpOutput    *float64
	soilMoisture  *float64
	params        map[string]float64
	nowForecast   map[string]interface{}
	histoForecast map[string]interface{}
}

// NewOpenWeatherSprinkler creates a new OpenWeatherSprinkler instance
func NewOpenWeatherSprinkler(apiKey string, lat, lon float64, barrelVolume, pumpOutput, soilMoisture *float64) *OpenWeatherSprinkler {
	sprinkler := &OpenWeatherSprinkler{
		apiKey:       apiKey,
		lat:          lat,
		lon:          lon,
		barrelVolume: barrelVolume,
		pumpOutput:   pumpOutput,
		soilMoisture: soilMoisture,
		params: map[string]float64{
			"soil_moisture_lower": 30,
			"soil_moisture_upper": 80,
			"barrel_buffer":       0.5,
		},
	}

	sprinkler.setNowForecast()
	sprinkler.setHistoForecast()

	return sprinkler
}

// urlNowForecast returns the URL for the current and forecast weather data
func (s *OpenWeatherSprinkler) urlNowForecast() string {
	return fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&exclude=minutely,hourly,alerts&appid=%s",
		s.lat, s.lon, s.apiKey)
}

// urlHisto returns the URL for the historical weather data
func (s *OpenWeatherSprinkler) urlHisto() string {
	yesterday := time.Now().AddDate(0, 0, -1).Unix()
	return fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall/timemachine?lat=%f&lon=%f&dt=%d&appid=%s",
		s.lat, s.lon, yesterday, s.apiKey)
}

// setNowForecast fetches and sets the current and forecast weather data
func (s *OpenWeatherSprinkler) setNowForecast() {
	resp, err := http.Get(s.urlNowForecast())
	if err != nil {
		log.Printf("Error fetching forecast data: %v", err)
		s.nowForecast = map[string]interface{}{
			"current": map[string]interface{}{},
			"daily":   []interface{}{},
		}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading forecast data: %v", err)
		s.nowForecast = map[string]interface{}{
			"current": map[string]interface{}{},
			"daily":   []interface{}{},
		}
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error parsing forecast data: %v", err)
		s.nowForecast = map[string]interface{}{
			"current": map[string]interface{}{},
			"daily":   []interface{}{},
		}
		return
	}

	s.nowForecast = data
}

// setHistoForecast fetches and sets the historical weather data
func (s *OpenWeatherSprinkler) setHistoForecast() {
	resp, err := http.Get(s.urlHisto())
	if err != nil {
		log.Printf("Error fetching historical data: %v", err)
		s.histoForecast = map[string]interface{}{
			"hourly": []interface{}{},
		}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading historical data: %v", err)
		s.histoForecast = map[string]interface{}{
			"hourly": []interface{}{},
		}
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error parsing historical data: %v", err)
		s.histoForecast = map[string]interface{}{
			"hourly": []interface{}{},
		}
		return
	}

	s.histoForecast = data
}

// GetNowForecast returns the current and forecast weather data
func (s *OpenWeatherSprinkler) GetNowForecast() map[string]interface{} {
	if s.nowForecast == nil {
		s.setNowForecast()
	}
	return s.nowForecast
}

// GetHistoForecast returns the historical weather data
func (s *OpenWeatherSprinkler) GetHistoForecast() map[string]interface{} {
	if s.histoForecast == nil {
		s.setHistoForecast()
	}
	return s.histoForecast
}

// GetLatLon returns the latitude and longitude
func (s *OpenWeatherSprinkler) GetLatLon() (float64, float64) {
	return s.lat, s.lon
}

// GetBarrelVolume returns the barrel volume
func (s *OpenWeatherSprinkler) GetBarrelVolume() *float64 {
	return s.barrelVolume
}

// GetPumpOutput returns the pump output
func (s *OpenWeatherSprinkler) GetPumpOutput() *float64 {
	return s.pumpOutput
}

// GetSoilMoisture returns the soil moisture
func (s *OpenWeatherSprinkler) GetSoilMoisture() *float64 {
	return s.soilMoisture
}

// RainsToday checks if it's raining today
func (s *OpenWeatherSprinkler) RainsToday(weatherData map[string]interface{}) bool {
	if weatherData == nil {
		weatherData = s.GetNowForecast()
	}

	current, ok := weatherData["current"].(map[string]interface{})
	if !ok {
		return false
	}

	// Check if rain exists in current.rain.1h or current.rain
	if rain, ok := current["rain"].(map[string]interface{}); ok {
		if _, ok := rain["1h"]; ok {
			return true
		}
	}
	if _, ok := current["rain"].(float64); ok {
		return true
	}

	return false
}

// RainsTomorrow checks if it will rain tomorrow
func (s *OpenWeatherSprinkler) RainsTomorrow(weatherData map[string]interface{}) bool {
	if weatherData == nil {
		weatherData = s.GetNowForecast()
	}

	daily, ok := weatherData["daily"].([]interface{})
	if !ok || len(daily) == 0 {
		return false
	}

	tomorrow, ok := daily[0].(map[string]interface{})
	if !ok {
		return false
	}

	// Check if rain exists in daily[0].rain.1h or daily[0].rain
	if rain, ok := tomorrow["rain"].(map[string]interface{}); ok {
		if _, ok := rain["1h"]; ok {
			return true
		}
	}
	if _, ok := tomorrow["rain"].(float64); ok {
		return true
	}

	return false
}

// RainedYesterday checks if it rained yesterday
func (s *OpenWeatherSprinkler) RainedYesterday(weatherData map[string]interface{}) bool {
	if weatherData == nil {
		weatherData = s.GetHistoForecast()
	}

	hourly, ok := weatherData["hourly"].([]interface{})
	if !ok {
		return false
	}

	hoursWithRain := 0
	for _, h := range hourly {
		hour, ok := h.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if rain exists in hour.rain.1h or hour.rain
		if rain, ok := hour["rain"].(map[string]interface{}); ok {
			if _, ok := rain["1h"]; ok {
				hoursWithRain++
			}
		}
		if _, ok := hour["rain"].(float64); ok {
			hoursWithRain++
		}
	}

	return hoursWithRain > 1
}

// DaysToNextRain returns the number of days until the next rain
func (s *OpenWeatherSprinkler) DaysToNextRain(weatherData map[string]interface{}) int {
	if weatherData == nil {
		weatherData = s.GetNowForecast()
	}

	daily, ok := weatherData["daily"].([]interface{})
	if !ok {
		return 7
	}

	for idx, d := range daily {
		day, ok := d.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if rain exists in day.rain.1h or day.rain
		if rain, ok := day["rain"].(map[string]interface{}); ok {
			if _, ok := rain["1h"]; ok {
				return idx + 1
			}
		}
		if _, ok := day["rain"].(float64); ok {
			return idx + 1
		}
	}

	return 7
}

// SprinkleNow determines if sprinkling should be done now
func (s *OpenWeatherSprinkler) SprinkleNow() bool {
	if s.RainsToday(nil) {
		return false
	}

	if s.soilMoisture != nil {
		if *s.soilMoisture <= s.params["soil_moisture_lower"] {
			return true
		}
		if *s.soilMoisture >= s.params["soil_moisture_upper"] {
			return false
		}
	}

	if s.RainedYesterday(nil) {
		return false
	}

	if s.RainsTomorrow(nil) {
		return false
	}

	return true
}

// GetSprinkleTime calculates the time to sprinkle in seconds
func (s *OpenWeatherSprinkler) GetSprinkleTime() float64 {
	if s.barrelVolume == nil || s.pumpOutput == nil {
		return 0
	}

	sprinkleSeconds := ((*s.barrelVolume * s.params["barrel_buffer"] / float64(s.DaysToNextRain(nil))) / *s.pumpOutput) * 3600
	return sprinkleSeconds
}
