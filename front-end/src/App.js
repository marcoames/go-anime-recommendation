import React, { useMemo, useState } from 'react';
import './App.css';

function App() {
  const [animeTitle, setAnimeTitle] = useState('');
  const [animeDetails, setAnimeDetails] = useState(null);
  const [recommendations, setRecommendations] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [backend, setBackend] = useState('production');

  const BACKENDS = {
    local: 'http://localhost:8080',
    production: process.env.REACT_APP_BACKEND_URL,
  };

  const backendUrl = BACKENDS[backend];

  const randomAnimeList = useMemo(
    () => [
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
    ],
    []
  );

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
    } catch (err) {
      setError(err.message);
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
      `${backendUrl}/?anime=${encodeURIComponent(animeTitle)}`
    );
  };

  const handleRandomAnime = async () => {
    const randomIndex = Math.floor(Math.random() * randomAnimeList.length);
    const randomTitle = randomAnimeList[randomIndex];
    setAnimeTitle(randomTitle);

    await fetchAnimeData(
      `${backendUrl}/api/?anime=${encodeURIComponent(randomTitle)}`
    );
  };

  const getTrailerLink = (anime) => {
    const ytId = anime?.trailer?.youtube_id;
    if (ytId) return `https://www.youtube.com/watch?v=${ytId}`;

    const url = anime?.trailer?.url;
    if (url) return url;

    const embed = anime?.trailer?.embed_url;
    const match = embed?.match(/\/embed\/([^?]+)/);
    if (match?.[1]) return `https://www.youtube.com/watch?v=${match[1]}`;

    return null;
  };

  const getYear = (anime) => {
    const from = anime?.aired?.from;
    if (!from) return null;
    const d = new Date(from);
    return Number.isNaN(d.getTime()) ? null : d.getFullYear();
  };

  const displayTitle = (anime) => anime?.title_english || anime?.title || 'Untitled';

  const formatStudios = (anime) => {
    const studios = anime?.studios?.map((s) => s.name).filter(Boolean) || [];
    return studios.length ? studios.join(', ') : '—';
  };

  const formatGenres = (anime) => {
    const genres = anime?.genres?.map((g) => g.name).filter(Boolean) || [];
    return genres;
  };

  const safeScore = (score) => (typeof score === 'number' ? score.toFixed(2) : score ?? '—');

  const trailerLink = getTrailerLink(animeDetails);

  return (
    <div className="container">
      <h1>go-anime-recommendation</h1>

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

      {loading && <div className="spinner" />}
      {error && <div className="error">{error}</div>}

      {animeDetails && !loading && (
        <section className="anime-details" aria-label="Selected anime">
          <div className="poster-col">
            <img
              src={animeDetails?.images?.jpg?.image_url}
              alt={displayTitle(animeDetails)}
              loading="lazy"
            />
          </div>

          <div className="content-col">
            <header className="title-block">
              <h2>{displayTitle(animeDetails)}</h2>
              {animeDetails?.title !== animeDetails?.title_english && animeDetails?.title ? (
                <p className="subtitle">{animeDetails.title}</p>
              ) : null}

              <div className="quick-stats" role="list">
                <span className="badge" role="listitem">
                  ⭐ {safeScore(animeDetails.score)}
                </span>
                {getYear(animeDetails) ? (
                  <span className="badge" role="listitem">
                    {getYear(animeDetails)}
                  </span>
                ) : null}
                {animeDetails?.episodes ? (
                  <span className="badge" role="listitem">
                    {animeDetails.episodes} eps
                  </span>
                ) : null}
                {animeDetails?.type ? (
                  <span className="badge" role="listitem">
                    {animeDetails.type}
                  </span>
                ) : null}
                {animeDetails?.status ? (
                  <span className="badge" role="listitem">
                    {animeDetails.status}
                  </span>
                ) : null}
              </div>
            </header>

            {animeDetails?.synopsis ? (
              <details className="synopsis" open>
                <summary>Synopsis</summary>
                <p>{animeDetails.synopsis}</p>
              </details>
            ) : null}

            <div className="info-grid" aria-label="Anime metadata">
              <div className="info-item">
                <span className="label">Rating</span>
                <span className="value">{animeDetails?.rating || '—'}</span>
              </div>
              <div className="info-item">
                <span className="label">Source</span>
                <span className="value">{animeDetails?.source || '—'}</span>
              </div>
              <div className="info-item">
                <span className="label">Studio</span>
                <span className="value">{formatStudios(animeDetails)}</span>
              </div>
              <div className="info-item">
                <span className="label">Aired</span>
                <span className="value">{animeDetails?.aired?.string || '—'}</span>
              </div>
            </div>

            {formatGenres(animeDetails).length > 0 ? (
              <div className="genre-list" aria-label="Genres">
                {formatGenres(animeDetails).map((g) => (
                  <span key={g} className="genre-chip">
                    {g}
                  </span>
                ))}
              </div>
            ) : null}

            <div className="actions" aria-label="External links">
              {animeDetails?.url ? (
                <a
                  className="action-link"
                  href={animeDetails.url}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  View on MyAnimeList
                </a>
              ) : null}

              {trailerLink ? (
                <a
                  className="action-link secondary"
                  href={trailerLink}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Watch trailer
                </a>
              ) : null}
            </div>
          </div>
        </section>
      )}

      {recommendations.length > 0 && !loading && (
        <section className="recommendations" aria-label="Recommendations">
          <h3>Recommendations</h3>

          <ul>
            {recommendations.map((anime, index) => {
              const recTrailerLink = getTrailerLink(anime);

              return (
                <li key={`${anime?.mal_id ?? 'rec'}-${index}`} className="anime-item">
                  <article className="anime-info">
                    <div className="image-container">
                      <img
                        src={anime?.images?.jpg?.image_url}
                        alt={displayTitle(anime)}
                        loading="lazy"
                      />
                    </div>

                    <div className="details">
                      <header className="rec-header">
                        <h2 title={displayTitle(anime)}>{displayTitle(anime)}</h2>

                        <div className="quick-stats" role="list">
                          <span className="badge" role="listitem">
                            ⭐ {safeScore(anime?.score)}
                          </span>
                          {getYear(anime) ? (
                            <span className="badge" role="listitem">
                              {getYear(anime)}
                            </span>
                          ) : null}
                          {anime?.episodes ? (
                            <span className="badge" role="listitem">
                              {anime.episodes} eps
                            </span>
                          ) : null}
                          {anime?.type ? (
                            <span className="badge" role="listitem">
                              {anime.type}
                            </span>
                          ) : null}
                        </div>
                      </header>

                      <div className="info-grid compact" aria-label="Recommendation metadata">
                        <div className="info-item">
                          <span className="label">Status</span>
                          <span className="value">{anime?.status || '—'}</span>
                        </div>
                        <div className="info-item">
                          <span className="label">Rating</span>
                          <span className="value">{anime?.rating || '—'}</span>
                        </div>
                        <div className="info-item">
                          <span className="label">Studio</span>
                          <span className="value">{formatStudios(anime)}</span>
                        </div>
                        <div className="info-item">
                          <span className="label">Aired</span>
                          <span className="value">{anime?.aired?.string || '—'}</span>
                        </div>
                      </div>

                      {anime?.synopsis ? (
                        <details className="synopsis compact">
                          <summary>Synopsis</summary>
                          <p>{anime.synopsis}</p>
                        </details>
                      ) : null}

                      {formatGenres(anime).length > 0 ? (
                        <div className="genre-list" aria-label="Genres">
                          {formatGenres(anime).map((g) => (
                            <span key={`${anime?.mal_id ?? index}-${g}`} className="genre-chip">
                              {g}
                            </span>
                          ))}
                        </div>
                      ) : null}

                      <div className="actions" aria-label="External links">
                        {anime?.url ? (
                          <a
                            className="action-link"
                            href={anime.url}
                            target="_blank"
                            rel="noopener noreferrer"
                          >
                            MyAnimeList
                          </a>
                        ) : null}

                        {recTrailerLink ? (
                          <a
                            className="action-link secondary"
                            href={recTrailerLink}
                            target="_blank"
                            rel="noopener noreferrer"
                          >
                            Trailer
                          </a>
                        ) : (
                          <span className="action-link disabled" aria-disabled="true">
                            No trailer
                          </span>
                        )}
                      </div>
                    </div>
                  </article>
                </li>
              );
            })}
          </ul>
        </section>
      )}

      <div className="container">
        {/* ... all existing JSX stays exactly the same ... */}

        {/* Dev toggle - only visible in development */}
        {process.env.NODE_ENV === 'development' && (
          <div style={{
            position: 'fixed',
            bottom: 8,
            right: 8,
            fontSize: '0.7rem',
            color: '#666',
            background: '#1a1a1a',
            padding: '4px 8px',
            borderRadius: 4,
            zIndex: 9999,
            cursor: 'pointer',
          }}
            onClick={() => setBackend(b => b === 'production' ? 'local' : 'production')}
          >
            {backend === 'production' ? '🌐 prod' : '💻 local'}
          </div>
        )}
      </div>

    </div>
  );
}



export default App;