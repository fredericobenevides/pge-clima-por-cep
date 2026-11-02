package weather_api

import (
	"context"
)

type WeatherApiUseCaseInterface interface {
	GetTemperatureByCity(ctx context.Context, city string) (float64, error)
}
