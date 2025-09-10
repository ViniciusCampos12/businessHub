FROM golang:1.25-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git bash curl

RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install gotest.tools/gotestsum@latest

COPY go.* ./
RUN go mod download

COPY . .

RUN adduser -D -u 1000 appuser

ENV GOMODCACHE=/app/.cache/go-mod
ENV GOCACHE=/app/.cache/go-build

RUN mkdir -p /app/.cache && chown -R appuser:appuser /app/.cache

USER appuser

FROM golang:1.25-alpine AS build

COPY --from=dev /app/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-app ./cmd/main.go

FROM alpine:latest AS prod

WORKDIR /app

COPY --from=build /app/go-app .

RUN chmod +x /app/go-app

EXPOSE 8080

CMD ["./go-app"]
