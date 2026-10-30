package mocks

import (
	"context"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/internal_error"
	"github.com/stretchr/testify/mock"
)

type MockWeatherApi struct {
	mock.Mock
}

func (m *MockWeatherApi) CarregarDadosClimaPelaLocalidade(ctx context.Context, localidade string) (float64, *internal_error.InternalError) {
	args := m.Called(ctx, localidade)

	var err *internal_error.InternalError

	if args.Get(1) != nil {
		err = args.Get(1).(*internal_error.InternalError)
	}

	return args.Get(0).(float64), err
}
