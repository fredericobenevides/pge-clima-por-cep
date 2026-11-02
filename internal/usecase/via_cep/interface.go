package via_cep

import (
	"context"
)

type ViaCepUseCaseInterface interface {
	FindCityByCep(ctx context.Context, cep string) (string, error)
}
