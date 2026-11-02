package handler

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/usecase/via_cep"
	"github.com/fredericobenevides/pge-clima-por-cep/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWeatherHandlerValido(t *testing.T) {
	mockViaUseCase := &mocks.MockViaUseCase{}
	mockWeatherApiUseCase := &mocks.MockWeatherApiUseCase{}

	var cep = "01001000"

	mockViaUseCase.On("FindCityByCep", mock.Anything, cep).Return("São Paulo", nil)
	mockWeatherApiUseCase.On("GetTemperatureByCity", mock.Anything, "São Paulo").Return(16.4, nil)

	weatherHandler := NewWeatherHandler(
		mockViaUseCase,
		mockWeatherApiUseCase,
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", fmt.Sprintf("/weather/%s", cep), nil)
	request.SetPathValue("cep", cep)

	weatherHandler.GetWeatherByCep(recorder, request)

	assert.Equal(t, 200, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "16.4")
}

func TestWeatherHandlerCepFormatoInvalido(t *testing.T) {
	mockViaUseCase := &mocks.MockViaUseCase{}
	mockWeatherApiUseCase := &mocks.MockWeatherApiUseCase{}

	var cep = "01001-000"

	weatherHandler := NewWeatherHandler(
		mockViaUseCase,
		mockWeatherApiUseCase,
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", fmt.Sprintf("/weather/%s", cep), nil)

	weatherHandler.GetWeatherByCep(recorder, request)

	assert.Equal(t, 422, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid zipcode")
}

func TestWeatherHandlerCepNaoEncontrado(t *testing.T) {
	mockViaUseCase := &mocks.MockViaUseCase{}
	mockWeatherApiUseCase := &mocks.MockWeatherApiUseCase{}

	var cep = "11111111"

	mockViaUseCase.On("FindCityByCep", mock.Anything, cep).Return("", via_cep.ErrCepNotFound)

	weatherHandler := NewWeatherHandler(
		mockViaUseCase,
		mockWeatherApiUseCase,
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", fmt.Sprintf("/weather/%s", cep), nil)
	request.SetPathValue("cep", cep)

	weatherHandler.GetWeatherByCep(recorder, request)

	assert.Equal(t, 404, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "can not find zipcode")
}
