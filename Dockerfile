FROM golang:1.27 as build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o clima

FROM scratch 

WORKDIR /app

COPY --from=build /app/clima /app/clima

ENTRYPOINT ["./clima"]