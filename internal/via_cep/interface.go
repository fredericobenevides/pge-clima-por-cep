package via_cep

import (
	"context"
)

type ViaCepServiceInterface interface {
	CarregarDadosLocalidadePorCep(ctx context.Context, cep string) (string, error)
}
