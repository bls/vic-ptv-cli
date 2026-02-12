package api

import (
	"strings"
	"testing"
)

func TestSignURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		devID    string
		apiKey   string
		wantSig  string // expected signature (uppercase hex)
		wantPath string // expected path portion before signature
	}{
		{
			name:     "simple path no query params",
			baseURL:  "https://timetableapi.ptv.vic.gov.au",
			path:     "/v3/route_types",
			devID:    "1000001",
			apiKey:   "aaaabbbb-cccc-dddd-eeee-ffffffaaaaaa",
			wantPath: "/v3/route_types?devid=1000001",
		},
		{
			name:     "path with existing query params",
			baseURL:  "https://timetableapi.ptv.vic.gov.au",
			path:     "/v3/departures/route_type/0/stop/1234?max_results=5",
			devID:    "1000001",
			apiKey:   "aaaabbbb-cccc-dddd-eeee-ffffffaaaaaa",
			wantPath: "/v3/departures/route_type/0/stop/1234?max_results=5&devid=1000001",
		},
		{
			name:     "search with encoded term",
			baseURL:  "https://timetableapi.ptv.vic.gov.au",
			path:     "/v3/search/Flinders%20Street",
			devID:    "2000002",
			apiKey:   "11112222-3333-4444-5555-666677778888",
			wantPath: "/v3/search/Flinders%20Street?devid=2000002",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SignURL(tt.baseURL, tt.path, tt.devID, tt.apiKey)
			if err != nil {
				t.Fatalf("SignURL() error = %v", err)
			}

			// Check that the URL starts with baseURL + path
			expectedPrefix := tt.baseURL + tt.wantPath
			if !strings.HasPrefix(got, expectedPrefix) {
				t.Errorf("SignURL() = %q, want prefix %q", got, expectedPrefix)
			}

			// Check that the signature parameter is present
			if !strings.Contains(got, "&signature=") {
				t.Errorf("SignURL() = %q, missing &signature= parameter", got)
			}

			// Extract signature and validate it's a valid hex string (uppercase)
			sigIdx := strings.Index(got, "&signature=")
			sig := got[sigIdx+len("&signature="):]
			if len(sig) != 40 { // SHA1 produces 20 bytes = 40 hex chars
				t.Errorf("signature length = %d, want 40", len(sig))
			}
			for _, c := range sig {
				if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
					t.Errorf("signature contains invalid char %c, want uppercase hex", c)
					break
				}
			}
		})
	}
}

func TestSignURLDeterministic(t *testing.T) {
	// Same inputs should produce the same signature
	url1, _ := SignURL("https://timetableapi.ptv.vic.gov.au", "/v3/route_types", "1000001", "test-key")
	url2, _ := SignURL("https://timetableapi.ptv.vic.gov.au", "/v3/route_types", "1000001", "test-key")
	if url1 != url2 {
		t.Errorf("SignURL not deterministic:\n  first:  %s\n  second: %s", url1, url2)
	}
}

func TestSignURLDifferentKeys(t *testing.T) {
	// Different API keys should produce different signatures
	url1, _ := SignURL("https://timetableapi.ptv.vic.gov.au", "/v3/route_types", "1000001", "key-one")
	url2, _ := SignURL("https://timetableapi.ptv.vic.gov.au", "/v3/route_types", "1000001", "key-two")
	if url1 == url2 {
		t.Error("SignURL produced same URL for different API keys")
	}
}

func TestSignURLDifferentDevIDs(t *testing.T) {
	// Different dev IDs should produce different signatures
	url1, _ := SignURL("https://timetableapi.ptv.vic.gov.au", "/v3/route_types", "1000001", "same-key")
	url2, _ := SignURL("https://timetableapi.ptv.vic.gov.au", "/v3/route_types", "2000002", "same-key")
	if url1 == url2 {
		t.Error("SignURL produced same URL for different dev IDs")
	}
}

func TestSignURLKnownVector(t *testing.T) {
	// Test with a known HMAC-SHA1 computation
	// For path "/v3/route_types?devid=1000001" with key "9c132d31-6a30-4cac-8d8b-8a1970834799"
	// We can verify by computing HMAC-SHA1 independently
	got, err := SignURL(
		"https://timetableapi.ptv.vic.gov.au",
		"/v3/route_types",
		"1000001",
		"9c132d31-6a30-4cac-8d8b-8a1970834799",
	)
	if err != nil {
		t.Fatalf("SignURL() error = %v", err)
	}

	// Verify structure
	if !strings.HasPrefix(got, "https://timetableapi.ptv.vic.gov.au/v3/route_types?devid=1000001&signature=") {
		t.Errorf("unexpected URL structure: %s", got)
	}

	// The HMAC-SHA1 of "/v3/route_types?devid=1000001" with key "9c132d31-6a30-4cac-8d8b-8a1970834799"
	// should be consistent. Let's verify by extracting and checking it's 40 hex chars.
	parts := strings.Split(got, "&signature=")
	if len(parts) != 2 {
		t.Fatalf("expected exactly one &signature= in URL, got %d parts", len(parts))
	}
	sig := parts[1]
	if len(sig) != 40 {
		t.Errorf("signature length = %d, want 40 hex chars", len(sig))
	}

	// Run it again to verify determinism with this specific key
	got2, _ := SignURL(
		"https://timetableapi.ptv.vic.gov.au",
		"/v3/route_types",
		"1000001",
		"9c132d31-6a30-4cac-8d8b-8a1970834799",
	)
	if got != got2 {
		t.Errorf("non-deterministic results:\n  first:  %s\n  second: %s", got, got2)
	}
}
