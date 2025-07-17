# CannySprinkler
This project now uses the OpenWeather API v3.0.

Check if sprinkling lawn is necessary.

Checks, if it rains today, rained yesterday or will rain tomorrow for a given latitude and longitude and tells you if your lawn should be sprinkled.
If you have values for soil moisture you can add them too.
You can also get the approx. sprinkle time.
The algorithm is designed for most economical pumping. Supposing you use something like this: https://www.gardena.com/de/produkte/bewasserung/pumpen/regenfasspumpe-4000-1/967974701/

You need an API Key for OpenWeather API v3.0. You can sign up at https://home.openweathermap.org/users/sign_up and subscribe to the appropriate plan that includes the One Call API 3.0. Make sure your API key has access to the One Call API 3.0.

The OpenWeather API can be used for free if you stay within their usage limits (typically 1,000 API calls per day for the free tier). This should be sufficient for most personal use cases. For more intensive usage, you may need to upgrade to a paid plan.

This project has been ported to Go from the original PHP script and API.

## Usage:

### Web API Usage

The project provides a web API with the following endpoints:

#### /sprinkleNow/:lat/:lon
Checks if it rains today, rained yesterday or will rain tomorrow for a given latitude and longitude and tells you if your lawn should be sprinkled. If you have values for soil moisture you can add them too.

Parameters:
- `lat`: Latitude of location in decimal
- `lon`: Longitude of location in decimal
- `soilm`: (Optional query parameter) Soil moisture in percent

Example: `http://localhost:8080/sprinkleNow/52.463/13.469?soilm=40`

Response:
```json
{
  "sprinkle": true
}
```

#### /sprinkleTime/:lat/:lon/:barrelv/:pumpo
Get the calculated sprinkle time in seconds. The calculation is based on the tank or barrel volume available, the pump power and the days until the next rain.

Parameters:
- `lat`: Latitude of location in decimal
- `lon`: Longitude of location in decimal
- `barrelv`: Volume of the rain tank barrel in l
- `pumpo`: Pump power of the rainwater pump in l/h
- `soilm`: (Optional query parameter) Soil moisture in percent

Example: `http://localhost:8080/sprinkleTime/52.463/13.469/300/4000?soilm=40`

Response:
```json
{
  "sprinkleTime": 135.0
}
```

### Installation

1. Make sure you have Go installed (version 1.16 or later recommended)
2. Clone this repository
3. Install dependencies:
   ```
   go mod init cannysprinkler
   go get github.com/gin-gonic/gin
   ```

### Library Usage

Example:

52.463, 13.469: Latitude and Longitude for Olching, Germany

300 l Barrel Volume

4000 l/h Pump Output

Null Soil Moisture

```go
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the API_KEY environment variable with your OpenWeather API key.")
	}

	barrelVolume := 300.0
	pumpOutput := 4000.0

	// Create a new sprinkler instance
	// In your actual code, you would import this package
	sprinkler := NewOpenWeatherSprinkler(apiKey, 52.463, 13.469, &barrelVolume, &pumpOutput, nil)

	// Check if sprinkling is needed
	fmt.Printf("SprinkleNow: %v\n", sprinkler.SprinkleNow())

	// Get the recommended sprinkle time in seconds
	fmt.Printf("SprinkleTime: %v seconds\n", sprinkler.GetSprinkleTime())
}
```

### API Setup

To set up the API, you need to:

1. Set the `API_KEY` environment variable with your OpenWeather API key
2. Run the application:
   ```
   go run *.go
   ```
   or build and run:
   ```
   go build -o cannysprinkler
   ./cannysprinkler
   ```
3. The API will be available at `http://localhost:8080`

### Docker Setup

You can run the application using Docker in two ways:

#### Option 1: Build locally

1. Build the Docker image:
   ```
   docker build -t cannysprinkler .
   ```

2. Run the container with your API key:
   ```
   docker run -p 8080:8080 -e API_KEY=your_api_key_here cannysprinkler
   ```

#### Option 2: Use pre-built image from GitHub Container Registry

1. Pull the image:
   ```
   docker pull ghcr.io/OWNER/cannysprinkler:latest
   ```
   (Replace "OWNER" with the GitHub username or organization that owns this repository)

2. Run the container with your API key:
   ```
   docker run -p 8080:8080 -e API_KEY=your_api_key_here ghcr.io/OWNER/cannysprinkler:latest
   ```

3. The API will be available at `http://localhost:8080`

### GitHub Actions

This project uses GitHub Actions to automatically build and publish Docker images to GitHub Container Registry when a new tag is created. The workflow:

1. Triggers when a tag starting with 'v' is pushed (e.g., v1.0.0)
2. Builds the Docker image using the Dockerfile
3. Pushes the image to GitHub Container Registry with the following tags:
   - Full semantic version (e.g., v1.2.3)
   - Major.minor version (e.g., 1.2)
   - 'latest' tag (if on the default branch)

To create a new release:

1. Create and push a new tag:
   ```
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The GitHub Actions workflow will automatically build and publish the Docker image
