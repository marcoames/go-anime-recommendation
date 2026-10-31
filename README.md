# Sistema de Recomendação de Anime

Um serviço web baseado em Go que fornece recomendações personalizadas de anime usando filtragem baseada em conteúdo e similaridade de cosseno.

## Visão Geral

Este projeto é um sistema de recomendação de anime que obtém dados da Jikan API (API não oficial do MyAnimeList) e fornece recomendações personalizadas com base em diversos recursos, incluindo gêneros, demografias, estúdios, popularidade e pontuações.

## Funcionalidades

- **Coleta de Dados de Anime**: Busca e armazena dados de anime da Jikan API através do modo `fetch`.
- **Recomendações Baseadas em Conteúdo**: Usa similaridade de cosseno para encontrar animes semelhantes com base em múltiplos recursos.
- **API RESTful**: Endpoint HTTP simples para obter recomendações de anime.
- **Integração com MongoDB:**: Armazenamento persistente dos dados de anime.
- **Suporte a Docker**:  Implantação facilitada usando contêineres Docker.
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

## Uso da API

Envie uma requisição GET para obter recomendações:

```http
GET /api/?anime=<anime_title>
````
Exemplo de resposta
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
    - Evita múltiplos títulos da mesma franquia quando existem alternativas
      igualmente relevantes.
    - Retorna até 6 animes mais similares como recomendações.

### Diversidade e franquias

A diversificação acontece depois da busca por similaridade. O sistema usa, nesta
ordem, `FranchiseID` quando disponível, títulos em inglês/japonês, título
principal e sinônimos. Também reconhece sufixos comuns como `Season`, `Movie`,
`Arc`, `OVA` e `Special`.

O coletor atual não consulta relações extras da Jikan para cada anime. Portanto,
`FranchiseID` é um campo opcional preparado para uma futura carga enriquecida;
sem ele, a identificação depende dos títulos e sinônimos retornados pela Jikan.
Quando a franquia não pode ser identificada, a recomendação ainda funciona, mas
não é possível garantir a remoção de todos os títulos relacionados.

As features ficam em cache na memória. Após alterar ou enriquecer os dados,
execute novamente o fetch e reinicie o backend:

```powershell
cd "back-end"
go run ./cmd/server -fetch
go run ./cmd/server
```


## Implantação

### Execução local com Docker

Suba o MongoDB e a API normalmente. Para carregar ou atualizar os dados antes de usar
as recomendações, execute o coletor uma vez:

```powershell
docker compose up -d mongodb
docker compose run --rm fetcher
docker compose up -d api frontend
```

O coletor grava a nova carga de forma atômica; uma falha durante a atualização
preserva a carga anterior.

### Execução local com MongoDB na nuvem

Para executar somente o backend e o frontend localmente, configure a URI do
MongoDB na sessão do PowerShell:

Copie `back-end/.env.example` para `back-end/.env` e preencha `MONGODB_URI`.
O arquivo `.env` real não deve ser versionado.

```powershell
cd "back-end"
$env:MONGODB_URI="mongodb+srv://usuario:senha@cluster.mongodb.net/?retryWrites=true&w=majority"
$env:LOG_COLOR="1"
go run ./cmd/server
```

O backend ficará disponível em `http://localhost:8080`. Em outro terminal:

```powershell
cd "front-end"
npm install
npm start
```

O frontend inicia usando o backend local (`http://localhost:8080`). A API pode
ser verificada em `http://localhost:8080/health`.


### Testes

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

Frontend:

```powershell
cd "front-end"
npm run build
```

O sistema está implantado e acessível via Render:

**Frontend (React)**: https://go-anime-recommendation-1.onrender.com/

**Backend (API)**: https://go-anime-recommendation.onrender.com/

Sinta-se à vontade para visitar o link do frontend e testar o sistema de recomendação de anime interativamente.
