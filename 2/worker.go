package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
)

func (h *WorkerHandler) cutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CutRequest

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fields, err := parseFields(req.Fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := os.Open(h.dataFile)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	result := processLinesConcurrent(lines, fields, req.Delimiter, req.Separated)
	resp := CutResponse{
		Lines: result,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *WorkerHandler) StartServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/cut", h.cutHandler)
	return http.ListenAndServe(h.addr, mux)
}
