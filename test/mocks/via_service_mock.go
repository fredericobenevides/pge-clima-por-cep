package mocks

import (
	"context"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/internal_error"
	"github.com/stretchr/testify/mock"
)

type MockViaService struct {
	mock.Mock
}

func (m *MockViaService) CarregarDadosLocalidadePorCep(ctx context.Context, cep string) (string, *internal_error.InternalError) {
	args := m.Called(ctx, cep)

	var err *internal_error.InternalError

	if args.Get(1) != nil {
		err = args.Get(1).(*internal_error.InternalError)
	}

	return args.String(0), err
}
