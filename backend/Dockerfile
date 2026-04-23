# Build API binary (static, no CGO).
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /mustore .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
ENV TZ=Europe/Moscow
ENV CONFIG_PATH=/app/config.yaml

COPY --from=builder /mustore ./mustore
COPY config/docker.yaml ./config.yaml
RUN chown nobody:nobody /app/mustore /app/config.yaml && chmod 500 /app/mustore && chmod 400 /app/config.yaml

EXPOSE 8000
USER nobody
CMD ["./mustore"]
