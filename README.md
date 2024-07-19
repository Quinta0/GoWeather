# GoWeather

This is a simple weather application built with Go, featuring a web-based GUI. The app allows users to enter a city name and view the current temperature and weather conditions for that location.

## Table of Contents

1. [Features](#features)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Code Structure](#code-structure)
6. [How It Works](#how-it-works)
7. [Customization](#customization)
8. [Contributing](#contributing)
9. [License](#license)

## Features

- Web-based user interface
- Real-time weather data fetching
- Display of temperature and weather conditions for a specified city

## Prerequisites

- Go (version 1.16 or later)
- OpenWeatherMap API key

## Installation

1. Clone this repository or download the source code.
2. Navigate to the project directory.
3. Replace `"YOUR_API_KEY_HERE"` in `main.go` with your actual OpenWeatherMap API key.

## Usage

1. Open a terminal and navigate to the project directory.
2. Run the following command to start the server:
    ```bash
    go run main.go
    ```
3. Open a web browser and go to `http://localhost:8080`.
4. Enter a city name in the input field and click "Get Weather" to see the current weather information.

## Code Structure

The project consists of two main files:

1. `main.go`: The Go server that handles API requests and serves the web interface.
2. `index.html`: The HTML file that provides the user interface.

### main.go

This file contains the Go code for the server. Here's a breakdown of its main components:

#### WeatherData struct

```go
type WeatherData struct {
 Main struct {
     Temp float64 `json:"temp"`
 } `json:"main"`
 Weather []struct {
     Description string `json:"description"`
 } `json:"weather"`
}
```
This struct defines the structure for parsing the JSON response from the OpenWeatherMap API. It extracts the temperature and weather description from the API response.

#### main() function
````go
func main() {
    http.HandleFunc("/", serveHome)
    http.HandleFunc("/weather", handleWeather)

    fmt.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
````
The main() function sets up the HTTP server:

It registers two route handlers: one for the home page ("/") and one for weather data requests ("/weather").
It starts the server on port 8080 and logs any fatal errors.

#### serveHome() function
````go
func serveHome(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}
````
This function serves the index.html file when a request is made to the root ("/") of the server.
#### handleWeather() function
````go
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
````
The handleWeather() function processes weather data requests:

1) It extracts the city name from the query parameters.
2) It constructs the OpenWeatherMap API URL with the city name and API key.
3) It sends a GET request to the API and reads the response.
4) It parses the JSON response into the WeatherData struct.
5) It creates a new struct with the relevant weather information (temperature and condition).
6) Finally, it sends this data back to the client as a JSON response.

### index.html
This file contains the HTML, CSS, and JavaScript for the user interface. Key components include:
#### HTML Structure
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <!-- Head content -->
</head>
<body>
    <div class="container">
        <h1>Weather App</h1>
        <input type="text" id="cityInput" placeholder="Enter city name">
        <button onclick="getWeather()">Get Weather</button>
        <div id="weatherInfo"></div>
    </div>
    <script>
        // JavaScript code
    </script>
</body>
</html>
```

This structure creates a simple interface with an input field for the city name, a button to fetch weather data, and a div to display the weather information.

#### css Styles
```css
<style>
    body {
        font-family: Arial, sans-serif;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        margin: 0;
        background-color: #f0f0f0;
    }
    .container {
        background-color: white;
        padding: 20px;
        border-radius: 5px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
    input, button {
        margin: 10px 0;
        padding: 5px;
    }
</style>
```

These styles create a centered layout with a white container for the weather app interface.

#### JavaScript Function
```javascript
function getWeather() {
    const city = document.getElementById('cityInput').value;
    fetch(`/weather?city=${encodeURIComponent(city)}`)
        .then(response => response.json())
        .then(data => {
            const weatherInfo = document.getElementById('weatherInfo');
            weatherInfo.innerHTML = `
                <h2>Weather in ${city}</h2>
                <p>Temperature: ${data.temperature}°C</p>
                <p>Condition: ${data.condition}</p>
            `;
        })
        .catch(error => {
            console.error('Error:', error);
            const weatherInfo = document.getElementById('weatherInfo');
            weatherInfo.innerHTML = '<p>Error fetching weather data</p>';
        });
}
```

This function is called when the user clicks the "Get Weather" button:

1. It retrieves the city name from the input field.
2. It sends a fetch request to the `/weather` endpoint with the city name.
3. When it receives the response, it updates the `weatherInfo` div with the temperature and weather condition.
4. If there's an error, it displays an error message.

## How It Works

1. The user enters a city name and clicks "Get Weather".
2. The JavaScript function `getWeather()` sends a request to the `/weather` endpoint of the Go server.
3. The Go server receives the request and fetches weather data from the OpenWeatherMap API.
4. The server processes the API response and sends back the relevant weather information as JSON.
5. The JavaScript in the browser receives this data and updates the HTML to display the weather information.

## Customization

You can customize this app by:

- Modifying the HTML and CSS in `index.html` to change the appearance of the app.
- Adding more weather data fields in both the Go code and HTML to display additional information.
- Implementing error handling for cases like invalid city names or API errors.

## Contributing

Contributions to improve the app are welcome. Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).