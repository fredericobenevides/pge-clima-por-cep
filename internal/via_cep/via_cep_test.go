package via_cep

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/fredericobenevides/pge-clima-por-cep/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCarregarDadosLocalidadePorCepValido(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"localidade":"São Paulo"}`)),
	}, nil)

	viaCepService := NewViaCepService(httpClientMock)
	localidade, err := viaCepService.CarregarDadosLocalidadePorCep(context.Background(), "01001000")

	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", localidade)
}

func TestCarregarDadosLocalidadePorCepVazio(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}

	viaCepService := NewViaCepService(httpClientMock)
	_, err := viaCepService.CarregarDadosLocalidadePorCep(context.Background(), "")

	assert.NotNil(t, err)
}

func TestCarregarDadosLocalidadePorCepInvalido(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(`{"erro":"true"}`)),
	}, nil)

	viaCepService := NewViaCepService(httpClientMock)
	_, err := viaCepService.CarregarDadosLocalidadePorCep(context.Background(), "invalid-cep")

	assert.NotNil(t, "Erro ao buscar localidade para o CEP: invalid-cep", err)
}
