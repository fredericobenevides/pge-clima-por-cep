package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/usecase/via_cep"
	"github.com/fredericobenevides/pge-clima-por-cep/internal/usecase/weather_api"
	"github.com/fredericobenevides/pge-clima-por-cep/pkg/converter"
)

type WeatherDTOResponse struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type WeatherHandler struct {
	viaCepUseCase  via_cep.ViaCepUseCaseInterface
	weatherUseCase weather_api.WeatherApiUseCaseInterface
}

func NewWeatherHandler(viaCepUseCase via_cep.ViaCepUseCaseInterface, weatherUseCase weather_api.WeatherApiUseCaseInterface) *WeatherHandler {
	return &WeatherHandler{
		viaCepUseCase,
		weatherUseCase,
	}
}

func (ch *WeatherHandler) GetWeatherByCep(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	cep := r.PathValue("cep")

	if regexp.MustCompile(`^\d{8}$`).MatchString(cep) == false {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := ch.viaCepUseCase.FindCityByCep(ctx, cep)
	if err != nil {
		if errors.Is(err, via_cep.ErrCepNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Println("Error fetching city for the specified zipcode:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp, err := ch.weatherUseCase.GetTemperatureByCity(ctx, city)
	if err != nil {
		if errors.Is(err, weather_api.ErrCityIsRequired) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("Error fetching temperature for the specified city:", err)
		http.Error(w, "error fetching temperature for the specified city", http.StatusInternalServerError)
		return
	}

	dtoResp := WeatherDTOResponse{
		Celsius:    temp,
		Fahrenheit: converter.CelsiusToFahrenheit(temp),
		Kelvin:     converter.CelsiusToKelvin(temp),
	}

	json.NewEncoder(w).Encode(dtoResp)

}
