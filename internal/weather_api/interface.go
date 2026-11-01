package weather_api

import (
	"context"
)

type WeatherApiServiceInterface interface {
	CarregarDadosClimaPelaLocalidade(ctx context.Context, localidade string) (float64, error)
}
