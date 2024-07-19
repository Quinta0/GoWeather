package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Struct to hold the weather data
type WeatherData struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	apiKey := os.Getenv("openWeatherKey")
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		return
	}

	response := struct {
		Temperature float64 `json:"temperature"`
		Condition   string  `json:"condition"`
	}{
		Temperature: weatherData.Main.Temp,
		Condition:   weatherData.Weather[0].Description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/weather", handleWeather)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
