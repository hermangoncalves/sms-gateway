package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	Listen string
	APIKey string
}

type SendSMSRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func (s *SendSMSRequest) isValid() bool {
	if s.To == "" || s.Message == "" {
		return false
	}

	return true
}

func main() {
	configPath := flag.String("config", "", "Path to JSON config file (optional)")
	listen := flag.String("listen", ":8081", "HTTP listen address (overrides config)")

	flag.Parse()

	cfg := Config{
		Listen: *listen,
	}

	if *configPath != "" {
		log.Printf("[DEBUG] Loading config file: %s", *configPath)
		if err := loadJSON(*configPath, &cfg); err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
	}

	if *listen != "" {
		cfg.Listen = *listen
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /send-sms", SendSmsHandler)

	server := &http.Server{
		Addr:         cfg.Listen,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("[INFO] CaddySMS v0.1 listening on %s", cfg.Listen)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

type Response map[string]any

func writeJson(w http.ResponseWriter, code int, obj any) {
	w.WriteHeader(code)
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		log.Printf("[ERROR] failed to encode JSON response: %v", err)
	}
}

func loadJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func SendSmsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] Incoming request: %s %s", r.Method, r.URL.Path)
	var req SendSMSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[WARN] Invalid JSON body: %v", err)
		writeJson(w, http.StatusBadRequest, Response{
			"error": "Invalid json",
		})
		return
	}

	log.Printf("[DEBUG] Parsed request: to=%s message=%s", req.To, req.Message)
	if !req.isValid() {
		log.Printf("[WARN] Invalid request: missing 'to' or 'message'")
		writeJson(w, http.StatusBadRequest, Response{
			"error": "'to' and 'message' are required",
		})
		return
	}

	if err := sendSMS(req.To, req.Message); err != nil {
		log.Printf("[ERROR] Failed to send SMS: %v", err)
		writeJson(w, http.StatusInternalServerError, Response{
			"error": "Internal Server error",
		})
		return
	}

	log.Printf("[INFO] SMS successfully sent to %s", req.To)
	writeJson(w, http.StatusOK, Response{
		"status": "sent",
		"to":     req.To,
	})
}

func haveInPath(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

func sendSMS(to, message string) error {
	log.Printf("[DEBUG] Preparing to send SMS to %s", to)
	if !haveInPath("termux-sms-send") {
		return fmt.Errorf("termux-sms-send not found; install termux-api and grant permissions")
	}

	cmd := exec.Command("termux-sms-send", "-n", to, message)
	log.Printf("[DEBUG] Executing command: %s", strings.Join(cmd.Args, " "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute termux-sms-send: %v | %s", err, strings.TrimSpace(string(out)))
	}

	log.Printf("[DEBUG] Command output: %s", strings.TrimSpace(string(out)))
	return nil
}
