package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type SimulateRequest struct {
	Endpoint string            `json:"endpoint"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
	DB_entry []DBEntry         `json:"defaultDB"`
}

func simulateHandler(w http.ResponseWriter, r *http.Request) {
	var simReq SimulateRequest
	if err := json.NewDecoder(r.Body).Decode(&simReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var found *APIDoc
	for _, doc := range docs {
		if doc.Endpoint == simReq.Endpoint {
			found = &doc
			break
		}
	}

	if found == nil {
		http.Error(w, "Endpoint not found", http.StatusNotFound)
		return
	}

	luaScript := "./scripts/" + strings.Split(found.LuaFunc, "(")[0] + ".lua"

	result, err := runLuaFunctionWithContext(luaScript, found.LuaFunc, simReq, simReq.DB_entry)
	if err != nil {
		http.Error(w, "Lua error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
}
