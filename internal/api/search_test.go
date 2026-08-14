package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/martinushka/ios-rag/internal/search"
)

func TestSearchReturnsServiceResults(t *testing.T) {
	service := &fakeSearchService{
		results: []string{
			"Lenovo ThinkPad X1",
			"Lenovo IdeaPad 5",
		},
	}

	handler := NewHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/search",
		strings.NewReader(`{"query":"ноутбук Lenovo","limit":10}`),
	)

	recorder := httptest.NewRecorder()

	handler.Search(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	expected := `{"query":"ноутбук Lenovo","results":["Lenovo ThinkPad X1","Lenovo IdeaPad 5"]}` + "\n"

	if recorder.Body.String() != expected {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
}

type fakeSearchService struct {
	results []string
}

func (f *fakeSearchService) Search(
	ctx context.Context,
	query string,
	limit int,
) ([]string, error) {
	return f.results, nil
}

func newTestHandler() *Handler {
	return NewHandler(search.NewInMemoryService(nil))
}

func TestSearch(t *testing.T) {
	body := `{"query":"ноутбук Lenovo","limit":10}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/search",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	newTestHandler().Search(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestSearchRejectsInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/search",
		strings.NewReader(`{"query":`),
	)

	recorder := httptest.NewRecorder()

	newTestHandler().Search(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestSearchRequiresQuery(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/search",
		strings.NewReader(`{"limit":10}`),
	)

	recorder := httptest.NewRecorder()

	newTestHandler().Search(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestSearchRejectsWrongMethod(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/search",
		nil,
	)

	recorder := httptest.NewRecorder()

	newTestHandler().Search(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
	}
}
