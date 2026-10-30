package handler

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/internal_error"
	"github.com/fredericobenevides/pge-clima-por-cep/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClimaHandlerValido(t *testing.T) {
	mockViaService := &mocks.MockViaService{}
	mockWeatherApi := &mocks.MockWeatherApi{}

	var cep = "01001000"

	mockViaService.On("CarregarDadosLocalidadePorCep", mock.Anything, cep).Return("São Paulo", nil)
	mockWeatherApi.On("CarregarDadosClimaPelaLocalidade", mock.Anything, "São Paulo").Return(16.4, nil)

	climaHandler := NewClimaHandler(
		mockViaService,
		mockWeatherApi,
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", fmt.Sprintf("/clima/%s", cep), nil)
	request.SetPathValue("cep", cep)

	climaHandler.BuscarClima(recorder, request)

	assert.Equal(t, 200, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "16.4")
}

func TestClimaHandlerCepFormatoInvalido(t *testing.T) {
	mockViaService := &mocks.MockViaService{}
	mockWeatherApi := &mocks.MockWeatherApi{}

	var cep = "01001-000"

	climaHandler := NewClimaHandler(
		mockViaService,
		mockWeatherApi,
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", fmt.Sprintf("/clima/%s", cep), nil)

	climaHandler.BuscarClima(recorder, request)

	assert.Equal(t, 422, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "CEP inválido. Cep deve ter 8 caracteres no formato numérico")
}

func TestClimaHandlerCepNaoEncontrado(t *testing.T) {
	mockViaService := &mocks.MockViaService{}
	mockWeatherApi := &mocks.MockWeatherApi{}

	var cep = "11111111"

	mockViaService.On("CarregarDadosLocalidadePorCep", mock.Anything, cep).Return("", internal_error.NewNotFoundError("CEP não encontrado"))

	climaHandler := NewClimaHandler(
		mockViaService,
		mockWeatherApi,
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", fmt.Sprintf("/clima/%s", cep), nil)
	request.SetPathValue("cep", cep)

	climaHandler.BuscarClima(recorder, request)

	assert.Equal(t, 404, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Não foi possível encontrar a localidade para o CEP informado")
}
