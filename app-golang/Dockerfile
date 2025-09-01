FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-app ./cmd/main.go
RUN ls -lh /app

FROM alpine:latest AS prod

WORKDIR /app

COPY --from=builder /app/go-app .

RUN chmod +x /app/go-app

EXPOSE 8080

CMD ["./go-app"]
