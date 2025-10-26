# Desafio Clean Architecture - Listagem de Orders

Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
  Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

## Como Executar a Aplicação

### Pré-requisitos
- Go 1.24 ou superior
- Docker e Docker Compose
- MySQL (ou usar via Docker)

### Passos para executar

**Configurar as variáveis de ambiente:**

```bash
cp .env.example .env
```

#### Opção 1: Executar com Docker Compose (Recomendado)

1. **Subir todos os serviços com Docker Compose**
```bash
docker compose up -d
```

Este comando irá subir automaticamente:
- MySQL na porta 3307
- RabbitMQ na porta 5672 (5672 para conexão e 15672 para management)
- Aplicação na porta 8000 (REST), 8080 (GraphQL) e 50051 (gRPC)

As migrations serão executadas automaticamente pela aplicação usando `golang-migrate` quando ela iniciar.

2. **Verificar os logs**
```bash
docker compose logs -f app
```

3. **Parar os serviços**
```bash
docker compose down
```

4. **Parar e remover volumes (limpar dados)**
```bash
docker compose down -v
```

#### Opção 2: Executar localmente

1. **Subir apenas MySQL e RabbitMQ**
```bash
docker compose up -d mysql rabbitmq
```

4. **Executar a aplicação**
```bash
cd cmd/ordersystem && go run main.go wire_gen.go
```

### Portas dos Serviços

- **REST API**: `http://localhost:8000`
  - POST `/order` - Criar uma order
  - GET `/orders` - Listar todas as orders

- **gRPC Server**: `localhost:50051`

- **GraphQL Playground**: `http://localhost:8080`
  - Query: `listOrders` - Listar todas as orders
  - Mutation: `createOrder` - Criar uma order

### Exemplos de Uso

#### REST API

**Criar uma order:**
```bash
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{
    "id": "order-10",
    "price": 100.5,
    "tax": 10.5
  }'
```

**Listar todas as orders:**
```bash
curl http://localhost:8000/orders
```

#### GraphQL

**Criar uma order:**
```graphql
mutation createOrder {
    createOrder(input: { id: "ccc", Price: 12.2, Tax: 2 }) {
        id
        Price
        Tax
        FinalPrice
    }
}
```

**Listar todas as orders:**
```graphql
query listOrders {
    listOrders {
        id
        Price
        Tax
        FinalPrice
    }
}
```

#### gRPC

Use as ferramentas como `grpcurl` ou crie um cliente para testar os métodos:
- `CreateOrder` - Criar uma order
- `ListOrders` - Listar todas as orders
