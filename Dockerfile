# BUILDER BINARKI
FROM golang:1.23-alpine as builder

WORKDIR /app

# build-base zestaw podswatowych narzędzi (trzeba żeby go skompilować)
RUN apk add --no-cache build-base git

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# KONTENER DO RUNOWANIA SERVERA
FROM alpine:latest

WORKDIR /app

# DO HTTPS
RUN apk add --no-cache ca-certificates

# KOPIOWANIE BINARKI
COPY --from=builder /app/main .

COPY static ./static
EXPOSE 8080

CMD ["./main"]
