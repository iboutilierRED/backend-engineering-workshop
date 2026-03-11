package stats

import (
	"sync"

	"github.com/iboutilierRED/backend-workshop/ch-3/wikievents/internal/models"
)

type StatsCollector struct {
	mu           sync.RWMutex
	MessageCount int
	UniqueUsers  map[string]struct{} // using blank struct instead of bool to save memory (bool = 1 byte, struct{} = 0 bytes)
	UniqueUris   map[string]struct{}
	Bots         int
	NonBots      int
}

// add messges to the StatsCollector
func (sc *StatsCollector) RecordChange(change models.WikipediaChange) {
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
func (sc *StatsCollector) GetStats() models.Stats {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return models.Stats{
		Messages:    sc.MessageCount,
		UniqueUsers: len(sc.UniqueUsers),
		UniqueUris:  len(sc.UniqueUris),
		Bots:        sc.Bots,
		NonBots:     sc.NonBots,
	}
}
