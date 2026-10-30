package via_cep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/http_client"
	"github.com/fredericobenevides/pge-clima-por-cep/internal/internal_error"
)

var ErrCepNaoEncontrado = fmt.Errorf("Erro ao buscar localidade para o CEP informado")

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

func (vc *ViaCepService) CarregarDadosLocalidadePorCep(ctx context.Context, cep string) (string, *internal_error.InternalError) {
	log.Println("Buscando localidade para o CEP:", cep)

	if cep == "" {
		return "", internal_error.NewBadRequestError("CEP deve ser fornecido")
	}

	var url = fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", internal_error.NewInternalError(err.Error())
	}

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return "", internal_error.NewInternalError(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", internal_error.NewNotFoundError(ErrCepNaoEncontrado.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return "", internal_error.NewInternalError(fmt.Sprintf("Erro ao buscar localidade para o CEP: %s", cep))
	}

	var viaCep ViaCep

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", internal_error.NewInternalError(err.Error())
	}

	if err := json.Unmarshal(body, &viaCep); err != nil {
		return "", internal_error.NewInternalError(err.Error())
	}

	if viaCep.Localidade == "" {
		return "", internal_error.NewNotFoundError(ErrCepNaoEncontrado.Error())
	}

	log.Println("Localidade encontrada:", viaCep.Localidade)
	return viaCep.Localidade, nil
}
