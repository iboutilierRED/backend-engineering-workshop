package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWikipediaChangeUnmarshal tests JSON parsing of Wikipedia change events
func TestWikipediaChangeUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    WikipediaChange
		wantErr bool
	}{
		{
			name:    "valid bot change",
			jsonStr: `{"user":"BotName","uri":"https://en.wikipedia.org/wiki/Test","bot":true,"meta":{"uri":"https://en.wikipedia.org/wiki/Test_v1"}}`,
			want: WikipediaChange{
				User: "BotName",
				Uri:  "https://en.wikipedia.org/wiki/Test",
				Bot:  true,
				Meta: Meta{Uri: "https://en.wikipedia.org/wiki/Test_v1"},
			},
			wantErr: false,
		},
		{
			name:    "valid human change",
			jsonStr: `{"user":"HumanEditor","uri":"https://en.wikipedia.org/wiki/Physics","bot":false,"meta":{"uri":"https://en.wikipedia.org/wiki/Physics_v2"}}`,
			want: WikipediaChange{
				User: "HumanEditor",
				Uri:  "https://en.wikipedia.org/wiki/Physics",
				Bot:  false,
				Meta: Meta{Uri: "https://en.wikipedia.org/wiki/Physics_v2"},
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			jsonStr: `not valid json`,
			wantErr: true,
		},
		{
			name:    "empty json",
			jsonStr: `{}`,
			want: WikipediaChange{
				User: "",
				Uri:  "",
				Bot:  false,
				Meta: Meta{Uri: ""},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got WikipediaChange
			err := json.Unmarshal([]byte(tt.jsonStr), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("Unmarshal() got %+v, want %+v", got, tt.want)
			}
		})
	}
}

// TestStatsCollectorRecordChange tests the RecordChange method
func TestStatsCollectorRecordChange(t *testing.T) {
	sc := &StatsCollector{
		UniqueUsers: make(map[string]struct{}),
		UniqueUris:  make(map[string]struct{}),
	}

	tests := []struct {
		name        string
		change      WikipediaChange
		wantBots    int
		wantNonBots int
	}{
		{
			name: "record bot change",
			change: WikipediaChange{
				User: "BotName",
				Uri:  "https://en.wikipedia.org/wiki/Test",
				Bot:  true,
				Meta: Meta{Uri: "https://en.wikipedia.org/wiki/Test_v1"},
			},
			wantBots:    1,
			wantNonBots: 0,
		},
		{
			name: "record human change",
			change: WikipediaChange{
				User: "HumanEditor",
				Uri:  "https://en.wikipedia.org/wiki/Physics",
				Bot:  false,
				Meta: Meta{Uri: "https://en.wikipedia.org/wiki/Physics_v1"},
			},
			wantBots:    0,
			wantNonBots: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc.RecordChange(tt.change)
		})
	}

	stats := sc.GetStats()
	if stats.Messages != 2 {
		t.Errorf("Messages: got %d, want 2", stats.Messages)
	}
	if stats.UniqueUsers != 2 {
		t.Errorf("UniqueUsers: got %d, want 2", stats.UniqueUsers)
	}
	if stats.UniqueUris != 2 {
		t.Errorf("UniqueUris: got %d, want 2", stats.UniqueUris)
	}
	if stats.Bots != 1 {
		t.Errorf("Bots: got %d, want 1", stats.Bots)
	}
	if stats.NonBots != 1 {
		t.Errorf("NonBots: got %d, want 1", stats.NonBots)
	}
}

// TestStatsCollectorGetStats tests the GetStats method
func TestStatsCollectorGetStats(t *testing.T) {
	sc := &StatsCollector{
		MessageCount: 30,
		UniqueUsers: map[string]struct{}{
			"alice": {},
			"bob":   {},
		},
		UniqueUris: map[string]struct{}{
			"https://en.wikipedia.org/wiki/Go":          {},
			"https://en.wikipedia.org/wiki/Concurrency": {},
		},
		Bots:    10,
		NonBots: 20,
	}

	stats := sc.GetStats()

	want := Stats{
		Messages:    30,
		UniqueUsers: 2,
		UniqueUris:  2,
		Bots:        10,
		NonBots:     20,
	}

	if stats != want {
		t.Errorf("GetStats() got %+v, want %+v", stats, want)
	}
}

// TestStatsHandler tests the statsHandler HTTP handler
func TestStatsHandler(t *testing.T) {
	sc := &StatsCollector{
		MessageCount: 30,
		UniqueUsers: map[string]struct{}{
			"alice": {},
			"bob":   {},
		},
		UniqueUris: map[string]struct{}{
			"https://en.wikipedia.org/wiki/Go":          {},
			"https://en.wikipedia.org/wiki/Concurrency": {},
		},
		Bots:    10,
		NonBots: 20,
	}

	handler := statsHandler(sc)
	req := httptest.NewRequest("GET", "/stats", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want application/json", ct)
	}

	var got Stats
	err := json.Unmarshal(rr.Body.Bytes(), &got)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	want := Stats{
		Messages:    30,
		UniqueUsers: 2,
		UniqueUris:  2,
		Bots:        10,
		NonBots:     20,
	}

	if got != want {
		t.Errorf("handler returned unexpected body: got %+v want %+v", got, want)
	}
}

// TestStatsStructure tests the Stats struct JSON marshaling
func TestStatsStructure(t *testing.T) {
	stats := Stats{
		Messages:    100,
		UniqueUsers: 50,
		UniqueUris:  75,
		Bots:        25,
		NonBots:     75,
	}

	data, err := json.Marshal(stats)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	expectedKeys := []string{
		"messages_consumed",
		"distinct_users",
		"distinct_uris",
		"num_bots",
		"num_non_bots",
	}

	for _, key := range expectedKeys {
		if _, exists := result[key]; !exists {
			t.Errorf("expected key %q not found in JSON output", key)
		}
	}

	if int(result["messages_consumed"].(float64)) != 100 {
		t.Errorf("messages_consumed: got %v, want 100", result["messages_consumed"])
	}
	if int(result["distinct_users"].(float64)) != 50 {
		t.Errorf("distinct_users: got %v, want 50", result["distinct_users"])
	}
}
