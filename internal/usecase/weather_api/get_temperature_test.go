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

func TestNewWeatherApiUseCaseWithEmptyApiKey(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}

	weatherApiUseCase, err := NewWeatherApiUseCase(httpClientMock, "")
	assert.Nil(t, weatherApiUseCase)
	assert.NotNil(t, err)
	assert.Equal(t, "apiKey is required", err.Error())
}

func TestGetTemperatureByCityValid(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"current":{"temp_c": 14.6}}`)),
	}, nil)

	weatherApiUseCase, err := NewWeatherApiUseCase(httpClientMock, "apiKey")
	assert.Nil(t, err)

	weather, err := weatherApiUseCase.GetTemperatureByCity(context.Background(), "São Paulo")
	assert.Nil(t, err)
	assert.Equal(t, 14.6, weather)
}

func TestGetTemperatureByCityNotProvided(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}

	weatherApiUseCase, err := NewWeatherApiUseCase(httpClientMock, "apiKey")
	assert.Nil(t, err)

	_, err = weatherApiUseCase.GetTemperatureByCity(context.Background(), "")
	assert.NotNil(t, err)
	assert.Equal(t, "city is required", err.Error())
}

func TestGetTemperatureByCityNotFound(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"code":1006,"message":"No matching location found."}}`)),
	}, nil)

	weatherApiUseCase, err := NewWeatherApiUseCase(httpClientMock, "apiKey")
	assert.Nil(t, err)

	_, err = weatherApiUseCase.GetTemperatureByCity(context.Background(), "Cidade inválida")
	assert.NotNil(t, err)
	assert.Equal(t, "error fetching temperature for the specified city", err.Error())
}
