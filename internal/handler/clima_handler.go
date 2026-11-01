package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/via_cep"
	"github.com/fredericobenevides/pge-clima-por-cep/internal/weather_api"
	"github.com/fredericobenevides/pge-clima-por-cep/pkg/converter"
)

type ClimaResponse struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type ClimaHandler struct {
	viaCepService  via_cep.ViaCepServiceInterface
	weatherService weather_api.WeatherApiServiceInterface
}

func NewClimaHandler(viaCepService via_cep.ViaCepServiceInterface, weatherService weather_api.WeatherApiServiceInterface) *ClimaHandler {
	return &ClimaHandler{
		viaCepService:  viaCepService,
		weatherService: weatherService,
	}
}

func (ch *ClimaHandler) BuscarClima(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	cep := r.PathValue("cep")

	if regexp.MustCompile(`^\d{8}$`).MatchString(cep) == false {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	cepData, err := ch.viaCepService.CarregarDadosLocalidadePorCep(ctx, cep)
	if err != nil {
		if errors.Is(err, via_cep.ErrCepNaoEncontrado) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Println("Erro ao buscar localidade:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temperatura, err := ch.weatherService.CarregarDadosClimaPelaLocalidade(ctx, cepData)
	if err != nil {
		if errors.Is(err, weather_api.ErrCidadeObrigatoria) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("Erro ao buscar dados do clima:", err)
		http.Error(w, "Não foi possível buscar os dados do clima para a localidade informada", http.StatusInternalServerError)
		return
	}

	climaResponse := ClimaResponse{
		Celsius:    temperatura,
		Fahrenheit: converter.CelsiusToFahrenheit(temperatura),
		Kelvin:     converter.CelsiusToKelvin(temperatura),
	}

	json.NewEncoder(w).Encode(climaResponse)

}
