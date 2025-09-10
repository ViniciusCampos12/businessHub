FROM golang:1.25-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git bash curl

RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install gotest.tools/gotestsum@latest

ENV GOPROXY=https://goproxy.io,direct
COPY go.* ./
RUN go mod download

COPY . .

RUN adduser -D -u 1000 appuser

ENV GOPATH=/app/.cache/go
ENV GOMODCACHE=$GOPATH/pkg/mod
ENV GOCACHE=$GOPATH/cache
RUN mkdir -p /app/.cache/go && chown -R appuser:appuser /app/.cache

USER appuser

FROM golang:1.25-alpine AS build

COPY --from=dev /app/. .

ENV GOPROXY=https://goproxy.io,direct
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-app ./cmd/main.go

FROM alpine:latest AS prod

WORKDIR /app

COPY --from=build /app/go-app .

RUN chmod +x /app/go-app

EXPOSE 8080

CMD ["./go-app"]
