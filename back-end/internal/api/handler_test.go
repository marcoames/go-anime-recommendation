package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
)

type fakeAnimeRepository struct {
	data           []anime.Anime
	loadErr        error
	suggestions    []string
	suggestionsErr error
}

func (f fakeAnimeRepository) LoadAnimeData() ([]anime.Anime, error) {
	return f.data, f.loadErr
}

func (f fakeAnimeRepository) GetAnimeIndex(title string, data []anime.Anime) (int, error) {
	for index, item := range data {
		if item.Title == title {
			return index, nil
		}
	}
	return -1, errors.New("anime not found")
}

func (f fakeAnimeRepository) SearchAnimeTitles(string, int64) ([]string, error) {
	return f.suggestions, f.suggestionsErr
}

func testAnimeData() []anime.Anime {
	return []anime.Anime{
		{Title: "Naruto", Score: 8, Popularity: 10},
		{Title: "Bleach", Score: 7.5, Popularity: 20},
	}
}

func TestHandleRequestReturnsRecommendations(t *testing.T) {
	handler := NewHandler(fakeAnimeRepository{data: testAnimeData()})
	request := httptest.NewRequest(http.MethodGet, "/api/?anime=Naruto", nil)
	recorder := httptest.NewRecorder()

	handler.HandleRequest(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected JSON content type, got %q", recorder.Header().Get("Content-Type"))
	}
	if recorder.Body.Len() == 0 {
		t.Fatal("expected a JSON response body")
	}
}

func TestHandleRequestValidatesAnimeQuery(t *testing.T) {
	handler := NewHandler(fakeAnimeRepository{data: testAnimeData()})
	request := httptest.NewRequest(http.MethodGet, "/api/", nil)
	recorder := httptest.NewRecorder()

	handler.HandleRequest(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestHandleRequestReturnsNotFoundForUnknownAnime(t *testing.T) {
	handler := NewHandler(fakeAnimeRepository{data: testAnimeData()})
	request := httptest.NewRequest(http.MethodGet, "/api/?anime=Unknown", nil)
	recorder := httptest.NewRecorder()

	handler.HandleRequest(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestHandleRequestReturnsInternalServerErrorWhenRepositoryFails(t *testing.T) {
	handler := NewHandler(fakeAnimeRepository{loadErr: errors.New("database unavailable")})
	request := httptest.NewRequest(http.MethodGet, "/api/?anime=Naruto", nil)
	recorder := httptest.NewRecorder()

	handler.HandleRequest(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestHandleSuggestionsReturnsTitles(t *testing.T) {
	handler := NewHandler(fakeAnimeRepository{
		data:        testAnimeData(),
		suggestions: []string{"Naruto", "Naruto: Shippuden"},
	})
	request := httptest.NewRequest(http.MethodGet, "/api/suggestions?q=nar", nil)
	recorder := httptest.NewRecorder()

	handler.HandleSuggestions(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if recorder.Body.String() != `{"suggestions":["Naruto","Naruto: Shippuden"]}`+"\n" {
		t.Fatalf("unexpected suggestions response: %s", recorder.Body.String())
	}
}

func TestHandleSuggestionsReturnsInternalServerErrorWhenRepositoryFails(t *testing.T) {
	handler := NewHandler(fakeAnimeRepository{
		suggestionsErr: errors.New("database unavailable"),
	})
	request := httptest.NewRequest(http.MethodGet, "/api/suggestions?q=nar", nil)
	recorder := httptest.NewRecorder()

	handler.HandleSuggestions(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestHandleRequestRejectsNonGet(t *testing.T) {
	handler := NewHandler(fakeAnimeRepository{data: testAnimeData()})
	request := httptest.NewRequest(http.MethodPost, "/api/?anime=Naruto", nil)
	recorder := httptest.NewRecorder()

	handler.HandleRequest(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", recorder.Code)
	}
}
