package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func runLocal(cfg *Config) {
	fields, err := parseFields(cfg.Fields)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
		os.Exit(1)
	}

	results := processLinesConcurrent(lines, fields, cfg.Delimiter, cfg.Separated)
	for _, line := range results {
		fmt.Println(line)
	}
}

func runWorker(cfg *Config) {
	if cfg.File == "" {
		fmt.Fprintln(os.Stderr, "worker mode requires -file")
		os.Exit(1)
	}

	handler := WorkerHandler{
		addr:     cfg.Addr,
		dataFile: cfg.File,
	}

	if err := handler.StartServer(); err != nil {
		fmt.Fprintln(os.Stderr, "worker error:", err)
		os.Exit(1)
	}
}

func runClient(cfg *Config) {
	if cfg.Servers == "" {
		fmt.Fprintln(os.Stderr, "client mode requires -servers")
		os.Exit(1)
	}

	servers := strings.Split(cfg.Servers, ",")
	for i := range servers {
		servers[i] = strings.TrimSpace(servers[i])
	}

	quorum := cfg.Quorum
	if quorum <= 0 {
		quorum = len(servers)/2 + 1
	}

	req := CutRequest{
		Fields:    cfg.Fields,
		Delimiter: cfg.Delimiter,
		Separated: cfg.Separated,
	}

	respCh := make(chan *serverResult, len(servers))
	errCh := make(chan error, len(servers))

	var wg sync.WaitGroup
	for i, server := range servers {
		wg.Add(1)
		go func(i int, server string) {
			defer wg.Done()

			cutResp, err := sendRequest(server, req)
			if err != nil {
				errCh <- err
				return
			}

			respCh <- &serverResult{
				index: i,
				resp:  cutResp,
			}
		}(i, server)
	}

	wg.Wait()
	close(respCh)
	close(errCh)

	ordered := make([]*CutResponse, len(servers))
	for sr := range respCh {
		ordered[sr.index] = sr.resp
	}

	successResp := 0
	for _, resp := range ordered {
		if resp != nil {
			successResp++
		}
	}

	for err := range errCh {
		fmt.Fprintln(os.Stderr, err)
	}

	if successResp < quorum {
		fmt.Fprintln(os.Stderr, "quorum not reached:", successResp, "/", quorum)
		os.Exit(1)
	}

	for _, resp := range ordered {
		if resp == nil {
			continue
		}
		for _, line := range resp.Lines {
			fmt.Println(line)
		}
	}
}

func main() {
	config := Parse()

	switch config.Mode {
	case "local":
		runLocal(config)
	case "worker":
		runWorker(config)
	case "client":
		runClient(config)
	default:
		fmt.Fprintln(os.Stderr, "unknown mode:", config.Mode)
		os.Exit(1)
	}
}
