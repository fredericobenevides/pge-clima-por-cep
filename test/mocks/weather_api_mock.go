package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockWeatherApi struct {
	mock.Mock
}

func (m *MockWeatherApi) CarregarDadosClimaPelaLocalidade(ctx context.Context, localidade string) (float64, error) {
	args := m.Called(ctx, localidade)
	return args.Get(0).(float64), args.Error(1)
}
