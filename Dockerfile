# Etapa 1: Construcción del binario
FROM golang:1.23-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o go-app .

# Etapa 2: Imagen distroless para la ejecución
FROM gcr.io/distroless/base-debian10

WORKDIR /app
COPY --from=build /app/go-app .

EXPOSE 8080
USER nonroot:nonroot
CMD ["./go-app"]