package main

import "flag"

type Config struct {
	Mode      string
	Fields    string
	Delimiter string
	Separated bool
	Addr      string
	File      string
	Servers   string
	Quorum    int
}

func Parse() *Config {
	var cfg Config

	flag.StringVar(&cfg.Mode, "mode", "local", "Mode: local | worker | client")
	flag.StringVar(&cfg.Fields, "f", "", "Fields to use in the input")
	flag.StringVar(&cfg.Delimiter, "d", "\t", "Delimiter to use in the input")
	flag.BoolVar(&cfg.Separated, "s", false, "Only output lines containing the delimiter")
	flag.StringVar(&cfg.Addr, "addr", "localhost:9001", "Address for worker server")
	flag.StringVar(&cfg.File, "file", "", "Data file for worker")
	flag.StringVar(&cfg.Servers, "servers", "", "Comma-separated worker addresses for client")
	flag.IntVar(&cfg.Quorum, "quorum", 0, "Quorum size for client")

	flag.Parse()
	return &cfg
}
