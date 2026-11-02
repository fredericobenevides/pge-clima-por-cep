package via_cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/http_client"
)

var ErrCepIsRequired = errors.New("zipcode is required")
var ErrCepNotFound = errors.New("can not find zipcode")
var ErrFetchingCity = errors.New("error fetching city for the specified zipcode")

type ViaCep struct {
	City string `json:"localidade"`
}

type ViaCepUseCase struct {
	httpClient http_client.HttpClientInterface
}

func NewViaCepUseCase(httpClient http_client.HttpClientInterface) *ViaCepUseCase {
	return &ViaCepUseCase{
		httpClient: httpClient,
	}
}

func (vc *ViaCepUseCase) FindCityByCep(ctx context.Context, cep string) (string, error) {
	log.Println("Fetching city for the zipcode:", cep)

	if cep == "" {
		return "", ErrCepIsRequired
	}

	var url = fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrCepNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrFetchingCity
	}

	var viaCep ViaCep

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &viaCep); err != nil {
		return "", err
	}

	if viaCep.City == "" {
		return "", ErrCepNotFound
	}

	log.Println("City found:", viaCep.City)
	return viaCep.City, nil
}
