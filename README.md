# PGE Clima por CEP

Projeto implementado com base nos requisitos da atividade **"Desafio Clima por CEP"** do curso Pós Go Expert da FullCycle.

A aplicação expõe uma API HTTP que recebe um CEP brasileiro, consulta a localidade na [ViaCEP](https://viacep.com.br/) e busca a temperatura atual dessa localidade na [WeatherAPI](https://www.weatherapi.com/). A resposta informa as temperaturas em Celsius, Fahrenheit e Kelvin.

## Arquitetura

O fluxo da requisição é:

1. A API recebe um CEP pela rota `GET /clima/{cep}`.
2. O CEP é validado: deve ter exatamente oito dígitos numéricos.
3. A ViaCEP informa a localidade correspondente ao CEP.
4. A WeatherAPI fornece a temperatura atual da localidade.
5. A aplicação converte a temperatura para Fahrenheit e Kelvin e devolve a resposta em JSON.

As chamadas a serviços externos têm timeout de 3 segundos.

## Estrutura do projeto

```text
.
├── cmd/
│   └── weather/
│       └── main.go                 # Inicialização do servidor HTTP
├── internal/
│   ├── handler/                    # Handler da rota de clima
│   ├── http_client/                # Cliente HTTP reutilizável
│   ├── internal_error/             # Erros HTTP da aplicação
│   ├── via_cep/                    # Integração com a ViaCEP
│   └── weather_api/                # Integração com a WeatherAPI
├── pkg/
│   └── converter/                  # Conversões de temperatura
├── test/
│   ├── mocks/                      # Mocks utilizados nos testes
│   └── temperatura.http            # Requisições HTTP de exemplo
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Pré-requisitos

- Go 1.27 ou superior;
- uma chave de API da [WeatherAPI](https://www.weatherapi.com/);
- Docker (opcional, para execução em container).

## Configuração

Defina a variável de ambiente `WEATHER_API_KEY` com sua chave da WeatherAPI.

Para desenvolvimento local, crie um arquivo `.env` na raiz do projeto:

```env
WEATHER_API_KEY=sua_chave_da_weatherapi
```

## Como executar localmente

Na raiz do projeto, execute:

```bash
go run ./cmd/weather
```

O servidor ficará disponível em `http://localhost:8080`.

Em outro terminal, consulte um CEP:

```bash
curl http://localhost:8080/clima/01001000
```

## Como executar com Docker Compose

Com o arquivo `.env` configurado, inicie o serviço definido no `docker-compose.yml`:

```bash
docker compose up -d
```

O serviço ficará disponível em `http://localhost:8080`.

## Endpoint disponível

| Método | Rota           | Descrição                                                        |
| ------ | -------------- | ---------------------------------------------------------------- |
| `GET`  | `/clima/{cep}` | Consulta a temperatura atual para a localidade do CEP informado. |

### Exemplo de requisição

```bash
curl http://localhost:8080/clima/01001000
```

## Deploy no Google Cloud Run

A aplicação está publicada no Google Cloud Run. A URL do serviço é:

```bash
curl https://pge-clima-por-cep-841576980828.us-central1.run.app/clima/01001000
```
