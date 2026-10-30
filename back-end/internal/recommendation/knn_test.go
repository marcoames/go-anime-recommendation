package recommendation

import "testing"

func TestFindRecommendationsRejectsInvalidInput(t *testing.T) {
	if got := FindRecommendations([][]float64{{1}}, -1, 1); got != nil {
		t.Fatalf("expected nil for invalid query index, got %v", got)
	}
	if got := FindRecommendations([][]float64{{1}}, 0, 0); got != nil {
		t.Fatalf("expected nil for non-positive k, got %v", got)
	}
}
