package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockViaService struct {
	mock.Mock
}

func (m *MockViaService) CarregarDadosLocalidadePorCep(ctx context.Context, cep string) (string, error) {
	args := m.Called(ctx, cep)
	return args.String(0), args.Error(1)
}
