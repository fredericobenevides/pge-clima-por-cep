package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockWeatherApiUseCase struct {
	mock.Mock
}

func (m *MockWeatherApiUseCase) GetTemperatureByCity(ctx context.Context, city string) (float64, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(float64), args.Error(1)
}
