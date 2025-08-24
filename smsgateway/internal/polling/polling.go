package polling

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hermangoncalves/sms-gateway/internal/config"
	"github.com/hermangoncalves/sms-gateway/internal/sms"
)

type Poller struct {
	cfg    *config.Config
	client *http.Client
}

func NewPoller(cfg *config.Config) *Poller {
	return &Poller{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *Poller) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(p.cfg.PollingInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Poller stopped")
			return
		case <-ticker.C:
			messages, err := p.fetchMessages()
			if err != nil {
				log.Printf("Failed to fetch messages: %v", err)
				continue
			}

			for _, msg := range messages {
				if err := p.processMessage(msg); err != nil {
					log.Printf("Failed to process message for %s: %v", msg.Number, err)
				}
			}
		}
	}
}

func (p *Poller) fetchMessages() ([]sms.SMSMessage, error) {
	resp, err := p.client.Get(p.cfg.WorkerURL)
	if err != nil {
		return nil, fmt.Errorf("worker request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("worker returned status %d: %s", resp.StatusCode, string(body))
	}

	var messages []sms.SMSMessage
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	return messages, nil
}

func (p *Poller) processMessage(msg sms.SMSMessage) error {
	if err := sms.SendSMS(msg.Number, msg.Text); err != nil {
		return fmt.Errorf("send failed: %w", err)
	}
	return p.confirmDelivery(msg.Number)
}

func (p *Poller) confirmDelivery(number string) error {
	payload, err := json.Marshal(map[string]string{
		"number": number,
		"status": "sent",
	})
	if err != nil {
		return fmt.Errorf("marshal confirm payload: %w", err)
	}

	resp, err := p.client.Post(p.cfg.ConfirmURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("confirm request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("confirm returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
