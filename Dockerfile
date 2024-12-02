FROM golang:1.23-alpine as builder

WORKDIR /app

RUN apk add --no-cache build-base git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .

COPY templates ./templates
COPY static ./static

EXPOSE 8080

CMD ["./main"]
