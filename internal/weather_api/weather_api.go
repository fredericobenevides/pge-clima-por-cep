package weather_api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/internal_error"
)

var ErrApiKeyObrigatoria = errors.New("apiKey deve ser fornecida")
var ErrCidadeObrigatoria = errors.New("cidade deve ser fornecida")

type WeatherCurrent struct {
	CurrentTemperature CurrentTemperature `json:"current"`
}

type CurrentTemperature struct {
	TempCelsius float64 `json:"temp_c"`
}

type WeatherApiService struct {
	apiKey string
}

func NewWeatherApiService(apiKey string) (*WeatherApiService, *internal_error.InternalError) {
	if apiKey == "" {
		return nil, internal_error.NewBadRequestError(ErrApiKeyObrigatoria.Error())
	}

	return &WeatherApiService{apiKey: apiKey}, nil
}

func (ws *WeatherApiService) CarregarDadosClimaPelaLocalidade(ctx context.Context, localidade string) (float64, *internal_error.InternalError) {
	log.Println("Buscando clima para a seguinte localidade:", localidade)

	if localidade == "" {
		return 0, internal_error.NewBadRequestError(ErrCidadeObrigatoria.Error())
	}

	endpoint, err := url.Parse("https://api.weatherapi.com/v1/current.json")
	if err != nil {
		return 0, internal_error.NewInternalError(err.Error())
	}

	query := endpoint.Query()
	query.Set("key", ws.apiKey)
	query.Set("q", localidade)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint.String(), nil)
	if err != nil {
		return 0, internal_error.NewInternalError(err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, internal_error.NewInternalError(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, internal_error.NewInternalError(err.Error())
	}

	var weather WeatherCurrent
	if err := json.Unmarshal(body, &weather); err != nil {
		return 0, internal_error.NewInternalError(err.Error())
	}

	log.Printf("Clima encontrado para a localidade %s: %v\n", localidade, weather.CurrentTemperature.TempCelsius)
	return weather.CurrentTemperature.TempCelsius, nil
}
