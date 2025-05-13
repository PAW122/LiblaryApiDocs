package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var docs []APIDoc

func loadDocs() error {
	file, err := os.ReadFile("./docs/docs.json")
	if err != nil {
		return fmt.Errorf("błąd odczytu pliku: %w", err)
	}

	if err := json.Unmarshal(file, &docs); err != nil {
		return fmt.Errorf("błąd parsowania JSON: %w", err)
	}

	return nil
}
