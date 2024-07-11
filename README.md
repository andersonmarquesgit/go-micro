# Arquitetura do Sistema

![Diagrama de Arquitetura](path/to/your/image.png)

## Descrição

Este diagrama ilustra a arquitetura do sistema, detalhando as interações entre os diferentes componentes. Abaixo estão descrições das principais partes do sistema:

- **User**: Interage com o sistema via browser.
- **Broker**: Central de comunicação, rodando em Docker, que coordena as interações entre os componentes.
- **Auth**: Serviço de autenticação, rodando em Docker, que gerencia a autenticação dos usuários e se comunica com o PostgreSQL.
- **Postgres**: Banco de dados relacional, rodando em Docker, usado para armazenar informações críticas.
- **Mail**: Serviço responsável por enviar e-mails.
- **Listener**: Componente que escuta mensagens da fila RabbitMQ e realiza ações baseadas nas mensagens recebidas.
- **RabbitMQ**: Sistema de fila de mensagens usado para comunicação assíncrona entre os serviços.
- **Mongo**: Banco de dados NoSQL usado para armazenar logs e outras informações não relacionais.
- **Logger**: Serviço que registra logs de diferentes partes do sistema, comunicando-se com o MongoDB.

### Fluxo de Dados

1. O usuário interage com o sistema via browser.
2. As solicitações são gerenciadas pelo Broker, que encaminha as solicitações para os serviços apropriados.
3. O serviço Auth valida a autenticação dos usuários e interage com o PostgreSQL para armazenamento e recuperação de dados de autenticação.
4. O serviço Mail envia e-mails necessários.
5. O Listener escuta mensagens do RabbitMQ e processa as mensagens conforme necessário.
6. O Logger registra logs das operações e interage com o MongoDB para armazenamento de logs.

### Comunicação entre Serviços

A comunicação entre os serviços é realizada de várias formas para testar diferentes métodos de integração:

- **RPC**: Comunicação entre o Broker e o Logger.
- **REST HTTP com JSON**: Outra forma de comunicação entre o Broker e o Logger.
- **gRPC**: Um método adicional de comunicação entre o Broker e o Logger.
- **RabbitMQ**: Utilizado para comunicação assíncrona entre os serviços.

### Componentes em Docker

Os serviços Broker, Auth e Postgres estão rodando em containers Docker para garantir isolamento, escalabilidade e facilidade de gerenciamento.

### Configuração do Docker Compose

O arquivo `docker-compose.yml` configura todos os serviços necessários para o sistema:

```yaml
version: '3.8'

services:
  broker-service:
    build:
      context: ./../broker-service
      dockerfile: ./../broker-service/broker-service.dockerfile
    restart: always
    ports:
      - "8080:80"
    deploy:
      mode: replicated
      replicas: 1

  logger-service:
    build:
      context: ./../logger-service
      dockerfile: ./../logger-service/logger-service.dockerfile
    restart: always
    deploy:
      mode: replicated
      replicas: 1
    environment:
      MONGO_URI: "mongodb://admin:password@mongo:27017"

  mail-service:
    build:
      context: ./../mail-service
      dockerfile: ./../mail-service/mail-service.dockerfile
    restart: always
    deploy:
      mode: replicated
      replicas: 1
    environment:
      MAIL_DOMAIN: localhost
      MAIL_HOST: mailhog
      MAIL_PORT: 1025
      MAIL_ENCRYPTION: none
      MAIL_USERNAME: ""
      MAIL_PASSWORD: ""
      MAIL_FROM: "Anderson Marques"
      MAIL_ADDRESS: "anderson@example.com"

  authentication-service:
    build:
      context: ./../authentication-service
      dockerfile: ./../authentication-service/authentication-service.dockerfile
    restart: always
    ports:
      - "8081:80"
    deploy:
      mode: replicated
      replicas: 1
    environment:
      DSN: "host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"
    depends_on:
      - postgres

  listener-service:
    build:
      context: ./../listener-service
      dockerfile: ./../listener-service/listener-service.dockerfile
    restart: always
    deploy:
      mode: replicated
      replicas: 1

  postgres:
    image: 'postgres:14.0'
    ports:
      - "5432:5432"
    restart: always
    deploy:
      mode: replicated
      replicas: 1
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: users
    volumes:
      - postgres-data:/var/lib/postgresql/data

  mongo:
    image: 'mongo:4.2.16-bionic'
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_DATABASE: logs
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password
    volumes:
      - mongo-data:/data/db

  mailhog:
    image: 'mailhog/mailhog:latest'
    ports:
      - "1025:1025"
      - "8025:8025"

  rabbitmq:
    image: 'rabbitmq:3.13.4-management-alpine'
    ports:
      - "5672:5672"
      - "15672:15672"
    deploy:
      mode: replicated
      replicas: 1
    volumes:
      - rabbitmq-data:/var/lib/rabbitmq

volumes:
  postgres-data:
  mongo-data:
  rabbitmq-data:
```

### Interface de Teste dos Microserviços
Para facilitar os testes das várias formas de integração entre os microserviços, foi desenvolvida uma interface frontend. Esta interface permite testar facilmente as diferentes formas de comunicação, como mostra a imagem abaixo:


### Funcionalidades da Interface de Teste
- **Test Broker:** Testa a comunicação com o serviço Broker.
- **Test Auth:** Testa a comunicação com o serviço de autenticação.
- **Test Log:** Testa a comunicação com o serviço de log.
- **Test Mail:** Testa a comunicação com o serviço de envio de e-mails.
- **Test Log via RabbitMQ:** Testa a comunicação com o serviço de log via RabbitMQ.
- **Test Log via RPC:** Testa a comunicação com o serviço de log via RPC.

### Exibição dos Resultados
- **Sent:** Exibe as mensagens enviadas durante os testes.
- **Received:** Exibe as respostas recebidas durante os testes.
- **Output:** Área onde são exibidos os resultados dos testes realizados.

### Makefile
Para facilitar a gestão e a automação dos comandos do Docker e do Go, um Makefile foi configurado com as seguintes regras:

```yml
FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
LOGGER_BINARY=loggerServiceApp
MAIL_BINARY=mailServiceApp
LISTENER_BINARY=listenerApp
AUTH_BINARY=authApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mail build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -o ${LISTENER_BINARY} .
	@echo "Done!"

## build_front: builds the front end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -buildvcs=false -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"
```

### Como Usar o Makefile
- Para iniciar os containers Docker:
```yml
make up
```

- Para parar os containers Docker:
```yml
make down
```

- Para construir e iniciar todos os serviços:
```yml
make up_build
```

- Para construir serviços específicos:
```yml
make build_broker
make build_logger
make build_auth
make build_mail
make build_listener
make build_front
```

- Para iniciar o frontend:
```yml
make start
```

- Para parar o frontend:
```yml
make stop
```

Para mais detalhes sobre a implementação e configuração, consulte a documentação específica de cada componente.
