# Anime Recommendation System

A Go-based web service that provides personalized anime recommendations using content-based filtering and cosine similarity.

## Overview

This project is an anime recommendation system that fetches data from the Jikan API (MyAnimeList unofficial API) and provides personalized anime recommendations based on various features including genres, demographics, studios, popularity, and scores.

## Features

- **Anime Data Collection**: Automatically fetches and stores anime data from the Jikan API.
- **Content-Based Recommendations**: Uses cosine similarity to find similar anime based on multiple features.
- **RESTful API**: Simple HTTP endpoint for getting anime recommendations.
- **MongoDB Integration**: Persistent storage of anime data.
- **Docker Support**: Easy deployment using Docker containers.

## Architecture

The project consists of two main services:

1. **API Service**: Handles incoming HTTP requests and serves recommendations.
2. **Fetcher Service**: Periodically updates the anime database with fresh data.

## Technical Stack

- **Backend**: Go (Golang)
- **Database**: MongoDB
- **Container**: Docker
- **API**: RESTful HTTP endpoints
- **External API**: Jikan API (MyAnimeList)

## API Usage

Send a GET request to get recommendations:

```http
GET /api/?anime=<anime_title>
````
Example response
```json
{
  "anime": {
    "title": "Your searched anime",
    "...": "..."
  },
  "recommendations": [
    {
      "title": "Recommended anime 1",
      "...": "..."
    },
    {
      "title": "Recommended anime 2",
      "...": "..."
    }
  ]
}
```

## How to run 
```bash
docker-compose run fetcher
docker-compose up api -d
```

## How it works

1. The fetcher service collects anime data from the Jikan API.

2. Data is stored in MongoDB for quick access.

3. When a request comes in, the system:

    - Extracts features (genres, studios, etc.).

    - Applies weights to different features.

    - Uses cosine similarity to find similar anime.

    - Returns the top 4 most similar recommendations.
