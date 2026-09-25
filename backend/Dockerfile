FROM golang:1.27-alpine AS builder

WORKDIR /src

RUN apk add --no-cache \
    ca-certificates \
    git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/gowatch-api \
    ./cmd/api


FROM alpine:3.22 AS runner

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget \
    && addgroup -S gowatch \
    && adduser -S -G gowatch gowatch

WORKDIR /app

COPY --from=builder \
    /out/gowatch-api \
    ./gowatch-api

USER gowatch

EXPOSE 8081

CMD ["./gowatch-api"]