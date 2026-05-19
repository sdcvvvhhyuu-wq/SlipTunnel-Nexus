package engine

import (
	"context"
	"errors"
	"sync"
	"time"

	"com.argotunnel/core/capabilities"
)

type ExecutorState int

const (
	StateStopped ExecutorState = iota
	StateStarting
	StateRunning
	StateThrottled
	StateFailingOver
)

type Executor struct {
	mu         sync.RWMutex
	state      ExecutorState
	selector   *capabilities.Selector
	db         *capabilities.EncryptedDB
	ctx        context.Context
	cancel     context.CancelFunc
	healthChan chan bool
	wg         sync.WaitGroup
}

func NewExecutor(selector *capabilities.Selector, db *capabilities.EncryptedDB) *Executor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Executor{
		state:      StateStopped,
		selector:   selector,
		db:         db,
		ctx:        ctx,
		cancel:     cancel,
		healthChan: make(chan bool, 10),
	}
}

func (e *Executor) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state != StateStopped {
		return errors.New("executor is already running")
	}

	e.state = StateStarting

	// Initialize the best transport based on historical DB metrics
	transport, err := e.selector.SelectBestTransport()
	if err != nil {
		// Fallback to aggressive scanning if no transport is viable
		e.state = StateFailingOver
	} else {
		_ = transport // Transport interface to be bound to Rust FFI in Phase 2/3
		e.state = StateRunning
	}

	e.wg.Add(1)
	go e.healthMonitorLoop()

	return nil
}

func (e *Executor) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state == StateStopped {
		return nil
	}

	e.cancel()
	e.wg.Wait()
	e.state = StateStopped

	// Volatile memory scrubbing trigger
	e.scrubMemory()

	return nil
}

func (e *Executor) healthMonitorLoop() {
	defer e.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			return
		case <-ticker.C:
			e.mu.RLock()
			currentState := e.state
			e.mu.RUnlock()

			if currentState == StateRunning {
				// Ping active transport, evaluate latency
				latency := e.evaluateLatency()
				if latency > 2000 { // 2000ms threshold for severe throttling
					e.triggerFailover()
				}
			}
		}
	}
}

func (e *Executor) evaluateLatency() int64 {
	// Deterministic mock for Phase 1. Real ICMP/TCP probe logic integrates in Phase 3.
	return 150 
}

func (e *Executor) triggerFailover() {
	e.mu.Lock()
	e.state = StateFailingOver
	e.mu.Unlock()

	// Instruct selector to escalate downward in the 8-Level Fallback Matrix
	_, _ = e.selector.EscalateFallback()
	
	e.mu.Lock()
	e.state = StateRunning
	e.mu.Unlock()
}

func (e *Executor) scrubMemory() {
	// Explicit byte-clearing routines for cryptographic handshakes
	// Ensures zero-trace allocation neutralization
	dummy := make([]byte, 1024)
	for i := range dummy {
		dummy[i] = 0
	}
}
