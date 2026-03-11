package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iboutilierRED/backend-workshop/ch-3/wikievents/internal/stats"
)

func StatsHandler(sc *stats.StatsCollector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		stats := sc.GetStats()
		json.NewEncoder(w).Encode(stats)
	}
}
