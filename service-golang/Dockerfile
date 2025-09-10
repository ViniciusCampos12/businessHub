FROM golang:1.25-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git bash curl

RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install gotest.tools/gotestsum@latest

RUN adduser -D -u 1000 appuser

ENV GOPATH=/app/.cache/go
ENV GOMODCACHE=$GOPATH/pkg/mod
ENV GOCACHE=$GOPATH/cache
RUN mkdir -p $GOPATH/pkg/mod $GOCACHE && chown -R appuser:appuser /app/.cache

USER appuser

ENV GOPROXY=https://goproxy.io,direct
COPY --chown=appuser:appuser go.mod go.sum ./
RUN go mod download

COPY --chown=appuser:appuser . .

FROM golang:1.25-alpine AS build

COPY --from=dev /app/. .

ENV GOPROXY=https://goproxy.io,direct
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/service-golang ./cmd/main.go

FROM alpine:latest AS prod

WORKDIR /app

COPY --from=build /app/service-golang .

RUN chmod +x /app/service-golang

EXPOSE 8080

CMD ["./service-golang"]
