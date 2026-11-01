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

var ErrCepObrigatorio = errors.New("zipcode is required")
var ErrCepNaoEncontrado = errors.New("can not find zipcode")
var ErrBuscaCep = errors.New("Error fetching zipcode data")

type ViaCep struct {
	Localidade string `json:"localidade"`
}

type ViaCepService struct {
	httpClient http_client.HttpClientInterface
}

func NewViaCepService(httpClient http_client.HttpClientInterface) *ViaCepService {
	return &ViaCepService{
		httpClient: httpClient,
	}
}

func (vc *ViaCepService) CarregarDadosLocalidadePorCep(ctx context.Context, cep string) (string, error) {
	log.Println("Buscando localidade para o CEP:", cep)

	if cep == "" {
		return "", ErrCepObrigatorio
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
		return "", ErrCepNaoEncontrado
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrBuscaCep
	}

	var viaCep ViaCep

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &viaCep); err != nil {
		return "", err
	}

	if viaCep.Localidade == "" {
		return "", ErrCepNaoEncontrado
	}

	log.Println("Localidade encontrada:", viaCep.Localidade)
	return viaCep.Localidade, nil
}
