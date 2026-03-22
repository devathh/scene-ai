FROM golang:1.25.5-alpine AS builder

WORKDIR /app
RUN apk add --no-cache openssl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN openssl genrsa -out jwtRS256.key 2048 && \
    openssl rsa -in jwtRS256.key -pubout -out jwtRS256.key.pub

RUN go build -o app cmd/main.go

FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app ./app
COPY --from=builder /app/configs/ ./configs/
COPY --from=builder /app/prompts/ ./prompts/
COPY --from=builder /app/docs/ ./docs/
COPY --from=builder /app/jwtRS256.key ./jwtRS256.key
COPY --from=builder /app/jwtRS256.key.pub ./jwtRS256.key.pub

EXPOSE 8090

CMD ["./app"]