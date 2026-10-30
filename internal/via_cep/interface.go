package via_cep

import (
	"context"

	"github.com/fredericobenevides/pge-clima-por-cep/internal/internal_error"
)

type ViaCepServiceInterface interface {
	CarregarDadosLocalidadePorCep(ctx context.Context, cep string) (string, *internal_error.InternalError)
}
