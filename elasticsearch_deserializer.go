package main

import (
	"context"
	"encoding/json"
	"net/url"
	"time"
)

type LogEntry struct {
	Timestamp    time.Time
	URL          url.URL
	ResponseTime time.Duration
	ReturnCode   int
	Method       string
	Body         json.RawMessage
}

func normalizeSearchResult(js json.RawMessage) LogEntry {
	entryMap := make(map[string]string)
	json.Unmarshal(js, entryMap)
}
func deserialize(ctx context.Context, entries chan<- LogEntry, hits <-chan json.RawMessage) error {
	for hit := range hits {
		// Deserialize
		entryMap := make(map[string]string)
		err := json.Unmarshal(hit, &e)
		if err != nil {
			return err
		}

		entries <- e

		// Terminate early?
		select {
		default:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
