package weather_api

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

func TestNewWeatherApiServiceSemApiKey(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}

	weatherApiService, err := NewWeatherApiService(httpClientMock, "")
	assert.Nil(t, weatherApiService)
	assert.NotNil(t, "apiKey deve ser fornecida", err)
}

func TestCarregarDadosClimaPelaLocalidadeValido(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"current":{"temp_c": 14.6}}`)),
	}, nil)

	weatherApiService, err := NewWeatherApiService(httpClientMock, "apiKey")
	assert.Nil(t, err)

	weather, err := weatherApiService.CarregarDadosClimaPelaLocalidade(context.Background(), "São Paulo")
	assert.Nil(t, err)
	assert.Equal(t, 14.6, weather)
}

func TestCarregarDadosClimaPelaLocalidadeVazio(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}

	weatherApiService, err := NewWeatherApiService(httpClientMock, "apiKey")
	assert.Nil(t, err)

	_, err = weatherApiService.CarregarDadosClimaPelaLocalidade(context.Background(), "")
	assert.NotNil(t, "cidade deve ser fornecida", err)
}

func TestCarregarDadosClimaPelaLocalidadeInvalido(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"code":1006,"message":"No matching location found."}}`)),
	}, nil)

	weatherApiService, err := NewWeatherApiService(httpClientMock, "apiKey")
	assert.Nil(t, err)

	_, err = weatherApiService.CarregarDadosClimaPelaLocalidade(context.Background(), "Cidade inválida")
	assert.NotNil(t, "Erro ao buscar clima para a localidade: Cidade inválida", err)
}
