
# entities-api

## Descrição

O projeto "entities-api" contém a API de entidades do domínio (customer, user, vehicle, parts supply e service), com autenticação JWT e rotas HTTP compatíveis com o padrão do projeto `mecanica_xpto` para uso local via Insomnia.

Event Storming: <https://miro.com/app/board/uXjVIgU2y2I=/>

## Pré-requisitos

- Golang 1.24 ou superior  
- Docker  
- Docker Compose  
- Make (para facilitar comandos via Makefile)

## Instalação (Tech Challange 1)

1. Clone o repositório:

   ```bash
   git clone <url-do-repositorio>
   ```

2. Navegue até o diretório do projeto:

   ```bash
  cd entities-api
   ```

3. Inicialize o ambiente, que vai:
   - copiar o arquivo `.env.example` para `.env` (sem sobrescrever se já existir)
   - instalar o Swag CLI (se necessário)
   - gerar a documentação Swagger
   - subir os containers Docker (app e banco)
   - aguardar banco e app ficarem prontos
   - executar migrations e seeds

   Para isso, execute:

   ```bash
   make init
   ```

4. A aplicação estará disponível em `http://localhost:8080`

## Comandos úteis

- Para subir os containers (build + background):

  ```bash
  make up
  ```

- Para parar e remover os containers:

  ```bash
  make down
  ```

- Para acompanhar os logs do container da aplicação:

  ```bash
  make logs
  ```

## Endpoints disponíveis

Base path: `http://localhost:8080/v1`

- `POST /login`
- `GET/POST/PATCH/DELETE /users`
- `GET/POST/PATCH/DELETE /customers`
- `GET/POST/PATCH/DELETE /vehicles`
- `GET/POST/PUT/DELETE /parts-supply`
- `GET/POST/PUT/DELETE /service`
- `GET /ping`

## Testes

- Para rodar os testes automatizados do escopo migrado dentro do docker:

```bash
make test
```

- Para rodar os testes de usecases com **100% de cobertura** (resumo no terminal):

```bash
make coverage

```

Para gerar e abrir o relatório de cobertura em HTML dos usecases:

```bash
make coverage-html
```

## Documentação da API

- Para gerar a documentação Swagger (a partir dos comentários no código):

  ```bash
  make swag-generate
  ```

A documentação Swagger estará disponível em:  
`http://localhost:8080/swagger/index.html` enquanto a aplicação estiver rodando.
