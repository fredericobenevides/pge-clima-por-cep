package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockViaUseCase struct {
	mock.Mock
}

func (m *MockViaUseCase) FindCityByCep(ctx context.Context, cep string) (string, error) {
	args := m.Called(ctx, cep)
	return args.String(0), args.Error(1)
}
