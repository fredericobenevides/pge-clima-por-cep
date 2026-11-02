package via_cep

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/fredericobenevides/pge-clima-por-cep/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindCityByCepWhenValidCEP(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"localidade":"São Paulo"}`)),
	}, nil)

	viaCepUseCase := NewViaCepUseCase(httpClientMock)
	localidade, err := viaCepUseCase.FindCityByCep(context.Background(), "01001000")

	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", localidade)
}

func TestFindCityByCepWhenCEPIsEmpty(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}

	viaCepUseCase := NewViaCepUseCase(httpClientMock)
	_, err := viaCepUseCase.FindCityByCep(context.Background(), "")

	assert.NotNil(t, err)
}

func TestFindCityByCepWhenInvalidCEPForBadRequest(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader(`{"erro":"true"}`)),
	}, nil)

	viaCepUseCase := NewViaCepUseCase(httpClientMock)
	_, err := viaCepUseCase.FindCityByCep(context.Background(), "invalid-cep")

	assert.NotNil(t, err)
	assert.Equal(t, "can not find zipcode", err.Error())
}

func TestFindCityByCepWhenInvalidCEPForSuccessResponse(t *testing.T) {
	httpClientMock := &mocks.MockHttpClient{}
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"erro":"true"}`)),
	}, nil)

	viaCepUseCase := NewViaCepUseCase(httpClientMock)
	_, err := viaCepUseCase.FindCityByCep(context.Background(), "invalid-cep")

	assert.NotNil(t, err)
	assert.Equal(t, "can not find zipcode", err.Error())
}
