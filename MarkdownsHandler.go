package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func listMarkdownsHandler(w http.ResponseWriter, r *http.Request) {
	type MarkdownFile struct {
		Category string `json:"category"`
		Name     string `json:"name"` // plik bez .md
		Path     string `json:"path"` // np. "LiblaryManagerApi/how_to_add_book"
	}

	var files []MarkdownFile

	root := "./markdowns"
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		parts := strings.Split(rel, string(filepath.Separator))
		if len(parts) == 2 {
			files = append(files, MarkdownFile{
				Category: parts[0],
				Name:     strings.TrimSuffix(parts[1], ".md"),
				Path:     strings.TrimSuffix(rel, ".md"),
			})
		}
		return nil
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func viewMarkdownHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	if strings.Contains(path, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Konwersja slasha do lokalnego separatora
	path = filepath.FromSlash(path)

	// Usuń ".md" z końca jeśli już jest
	path = strings.TrimSuffix(path, ".md")

	fullPath := filepath.Join("markdowns", path+".md")

	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "Cannot read file", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}
