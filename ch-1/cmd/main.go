package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var listenPort string = ":7000"

type WikipediaChange struct {
	User string `json:"user"`
	Uri  string `json:"uri"`
	Bot  bool   `json:"bot"`
	Meta Meta   `json:"meta"`
}

type Meta struct {
	Uri string `json:"uri"`
}

type Stats struct {
	Messages    int `json:"messages_consumed"`
	UniqueUsers int `json:"distinct_users"`
	UniqueUris  int `json:"distinct_uris"`
	Bots        int `json:"num_bots"`
	NonBots     int `json:"num_non_bots"`
}

type StatsCollector struct {
	mu           sync.RWMutex
	MessageCount int
	UniqueUsers  map[string]struct{} // using blank struct instead of bool to save memory (bool = 1 byte, struct{} = 0 bytes)
	UniqueUris   map[string]struct{}
	Bots         int
	NonBots      int
}

// add messges to the StatsCollector
func (sc *StatsCollector) RecordChange(change WikipediaChange) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.MessageCount++
	sc.UniqueUsers[change.User] = struct{}{}
	sc.UniqueUris[change.Meta.Uri] = struct{}{}

	if change.Bot {
		sc.Bots++
	} else {
		sc.NonBots++
	}
}

// Retrieve the stats as a snapshot from the statscollector
func (sc *StatsCollector) GetStats() Stats {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return Stats{
		Messages:    sc.MessageCount,
		UniqueUsers: len(sc.UniqueUsers),
		UniqueUris:  len(sc.UniqueUris),
		Bots:        sc.Bots,
		NonBots:     sc.NonBots,
	}
}

func statsHandler(sc *StatsCollector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		stats := sc.GetStats()
		json.NewEncoder(w).Encode(stats)
	}
}

func main() {
	// Accept the port as an argument. Default is 7000 if no argument is provided.
	if len(os.Args) > 1 {
		listenPort = ":" + os.Args[1]
	}

	messages := make(chan WikipediaChange, 100) // buffered channel to handle bursts of messages

	statsCollector := &StatsCollector{
		UniqueUsers: make(map[string]struct{}),
		UniqueUris:  make(map[string]struct{}),
	}

	go consumeWikipedia(messages)

	go updateStats(statsCollector, messages)

	http.HandleFunc("/stats", statsHandler(statsCollector))

	log.Fatal(http.ListenAndServe(listenPort, nil))
}

// This function connects to the Wikipedia stream and sends messages to the channel.
func consumeWikipedia(ch chan<- WikipediaChange) {
	// adding defer close channel to ensure the channel gets closed if Wikipedia stream connection is lost
	defer close(ch)
	prefix := "data: "
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://stream.wikimedia.org/v2/stream/recentchange", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "WorkshopBot/1.0 (https://github.com/iboutilier22/workshop; ian.boutilier@redspace.com)")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	// ensure response body gets closed when we're done to prevent resource leaks
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	log.Println("Connected to Wikipedia stream")

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		var message WikipediaChange
		line := scanner.Text()

		if strings.HasPrefix(line, prefix) {
			jsonData, _ := strings.CutPrefix(line, prefix)
			err := json.Unmarshal([]byte(jsonData), &message)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v", err)
			} else {
				ch <- message
				log.Printf("User: %s, URI: %s, Bot: %t", message.User, message.Meta.Uri, message.Bot)
			}
		}
	}
}

// This function runs as a go routine and updates the stats based on the messages received from the channel. It locks the statsMutex to safely update shared state.
func updateStats(sc *StatsCollector, ch <-chan WikipediaChange) {
	for message := range ch {
		sc.RecordChange(message)
	}
}
