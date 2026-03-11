package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/iboutilierRED/backend-workshop/ch-3/wikievents/internal/data"
	"github.com/iboutilierRED/backend-workshop/ch-3/wikievents/internal/handlers"
	"github.com/iboutilierRED/backend-workshop/ch-3/wikievents/internal/middleware"
	"github.com/iboutilierRED/backend-workshop/ch-3/wikievents/internal/models"
	"github.com/iboutilierRED/backend-workshop/ch-3/wikievents/internal/stats"
)

func main() {
	c, err := loadConfigs()

	if err != nil {
		log.Fatalf("Error loading configs file: %v", err)
	}

	// Accept the port as an overide argument. Default is 7000 if no argument is provided.
	if len(os.Args) > 1 {
		c.StatsPort = os.Args[1]
	}

	_, err = data.InitDbConnection()

	if err != nil {
		log.Fatalf("Error communicating with database: %v", err)
	}

	messages := make(chan models.WikipediaChange)

	statsCollector := &stats.StatsCollector{
		UniqueUsers: make(map[string]struct{}),
		UniqueUris:  make(map[string]struct{}),
	}

	go consumeWikipedia(messages, c)

	go updateStats(statsCollector, messages)

	http.HandleFunc("/signup", handlers.SignUpHandler())
	http.HandleFunc("/login", handlers.LoginHandler())
	http.HandleFunc("/stats", middleware.Authenticate(handlers.StatsHandler(statsCollector)))

	log.Fatal(http.ListenAndServe(":"+c.StatsPort, nil))
}

// This function connects to the Wikipedia stream and sends messages to the channel.
func consumeWikipedia(ch chan<- models.WikipediaChange, c models.AppConfig) {
	// adding defer close channel to ensure the channel gets closed if Wikipedia stream connection is lost
	defer close(ch)
	prefix := "data: "
	client := &http.Client{}

	req, err := http.NewRequest("GET", c.FeedUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "WorkshopBot/1.0 (https://github.com/iboutilierRED/backend-workshop; ian.boutilier@redspace.com)")

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
		var message models.WikipediaChange
		line := scanner.Text()

		if strings.HasPrefix(line, prefix) {
			jsonData, _ := strings.CutPrefix(line, prefix)
			err := json.Unmarshal([]byte(jsonData), &message)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v", err)
			} else {
				ch <- message
				// Commenting out console log to reduce noise
				//log.Printf("User: %s, URI: %s, Bot: %t", message.User, message.Meta.Uri, message.Bot)
			}
		}
	}
}

// This function runs as a go routine and updates the stats based on the messages received from the channel. It locks the statsMutex to safely update shared state.
func updateStats(sc *stats.StatsCollector, ch <-chan models.WikipediaChange) {
	for message := range ch {
		sc.RecordChange(message)
	}
}

func loadConfigs() (models.AppConfig, error) {
	var config models.AppConfig
	f, err := os.Open("internal/config/config.json")

	if err != nil {
		return config, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	err = dec.Decode(&config)
	return config, err
}
