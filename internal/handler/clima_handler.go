package handler

import (
	"context"
	"encoding/json"
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
		http.Error(w, "CEP inválido. Cep deve ter 8 caracteres no formato numérico", http.StatusUnprocessableEntity)
		return
	}

	cepData, err := ch.viaCepService.CarregarDadosLocalidadePorCep(ctx, cep)
	if err != nil {
		log.Println("Erro ao buscar localidade:", err.Message)
		w.WriteHeader(err.Status)
		w.Write([]byte("Não foi possível encontrar a localidade para o CEP informado"))
		return
	}

	temperatura, err := ch.weatherService.CarregarDadosClimaPelaLocalidade(ctx, cepData)
	if err != nil {
		log.Printf("Erro ao buscar dados do clima: %v", err)
		w.WriteHeader(err.Status)
		w.Write([]byte("Não foi possível buscar os dados do clima para a localidade informada"))
		return
	}

	climaResponse := ClimaResponse{
		Celsius:    temperatura,
		Fahrenheit: converter.CelsiusToFahrenheit(temperatura),
		Kelvin:     converter.CelsiusToKelvin(temperatura),
	}

	json.NewEncoder(w).Encode(climaResponse)

}
