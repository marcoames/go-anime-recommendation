import React, { useState } from 'react';
import './index.css';

function App() {
  const [animeTitle, setAnimeTitle] = useState('');
  const [animeDetails, setAnimeDetails] = useState(null);
  const [recommendations, setRecommendations] = useState([]);
  const [error, setError] = useState('');

  const handleSearch = async () => {
    try {
      // Send request to the Go backend
      //const response = await fetch(`https://animerecommendation-443313.rj.r.appspot.com/api/?anime=${animeTitle}`);
      const response = await fetch(`http://localhost:8080/api/?anime=${animeTitle}`);
      // If the response is not OK, throw an error
      if (!response.ok) {
        throw new Error('Anime not found');
      }

      // Parse the JSON response
      const data = await response.json();

      console.log(data);

      // Set anime details (e.g., title, image, genres)
      setAnimeDetails(data.anime);

      // Set recommendations (assuming the data contains a 'recommendations' array)
      setRecommendations(data.recommendations || []);

      // Clear any previous errors
      setError('');
    } catch (error) {
      setAnimeDetails(null);
      setRecommendations([]);
      setError(error.message);
    }
  };

  return (
    <div className="container">
      <h1>Anime Recommendations</h1>

      <input
        type="text"
        value={animeTitle}
        onChange={(e) => setAnimeTitle(e.target.value)}
        placeholder="Enter an anime title"
      />
      <button onClick={handleSearch}>Get Recommendations</button>

      {error && <div className="error">{error}</div>}

      {/* Display anime details */}
      {animeDetails && (
        <div className="anime-details">
          <h2>{animeDetails.title_english || animeDetails.title}</h2>
          <img src={animeDetails.images.jpg.image_url} alt={animeDetails.title_english || animeDetails.title} />
          <p><strong>Genres:</strong> {animeDetails.genres?.map((genre) => genre.name).join(', ')}</p>
          <p><strong>Rating:</strong> {animeDetails.rating}</p>
          <p><strong>Status:</strong> {animeDetails.status}</p>
          <p><strong>Score:</strong> {animeDetails.score}</p>
          <p><strong>MAL Page:</strong> <a href={animeDetails.url} target="_blank" rel="noopener noreferrer">Link to MyAnimeList</a></p>
          <p><strong>Trailer:</strong> <a href={animeDetails.trailer.embed_url} target="_blank" rel="noopener noreferrer">Watch Trailer</a></p>
          <p><strong>Synopsis:</strong> {animeDetails.synopsis}</p>
        </div>
      )}


{/* Display recommendations */}
{recommendations.length > 0 && (
  <div className="recommendations">
    <h3>Recommendations:</h3>
    
    {recommendations.map((anime, index) => (
      <li key={index} className="anime-item">
        <div className="anime-info">
          <div className="image-container">
            <img src={anime.images.jpg.image_url} alt={anime.title} />
          </div>
          <div className="details">
            <h2>{anime.title}</h2>
            <p><strong>Genres:</strong> {anime.genres?.map((genre) => genre.name).join(', ')}</p>
            <p><strong>Rating:</strong> {anime.rating}</p>
            <p><strong>Status:</strong> {anime.status}</p>
            <p><strong>Score:</strong> {anime.score}</p>
            
            {/* Link to MyAnimeList */}
            <p>
              <strong>MAL Page: </strong> 
              <a href={anime.url} target="_blank" rel="noopener noreferrer">
                Link to MyAnimeList
              </a>
            </p>
            
            {/* Trailer (if available) */}
            <p>
              <strong>Trailer: </strong> 
              {anime.trailer?.embed_url ? (
                <a href={anime.trailer.embed_url} target="_blank" rel="noopener noreferrer">
                  Watch Trailer
                </a>
              ) : (
                <span>No trailer available</span>
              )}
            </p>
            
            {/* Synopsis */}
            <p><strong>Synopsis:</strong> {anime.synopsis}</p>
          </div>
        </div>
      </li>
    ))}
  </div>
)}



    </div>
  );
}

export default App;
