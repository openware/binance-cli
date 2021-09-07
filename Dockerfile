FROM golang:1.14-alpine AS builder

RUN apk add --no-cache curl

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bin/binance-cli .

FROM alpine:3.12

WORKDIR app

COPY --from=builder /build/bin/binance-cli /usr/local/bin

CMD ["binance"]
