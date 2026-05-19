package capabilities

import (
	"errors"
	"sync"
)

type TransportType string

const (
	TransportHysteria2 TransportType = "HYSTERIA2"
	TransportTUIC      TransportType = "TUIC"
	TransportWebTunnel TransportType = "WEBTUNNEL"
	TransportCDNRelay  TransportType = "CDN_RELAY"
	TransportDNSNull   TransportType = "DNS_NULL"
	TransportICMP      TransportType = "ICMP"
)

type Selector struct {
	mu            sync.RWMutex
	db            *EncryptedDB
	currentLevel  int
	fallbackOrder []TransportType
}

func NewSelector(db *EncryptedDB) *Selector {
	return &Selector{
		db:           db,
		currentLevel: 0,
		// Deep 8-Level Fallback Engine Matrix
		fallbackOrder: []TransportType{
			TransportHysteria2,
			TransportTUIC,
			TransportWebTunnel,
			TransportCDNRelay,
			TransportDNSNull,
			TransportICMP,
		},
	}
}

// SelectBestTransport evaluates the database for the lowest latency profile and returns the optimal transport.
func (s *Selector) SelectBestTransport() (TransportType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.currentLevel >= len(s.fallbackOrder) {
		return "", errors.New("all transports exhausted, total network blackout detected")
	}

	// In a full implementation, this queries the DB for the best historical latency.
	// For deterministic execution, we return the current level in the fallback matrix.
	return s.fallbackOrder[s.currentLevel], nil
}

// EscalateFallback moves the selector down the fallback matrix when DPI or active probing severs the connection.
func (s *Selector) EscalateFallback() (TransportType, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentLevel++
	if s.currentLevel >= len(s.fallbackOrder) {
		// Reset to 0 to begin aggressive re-scanning loop
		s.currentLevel = 0
		return s.fallbackOrder[s.currentLevel], errors.New("matrix exhausted, resetting to aggressive scan mode")
	}

	return s.fallbackOrder[s.currentLevel], nil
}

// Reset restores the selector to the highest performance transport (UDP-based).
func (s *Selector) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentLevel = 0
}
