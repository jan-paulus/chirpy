package main

import "net/http"

func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
}
