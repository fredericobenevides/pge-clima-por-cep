package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/handler"
	"github.com/fredericobenevides/pge-clima-por-cep/internal/http_client"
	"github.com/fredericobenevides/pge-clima-por-cep/internal/usecase/via_cep"
	"github.com/fredericobenevides/pge-clima-por-cep/internal/usecase/weather_api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("The file .env was not found, using environment variables. Error: %v", err)
	}

	viaCepUseCase, weatherApiUseCase := initDependencies()
	http.HandleFunc("/weather/{cep}", handler.NewWeatherHandler(viaCepUseCase, weatherApiUseCase).GetWeatherByCep)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func initDependencies() (viaCepUseCase *via_cep.ViaCepUseCase, weatherApiUseCase *weather_api.WeatherApiUseCase) {
	viaCepUseCase = via_cep.NewViaCepUseCase(http_client.NewHttpClient())

	apiKey := os.Getenv("WEATHER_API_KEY")
	weatherApiUseCase, err := weather_api.NewWeatherApiUseCase(http_client.NewHttpClient(), apiKey)
	if err != nil {
		log.Println("Error creating weatherApiUseCase:", err)
		panic(err)
	}

	return viaCepUseCase, weatherApiUseCase
}
