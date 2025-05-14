package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

type APIDoc struct {
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
	if err := loadDocs(); err != nil {
		fmt.Println("Błąd:", err)
		return
	}

	fmt.Println("server started on port 5899")

	http.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(docs)
	})

	http.HandleFunc("/api/simulate", simulateHandler)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.ListenAndServe(":5899", nil)
}
