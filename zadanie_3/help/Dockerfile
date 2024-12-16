# ==================================== BUILDER ===============================================
# Użyj minimalnego obrazu bazowego dla Go
FROM golang:1.23.2-alpine AS builder

# Dodanie użytkownika nieuprzywilejowanego dla większego bezpieczeństwa
RUN adduser -D -g '' appuser

# Ustawienie katalogu roboczego w kontenerze
WORKDIR /app

# Kopiowanie plików go.mod i go.sum w celu wykorzystania cache Dockera dla zależności
COPY go.mod go.sum ./

# Pobranie zależności (warstwa cache)
RUN go mod download

# Skopiowanie kodu źródłowego aplikacji
COPY . .

# Kompilacja binarki Go z flagami optymalizacyjnymi w celu zmniejszenia rozmiaru pliku wynikowego
RUN go build -ldflags="-s -w" -o main .


# ==================================== RUNTIME ===============================================
# Użycie minimalnego obrazu runtime dla mniejszego rozmiaru końcowego obrazu
FROM alpine:latest

# Instalacja niezbędnych zależności runtime (jeśli są potrzebne)
RUN apk --no-cache add ca-certificates

# Dodanie użytkownika w runtime, aby chown działało
RUN adduser -D -g '' appuser

# Ustawienie katalogu roboczego w kontenerze
WORKDIR /app

# Kopiowanie skompilowanej aplikacji z etapu budowy
COPY --from=builder /app/main ./

# Kopiowanie plików statycznych, uploadów i innych zasobów
COPY static ./static


# Kopiowanie plików użytkownika z etapu budowy
COPY --from=builder /etc/passwd /etc/passwd

# Stworzenie katalogu /shared-data i nadanie uprawnień dla użytkownika
RUN mkdir -p /shared-data && chown -R appuser:appuser /shared-data


# Przełącz na użytkownika nieuprzywilejowanego
USER appuser

# Udostępnienie portu aplikacji
EXPOSE 8080

# Ustawienie bezpiecznego entrypointa dla aplikacji
CMD ["./main"]
