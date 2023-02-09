# Star Wars API

## Requisitos

- [go](https://tip.golang.org/doc/go1.19)
- [mockgen](https://github.com/golang/mock)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [swaggo](https://github.com/swaggo/swag)

## Configuração

### Criando o arquivo .env

Para o ambiente local criar o arquivo `.env` como no exemplo abaixo:

```txt
PORT=3000
DB_USER=luke
DB_PASSWORD=xQlpKD95kp20Wa1JAX6O
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=star_wars
MIGRATE_URL=mysql://luke:xQlpKD95kp20Wa1JAX6O@tcp\(127.0.0.1:3306\)/star_wars?multiStatements=true
URL_STARWARS_API=https://swapi.dev/api
```

## Instalação

```bash
$ make install
```


### Migrate

Após executar o comando `docker-compose`, pode ser necessário esperar alguns segundos até o banco estar apto a receber comandos:

```bash
$ docker-compose up -d
$ make migrateup
```

## Executando

Antes de executar a API Rest, é recomendado executar o `make import` para buscar os dados da API [SWAPI](https://swapi.dev/):

```bash
$ make import

# API Rest
$ make run
```

Para visualizar a documentação das rotas localmente, após a API estiver em execução, basta acessar o [swagger](http://localhost:3000/swagger/index.html)

## Testes

```bash
# Testes unitários
$ make test

# Cobertura dos testes unitários
$ make test/cov
```

## Arquitetura do projeto

O código está organizado da seguinte forma:

- **db**: códigos referentes a banco de dados
    - **migrations**: SQLs para as `migrations`
- **docs**: arquivos swagger
- **log**: arquivos de logs
- **api**: path com as configuraçoes das rotas e handlers da api rest
- **app**: path com as regras de negócio
- **imports**: path com o comando de importação
- **model**: representações dos modelos
- **server**: path com os registers das camadas
- **store**: comunicação com o banco de dados e integrações com api de terceiros
- **mock**: arquivos `mock` para dar suporte aos testes unitários