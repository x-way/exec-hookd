package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	Port     int
	HookList []Hook
}

type Hook struct {
	Path string
	Exec []Exec
}

type Exec struct {
	Cmd     string
	Args    []string
	Timeout Duration
}

type Duration struct {
	time.Duration
}

func (duration *Duration) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON interface{}

	err := json.Unmarshal(b, &unmarshalledJSON)
	if err != nil {
		return err
	}

	switch value := unmarshalledJSON.(type) {
	case float64:
		duration.Duration = time.Duration(value)
	case string:
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration: %#v", unmarshalledJSON)
	}

	return nil
}

func runExec(e Exec) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout.Duration)
	defer cancel()

	return exec.CommandContext(ctx, e.Cmd, e.Args...).Run()
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	for _, hook := range cfg.HookList {
		if hook.Path == r.URL.Path {
			for _, exec := range hook.Exec {
				if err := runExec(exec); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			http.Error(w, "OK", http.StatusOK)
			return
		}
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}

func loadConfig(configfile string) {
	data, err := os.ReadFile(configfile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}
}

var cfg Config

func main() {
	configfile := flag.String("f", "exec-hookd.json", "path of the config file to use")
	flag.Parse()

	loadConfig(*configfile)

	http.HandleFunc("/", requestHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}
