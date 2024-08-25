# Etapa de construção
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY cmd/stresstest/ ./cmd/stresstest/

RUN CGO_ENABLED=0 GOOS=linux go build -o loadtester ./cmd/stresstest/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/loadtester .

ENTRYPOINT ["./loadtester"]
