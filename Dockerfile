# --- builder stage -----------------------------------------------------------
FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY app/ ./app/
RUN go build -o /gomw-gw ./app/cmd/server

# --- runtime stage -----------------------------------------------------------
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /gomw-gw /usr/local/bin/gomw-gw

ENV LISTEN_ADDR=:8080
EXPOSE 8080
ENTRYPOINT ["gomw-gw"]