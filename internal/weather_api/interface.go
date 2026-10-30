package weather_api

import (
	"context"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/internal_error"
)

type WeatherApiServiceInterface interface {
	CarregarDadosClimaPelaLocalidade(ctx context.Context, localidade string) (float64, *internal_error.InternalError)
}
