# Sistema de Recomendação de Anime

Um serviço web baseado em Go que fornece recomendações personalizadas de anime usando filtragem baseada em conteúdo e similaridade de cosseno.

## Visão Geral

Este projeto é um sistema de recomendação de anime que obtém dados da Jikan API (API não oficial do MyAnimeList) e fornece recomendações personalizadas com base em diversos recursos, incluindo gêneros, demografias, estúdios, popularidade e pontuações.

## Funcionalidades

- **Coleta de Dados de Anime**: Busca e armazena automaticamente dados de anime da Jikan API.
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

    - Aplica pesos a diferentes recursos.

    - Usa similaridade de cosseno para encontrar animes semelhantes.

    - Retorna os 4 animes mais similares como recomendações.


## Implantação

O sistema está implantado e acessível via Render:

**Frontend (React)**: https://go-anime-recommendation-1.onrender.com/

**Backend (API)**: https://go-anime-recommendation.onrender.com/

Sinta-se à vontade para visitar o link do frontend e testar o sistema de recomendação de anime interativamente.

