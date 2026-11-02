package weather_api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/http_client"
)

var ErrApiKeyIsRequired = errors.New("apiKey is required")
var ErrCityIsRequired = errors.New("city is required")
var ErrFetchingTemp = errors.New("error fetching temperature for the specified city")

type WeatherCurrentInputDTO struct {
	CurrentTemperature CurrentTemperature `json:"current"`
}

type CurrentTemperature struct {
	TempCelsius float64 `json:"temp_c"`
}

type WeatherApiUseCase struct {
	httpClient http_client.HttpClientInterface
	apiKey     string
}

func NewWeatherApiUseCase(httpClient http_client.HttpClientInterface, apiKey string) (*WeatherApiUseCase, error) {
	if apiKey == "" {
		return nil, ErrApiKeyIsRequired
	}

	return &WeatherApiUseCase{httpClient: httpClient, apiKey: apiKey}, nil
}

func (ws *WeatherApiUseCase) GetTemperatureByCity(ctx context.Context, city string) (float64, error) {
	log.Println("Fetching temperature for the city:", city)

	if city == "" {
		return 0, ErrCityIsRequired
	}

	endpoint, err := url.Parse("https://api.weatherapi.com/v1/current.json")
	if err != nil {
		return 0, err
	}

	query := endpoint.Query()
	query.Set("key", ws.apiKey)
	query.Set("q", city)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint.String(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, ErrFetchingTemp
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var weather WeatherCurrentInputDTO
	if err := json.Unmarshal(body, &weather); err != nil {
		return 0, err
	}

	log.Printf("Temperature found for the city %s: %v\n", city, weather.CurrentTemperature.TempCelsius)
	return weather.CurrentTemperature.TempCelsius, nil
}
