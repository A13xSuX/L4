package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func sendRequest(server string, req CutRequest) (*CutResponse, error) {
	reqByte, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	address := "http://" + server + "/cut"

	resp, err := http.Post(address, "application/json", bytes.NewReader(reqByte))
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %w", resp.Status)
	}
	var cutResp CutResponse
	err = json.NewDecoder(resp.Body).Decode(&cutResp)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	return &cutResp, nil
}
