package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusRecorder(t *testing.T) {
	recorder := httptest.NewRecorder()

	w := &statusRecorder{
		ResponseWriter: recorder,
	}

	w.WriteHeader(http.StatusCreated)

	if w.status != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.status)
	}

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected response status %d, got %d", http.StatusCreated, recorder.Code)
	}
}

func TestStatusRecorderDefaultsToOK(t *testing.T) {
	recorder := httptest.NewRecorder()

	w := &statusRecorder{
		ResponseWriter: recorder,
	}

	_, err := w.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if w.status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.status)
	}
}
