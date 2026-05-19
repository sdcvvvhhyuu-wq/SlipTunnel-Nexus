package engine

import (
	"errors"
	"log"
	"sync"

	"com.argotunnel/core/capabilities"
	"com.argotunnel/core/psiphon"
)

// =====================================================================
// 8-LEVEL FALLBACK MATRIX ORCHESTRATOR
// =====================================================================

type FallbackEngine struct {
	mu             sync.Mutex
	selector       *capabilities.Selector
	psiphonCtrl    *psiphon.Controller
	routingManager *RoutingManager
}

func NewFallbackEngine(selector *capabilities.Selector, db *capabilities.EncryptedDB) *FallbackEngine {
	// Initialize Psiphon Config with default JSON
	cfgData := []byte(`{"propagation_channel": "default", "sponsor_id": "nexus"}`)
	cfg, _ := psiphon.NewConfig(cfgData, capabilities.TransportCDNRelay)

	return &FallbackEngine{
		selector:       selector,
		psiphonCtrl:    psiphon.NewController(cfg, db),
		routingManager: NewRoutingManager(),
	}
}

func (fe *FallbackEngine) ExecuteTransport() error {
	fe.mu.Lock()
	defer fe.mu.Unlock()

	transport, err := fe.selector.SelectBestTransport()
	if err != nil {
		return err
	}

	log.Printf("Executing Transport Level: %s", transport)

	switch transport {
	case capabilities.TransportHysteria2, capabilities.TransportTUIC, capabilities.TransportWebTunnel:
		// High-speed UDP/TCP transports executed here
		// If they fail, we escalate downward
		_, err = fe.selector.EscalateFallback()
		return err

	case capabilities.TransportCDNRelay:
		// Execute Psiphon Meek/Domain Fronting as the mid-tier fallback
		err := fe.psiphonCtrl.StartTunnel()
		if err != nil {
			_, _ = fe.selector.EscalateFallback()
			return err
		}
		return nil

	case capabilities.TransportDNSNull, capabilities.TransportICMP:
		// Covert channels for severe throttling
		return errors.New("covert channels initialized")

	default:
		return errors.New("unknown transport")
	}
}

func (fe *FallbackEngine) StopCurrentTransport() {
	fe.mu.Lock()
	defer fe.mu.Unlock()
	_ = fe.psiphonCtrl.StopTunnel()
}
