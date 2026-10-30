package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/handler"
	"github.com/fredericobenevides/pge-clima-por-cep/internal/http_client"
	"github.com/fredericobenevides/pge-clima-por-cep/internal/via_cep"
	"github.com/fredericobenevides/pge-clima-por-cep/internal/weather_api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar variáveis de ambiente: %v", err)
	}

	viaCepService, weatherService := initDependencies()
	http.HandleFunc("/clima/{cep}", handler.NewClimaHandler(viaCepService, weatherService).BuscarClima)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func initDependencies() (viaCepService *via_cep.ViaCepService, weatherApiService *weather_api.WeatherApiService) {
	viaCepService = via_cep.NewViaCepService(http_client.NewHttpClient())

	apiKey := os.Getenv("WEATHER_API_KEY")
	weatherApiService, ierr := weather_api.NewWeatherApiService(http_client.NewHttpClient(), apiKey)
	if ierr != nil {
		log.Println("Error ao criar serviço de clima:", ierr.Message)
		panic(ierr)
	}

	return
}
