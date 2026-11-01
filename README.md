# Sistema de Recomendação de Anime

Um serviço web baseado em Go que fornece recomendações personalizadas de anime usando filtragem baseada em conteúdo e similaridade de cosseno.

## Visão Geral

Este projeto é um sistema de recomendação de anime que obtém dados da Jikan API (API não oficial do MyAnimeList) e fornece recomendações personalizadas com base em diversos recursos, incluindo gêneros, demografias, estúdios, popularidade e pontuações.

## Funcionalidades

- **Coleta de Dados de Anime**: Busca e armazena dados de anime da Jikan API através do modo `fetch`.
- **Recomendações Baseadas em Conteúdo**: Usa similaridade de cosseno para encontrar animes semelhantes com base em múltiplos recursos.
- **API RESTful**: Endpoint HTTP simples para obter recomendações de anime.
- **Integração com MongoDB**: Armazenamento persistente dos dados de anime.
- **Suporte a Docker**: Implantação facilitada usando contêineres Docker.
- **Frontend em React**: Interface amigável para buscar e exibir recomendações de anime.

## Arquitetura

O projeto consiste nos seguintes serviços:

1. **Serviço de API**: Lida com requisições HTTP e fornece recomendações.
2. **Serviço de Coleta (Fetcher)**: Atualiza o banco de dados de anime com informações recentes.
3. **Frontend React**: Fornece uma UI interativa para os usuários buscarem animes e visualizarem recomendações.

## Estrutura do projeto

```text
back-end/
├── cmd/server/                 # Entrada do servidor HTTP e configuração da aplicação
├── internal/
│   ├── anime/                  # Modelo, coleta da Jikan e persistência no MongoDB
│   ├── api/                    # Handlers HTTP, rotas de negócio e testes da API
│   └── recommendation/         # Features, similaridade KNN, diversidade e testes
├── package/jikan/              # Cliente auxiliar para chamadas à Jikan API
├── go.mod                      # Dependências e configuração do módulo Go
└── Dockerfile                  # Imagem do backend

front-end/
├── src/                        # Componentes React, estilos e entrada da aplicação
├── public/                     # Arquivos públicos e metadados da aplicação
├── package.json                # Scripts e dependências do frontend
└── Dockerfile                  # Imagem do frontend
```

O backend segue a convenção comum de Go: `cmd` contém os executáveis e
`internal` contém pacotes privados da aplicação. Os testes ficam próximos do
código que validam. O frontend é uma aplicação React independente que consome a
API HTTP do backend.

## Stack Técnica

- **Backend**: Go (Golang)
- **Frontend**: React.js
- **Banco de Dados**: MongoDB
- **Contêiner**: Docker
- **API**: RESTful HTTP endpoints
- **API Externa**: Jikan API (MyAnimeList)

## Requisitos

- Go 1.21 ou superior
- Node.js e npm
- MongoDB local ou MongoDB Atlas
- Docker e Docker Compose, apenas se quiser executar com contêineres

## Configuração

O backend lê as variáveis de ambiente de `back-end/.env`. Para começar,
copie o arquivo de exemplo e preencha a URI do MongoDB:

```powershell
Copy-Item back-end\.env.example back-end\.env
```

Edite `back-end/.env` e informe `MONGODB_URI`. O arquivo `.env.example` deve ser usado
como referência.

Variáveis disponíveis:

```env
MONGODB_URI="mongodb+srv://usuario:senha@cluster.mongodb.net/anime?retryWrites=true&w=majority"
PORT="8080"
LOG_COLOR="1"
```

`PORT` é opcional e assume `8080`. `LOG_COLOR=1` habilita cores ANSI nos logs
do terminal. Para desabilitá-las, use `NO_COLOR=1`.

Se estiver usando MongoDB localmente, inicie o serviço antes de executar o
backend. Com MongoDB Atlas, use a string de conexão do cluster em
`MONGODB_URI`.

## Uso da API

Envie uma requisição GET para obter recomendações:

```http
GET /api/?anime=<anime_title>
```

Exemplo:

```http
GET /api/?anime=Demon%20Slayer
```

Para consultar sugestões enquanto o usuário digita:

```http
GET /api/suggestions?q=demon
```

Para verificar se o backend está disponível:

```http
GET /health
```

Exemplo de resposta:
```json
{
  "anime": {
    "title": "Anime buscado",
    "...": "..."
  },
  "recommendations": [
    {
      "title": "Anime recomendado 1",
      "...": "..."
    },
    {
      "title": "Anime recomendado 2",
      "...": "..."
    }
  ]
}
```

## Como Funciona

1. O serviço de coleta obtém dados de anime da Jikan API.

2. Os dados são armazenados no MongoDB para acesso rápido.

3. Quando uma requisição é recebida, o sistema:

    - Extrai recursos (gêneros, estúdios, etc.).

    - Normaliza score e popularidade e codifica categorias por identidade.

    - Usa similaridade de cosseno para gerar um pool dos 50 candidatos mais relevantes.

    - Reordena as recomendações com MMR para equilibrar relevância e variedade,
      sem aceitar candidatos muito menos similares.
    - Retorna até 6 animes mais similares como recomendações.

As features ficam em cache na memória. Após alterar ou enriquecer os dados,
execute novamente o fetch e reinicie o backend:

```powershell
cd "back-end"
go run ./cmd/server -fetch
go run ./cmd/server
```

O primeiro comando atualiza a base e encerra o processo; o segundo inicia a
API.

## Execução local

Após iniciar o backend com os comandos acima, ele ficará disponível em
`http://localhost:8080`.

O frontend deve ser iniciado em outro terminal:

```powershell
cd front-end
npm install
npm start
```

Abra `http://localhost:3000`. Por padrão, o frontend usa o backend local
(`http://localhost:8080`). Em modo de desenvolvimento, o botão no rodapé
permite alternar entre o backend local e o ambiente de produção.

## Execução com Docker

```powershell
docker compose up -d mongodb
docker compose run --rm fetcher
docker compose up -d api frontend
```

Nesse modo, o MongoDB é executado pelo próprio Docker e o backend usa
`mongodb://mongodb:27017`. O arquivo `.env` local não é necessário para esses
serviços.

Endereços:

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- Health check: `http://localhost:8080/health`

Para parar os serviços:

```powershell
docker compose down
```

## Testes

Backend:

```powershell
cd "back-end"
go test ./...
go vet ./...
```

Os testes da API usam `httptest` e um repositório falso, cobrindo respostas de
sucesso, validação `400`, anime inexistente `404`, falhas internas `500`,
sugestões e métodos inválidos `405`. Eles não exigem MongoDB ou conexão com a
internet.

Para executar os testes do frontend:

```powershell
cd "front-end"
npm test
npm run build
```

## Deploy atual

- Frontend: https://go-anime-recommendation-1.onrender.com/
- Backend: https://go-anime-recommendation.onrender.com/
