package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type APIError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

type APIDoc struct {
	Markdown    []string                 `json:"markdown"`
	Method      string                   `json:"method"`
	Endpoint    string                   `json:"endpoint"`
	Description string                   `json:"description"`
	Permissions string                   `json:"permissions"`
	Body        string                   `json:"body"`
	Headers     string                   `json:"headers"`
	Response    string                   `json:"res"`
	Errors      []APIError               `json:"errors"`
	Category    string                   `json:"category"`
	LuaFunc     string                   `json:"luaFunc,omitempty"`
	DefaultDB   []map[string]interface{} `json:"defaultDB,omitempty"`
}

func main() {
	port := "5899"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	if err := loadDocs(); err != nil {
		fmt.Println("Błąd:", err)
		return
	}

	fmt.Println("server started on port", port)

	http.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(docs)
	})

	http.HandleFunc("/api/simulate", simulateHandler)
	http.HandleFunc("/api/markdowns", listMarkdownsHandler)
	http.HandleFunc("/api/markdowns/view", viewMarkdownHandler)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.ListenAndServe(":"+port, nil)
}
