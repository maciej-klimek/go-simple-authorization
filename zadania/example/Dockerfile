# Wykorzystanie oficjalnego obrazu Go
FROM golang:1.20

# Ustawienie katalogu roboczego
WORKDIR /app

# Kopiowanie plików go.mod i go.sum do obrazu
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Kopiowanie reszty plików aplikacji
COPY backend/ .

# Instalacja zależności
RUN go mod tidy

# Budowanie aplikacji
RUN go build -o main .

# Uruchamianie aplikacji
CMD ["./main"]