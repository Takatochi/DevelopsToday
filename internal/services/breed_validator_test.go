package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBreedValidator_IsValid(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"name":"Bengal"},{"name":"Siamese"}]`))
	}))
	defer mockServer.Close()

	// Create validator with custom URL (no race condition!)
	validator := NewBreedWithURL(mockServer.URL)

	// Test cases
	tests := []struct {
		name      string
		breedName string
		want      bool
	}{
		{
			name:      "Valid breed - exact match",
			breedName: "Bengal",
			want:      true,
		},
		{
			name:      "Valid breed - case insensitive",
			breedName: "bengal",
			want:      true,
		},
		{
			name:      "Invalid breed",
			breedName: "NonExistentBreed",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validator.IsValid(tt.breedName); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBreedValidator_FetchBreeds_Error(t *testing.T) {
	// Setup mock server that returns an error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Create validator with custom URL (no race condition!)
	validator := NewBreedWithURL(mockServer.URL)

	// Test that IsValid returns false when API call fails
	if validator.IsValid("Bengal") {
		t.Errorf("IsValid() should return false when API call fails")
	}
}
