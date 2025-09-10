# BusinessHub

Este projeto é uma aplicação em Go que utiliza Docker, MongoDB e RabbitMQ.

## Pré-requisitos

- Docker e Docker Compose instalados
- Git

## Instalação e Setup

1. Clone o repositório:  
```bash
git clone https://github.com/ViniciusCampos12/businessHub/
cd businessHub
```

2. Crie o arquivo de variáveis de ambiente:  
```bash
cp .env.example .env
```
Ajuste as variáveis conforme necessário.

3. Construa e inicie os containers:  
```bash
docker compose up -d --build
```

## Acessando a aplicação

- Aplicação disponível em: `http://localhost:<GO_PORT>` (porta definida no `.env`)
- Swagger: `http://localhost/api/swagger/index.html`
- RabbitMQ Management: `http://localhost:15672`  
Usuário e senha do RabbitMQ configurados no `.env`

## Rodando testes

1. Acesse o container de testes:  
```bash
docker exec -it go-app-test sh
```

2. Rode os testes com gotestsum para uma visualização mais amigável dos resultados.:  
```bash
gotestsum --format testname
```

3. Para gerar cobertura dos testes e relatório HTML:  
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

- O arquivo `coverage.html` pode ser aberto no navegador para visualizar a cobertura detalhada.

## Observações

- O serviço `go-app-test` é separado do container de produção (`go-app`) para permitir rodar testes e gerar coverage sem interferir na aplicação principal.  
- Certifique-se de que MongoDB e RabbitMQ estão acessíveis nas URLs definidas no `.env`.  
- Para reiniciar os containers:  
```bash
docker compose down
docker compose up -d --build
```

- O arquivo `.env` deve conter todas as variáveis de configuração necessárias, incluindo:  
  - `RABBITMQ_URI`  
  - `MONGO_URI`  
  - `RABBITMQ_USER`  
  - `RABBITMQ_PASW`  
  - `MONGO_DATABASE`  
  - `RABBITMQ_COMPANY_QUEUE`  
  - `WEBSOCKET_URL`
  - `TARGET`
  - `APP_GOLANG_PORT`
  - `SERVICE_GOLANG_PORT`
  - `MONGO_EXTERNAL_PORT`

## Integração via WebSockets com outro serviço

O serviço `service-golang` é um microsserviço que consome mensagens da fila do RabbitMQ, publicadas pelo `go-app`.

Para testar a integração:

1. Abra o arquivo `index.html` na raiz do projeto no navegador e conecte-se ao WebSocket.
2. Realize qualquer operação no `go-app`.
3. O evento será enviado para o `service-golang` via RabbitMQ, que:
   - Armazena o evento na tabela `outbox_events` do banco de dados.
   - Publica o evento para todos os clientes conectados ao WebSocket.
