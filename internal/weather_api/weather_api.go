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

var ErrApiKeyObrigatoria = errors.New("apiKey deve ser fornecida")
var ErrCidadeObrigatoria = errors.New("cidade deve ser fornecida")
var ErrBuscaClima = errors.New("Erro ao buscar clima para a localidade")

type WeatherCurrent struct {
	CurrentTemperature CurrentTemperature `json:"current"`
}

type CurrentTemperature struct {
	TempCelsius float64 `json:"temp_c"`
}

type WeatherApiService struct {
	httpClient http_client.HttpClientInterface
	apiKey     string
}

func NewWeatherApiService(httpClient http_client.HttpClientInterface, apiKey string) (*WeatherApiService, error) {
	if apiKey == "" {
		return nil, ErrApiKeyObrigatoria
	}

	return &WeatherApiService{httpClient: httpClient, apiKey: apiKey}, nil
}

func (ws *WeatherApiService) CarregarDadosClimaPelaLocalidade(ctx context.Context, localidade string) (float64, error) {
	log.Println("Buscando clima para a seguinte localidade:", localidade)

	if localidade == "" {
		return 0, ErrCidadeObrigatoria
	}

	endpoint, err := url.Parse("https://api.weatherapi.com/v1/current.json")
	if err != nil {
		return 0, err
	}

	query := endpoint.Query()
	query.Set("key", ws.apiKey)
	query.Set("q", localidade)
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
		return 0, ErrBuscaClima
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var weather WeatherCurrent
	if err := json.Unmarshal(body, &weather); err != nil {
		return 0, err
	}

	log.Printf("Clima encontrado para a localidade %s: %v\n", localidade, weather.CurrentTemperature.TempCelsius)
	return weather.CurrentTemperature.TempCelsius, nil
}
