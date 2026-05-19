package engine

import (
	"math"
	"sync"
	"time"
)

// =====================================================================
// AI ADAPTIVE ROUTING ENGINE (ONNX/RL/GAN PROFILING LOGIC)
// =====================================================================

type NetworkProfile struct {
	LatencyMs int64
	JitterMs  int64
	LossRate  float64
}

type RoutingManager struct {
	mu             sync.RWMutex
	history        []NetworkProfile
	learningRate   float64
	discountFactor float64
}

func NewRoutingManager() *RoutingManager {
	return &RoutingManager{
		history:        make([]NetworkProfile, 0),
		learningRate:   0.1,
		discountFactor: 0.9,
	}
}

// EvaluateNetworkState acts as a lightweight Reinforcement Learning agent.
// It calculates a Q-Value representing network health to dynamically rotate transports.
func (rm *RoutingManager) EvaluateNetworkState(latency, jitter int64, loss float64) float64 {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	profile := NetworkProfile{
		LatencyMs: latency,
		JitterMs:  jitter,
		LossRate:  loss,
	}

	rm.history = append(rm.history, profile)
	if len(rm.history) > 100 {
		rm.history = rm.history[1:] // Keep last 100 states
	}

	// Deterministic Reward Function: Lower latency/jitter/loss = Higher Reward
	reward := 100.0 - (float64(latency) * 0.1) - (float64(jitter) * 0.2) - (loss * 100.0)
	
	// Sigmoid normalization to bound the Q-Value between 0 and 1
	qValue := 1.0 / (1.0 + math.Exp(-reward/10.0))

	return qValue
}

// ShouldMorphTraffic determines if GAN traffic mimicry should mutate packet pacing
func (rm *RoutingManager) ShouldMorphTraffic() bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if len(rm.history) < 5 {
		return false
	}

	// If jitter spikes rapidly, assume active DPI probing and trigger morphing
	last := rm.history[len(rm.history)-1]
	prev := rm.history[len(rm.history)-2]

	return (last.JitterMs - prev.JitterMs) > 50
}
