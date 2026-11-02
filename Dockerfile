FROM golang:1.27 as build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-app cmd/weather/main.go

FROM scratch 

WORKDIR /app

COPY --from=build /app/weather-app .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["./weather-app"]