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
		log.Printf("Arquivo .env não encontrado, usando variáveis de ambiente. Erro: %v", err)
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
	weatherApiService, err := weather_api.NewWeatherApiService(http_client.NewHttpClient(), apiKey)
	if err != nil {
		log.Println("Error ao criar serviço de clima:", err)
		panic(err)
	}

	return
}
