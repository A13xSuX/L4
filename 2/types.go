package main

type serverResult struct {
	index int
	resp  *CutResponse
}
type WorkerHandler struct {
	addr     string
	dataFile string
}

type CutRequest struct {
	Fields    string `json:"fields"`
	Delimiter string `json:"delimiter"`
	Separated bool   `json:"separated"`
}

type CutResponse struct {
	Lines []string `json:"lines"`
}

type lineJob struct {
	index int
	text  string
}

type lineResult struct {
	index int
	text  string
}
