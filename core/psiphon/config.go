package psiphon

import (
	"encoding/json"
	"errors"
	"sync"

	"com.argotunnel/core/capabilities"
)

type Config struct {
	mu                 sync.RWMutex
	PropagationChannel string                     `json:"propagation_channel"`
	SponsorID          string                     `json:"sponsor_id"`
	FrontDomains       []string                   `json:"front_domains"`
	HostHeaders        []string                   `json:"host_headers"`
	TargetSNI          string                     `json:"target_sni"`
	FallbackLevel      capabilities.TransportType `json:"fallback_level"`
}

func NewConfig(jsonData []byte, currentFallback capabilities.TransportType) (*Config, error) {
	var cfg Config
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		return nil, errors.New("failed to parse psiphon configuration: " + err.Error())
	}

	cfg.FallbackLevel = currentFallback

	// Ensure defaults for severe censorship environments
	if len(cfg.FrontDomains) == 0 {
		cfg.FrontDomains = []string{"ajax.aspnetcdn.com", "cdnjs.cloudflare.com"}
	}
	if len(cfg.HostHeaders) == 0 {
		cfg.HostHeaders = []string{"req.psiphon.net"}
	}
	if cfg.TargetSNI == "" {
		cfg.TargetSNI = "random.sni.psiphon.net"
	}

	return &cfg, nil
}

func (c *Config) GetActiveFrontingSpec() (string, string, string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Deterministic selection for Phase 3
	return c.FrontDomains[0], c.HostHeaders[0], c.TargetSNI
}
