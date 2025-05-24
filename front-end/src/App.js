import React, { useState } from 'react';
import './App.css';

function App() {
  const [animeTitle, setAnimeTitle] = useState('');
  const [animeDetails, setAnimeDetails] = useState(null);
  const [recommendations, setRecommendations] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const randomAnimeList = [
    'Naruto',
    'One Piece',
    'Attack on Titan',
    'Fullmetal Alchemist: Brotherhood',
    'Demon Slayer: Kimetsu no Yaiba',
    'One Punch Man',
    'Death Note',
    'My Hero Academia',
    'Dragon Ball Z',
    'Tokyo Ghoul',
    'Sword Art Online',
    'Hunter x Hunter',
    'Jujutsu Kaisen',
    'Spirited Away',
  ];

  const fetchAnimeData = async (url) => {
    setLoading(true);
    setError('');
    setAnimeDetails(null);
    setRecommendations([]);

    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error('Anime not found');
      }
      const data = await response.json();

      setAnimeDetails(data.anime);
      setRecommendations(data.recommendations || []);
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async () => {
    if (!animeTitle.trim()) {
      setError('Please enter an anime title');
      return;
    }
    await fetchAnimeData(
      `https://go-anime-recommendation.onrender.com/api/?anime=${encodeURIComponent(
        animeTitle
      )}`
    );
  };

  const handleRandomAnime = async () => {
    const randomIndex = Math.floor(Math.random() * randomAnimeList.length);
    const randomTitle = randomAnimeList[randomIndex];
    setAnimeTitle(randomTitle);
    await fetchAnimeData(
      `https://go-anime-recommendation.onrender.com/api/?anime=${encodeURIComponent(
        randomTitle
      )}`
    );
  };

  return (
    <div className="container">
      <h1>Anime Recommendations</h1>
      
      <div className="input-group">
      <input
        type="text"
        value={animeTitle}
        onChange={(e) => setAnimeTitle(e.target.value)}
        placeholder="Enter an anime title"
      />
      <button onClick={handleSearch} disabled={loading}>
        {loading ? 'Loading...' : 'Get Recommendations'}
      </button>
      <button
        onClick={handleRandomAnime}
        disabled={loading}
        className="random-button"
      >
        🎲 Random
      </button>
    </div>

      {loading && <div className="spinner"></div>}

      {error && <div className="error">{error}</div>}

      {animeDetails && !loading && (
        <div className="anime-details">
          <h2>{animeDetails.title_english || animeDetails.title}</h2>
          <img
            src={animeDetails.images.jpg.image_url}
            alt={animeDetails.title_english || animeDetails.title}
          />
          <p>
            <strong>Genres:</strong> {animeDetails.genres?.map((g) => g.name).join(', ')}
          </p>
          <p><strong>Rating:</strong> {animeDetails.rating}</p>
          <p><strong>Status:</strong> {animeDetails.status}</p>
          <p><strong>Score:</strong> {animeDetails.score}</p>
          <p>
            <strong>MAL Page:</strong>{' '}
            <a href={animeDetails.url} target="_blank" rel="noopener noreferrer">
              Link to MyAnimeList
            </a>
          </p>
          <p>
            <strong>Trailer:</strong>{' '}
            <a href={animeDetails.trailer.embed_url} target="_blank" rel="noopener noreferrer">
              Watch Trailer
            </a>
          </p>
          <p><strong>Synopsis:</strong> {animeDetails.synopsis}</p>
        </div>
      )}

      {recommendations.length > 0 && !loading && (
        <div className="recommendations">
          <h3>Recommendations:</h3>
          <ul>
            {recommendations.map((anime, index) => (
              <li key={index} className="anime-item">
                <div className="anime-info">
                  <div className="image-container">
                    <img src={anime.images.jpg.image_url} alt={anime.title} />
                  </div>
                  <div className="details">
                    <h2>{anime.title}</h2>
                    <p>
                      <strong>Genres:</strong> {anime.genres?.map((g) => g.name).join(', ')}
                    </p>
                    <p><strong>Rating:</strong> {anime.rating}</p>
                    <p><strong>Status:</strong> {anime.status}</p>
                    <p><strong>Score:</strong> {anime.score}</p>
                    <p>
                      <strong>MAL Page: </strong>{' '}
                      <a href={anime.url} target="_blank" rel="noopener noreferrer">
                        Link to MyAnimeList
                      </a>
                    </p>
                    <p>
                      <strong>Trailer: </strong>{' '}
                      {anime.trailer?.embed_url ? (
                        <a href={anime.trailer.embed_url} target="_blank" rel="noopener noreferrer">
                          Watch Trailer
                        </a>
                      ) : (
                        <span>No trailer available</span>
                      )}
                    </p>
                    <p><strong>Synopsis:</strong> {anime.synopsis}</p>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

export default App;
