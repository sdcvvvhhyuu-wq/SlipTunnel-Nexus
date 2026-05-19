package psiphon

import (
	"context"
	"errors"
	"sync"
	"time"

	"com.argotunnel/core/capabilities"
)

type ControllerState int

const (
	StateDisconnected ControllerState = iota
	StateConnecting
	StateConnected
)

type Controller struct {
	mu        sync.RWMutex
	state     ControllerState
	config    *Config
	api       *ServerAPI
	dialer    *TLSDialer
	meekConn  *MeekConn
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewController(cfg *Config, db *capabilities.EncryptedDB) *Controller {
	ctx, cancel := context.WithCancel(context.Background())
	return &Controller{
		state:  StateDisconnected,
		config: cfg,
		api:    NewServerAPI(db),
		dialer: NewTLSDialer(10 * time.Second),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *Controller) StartTunnel() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state != StateDisconnected {
		return errors.New("tunnel already active or connecting")
	}
	c.state = StateConnecting

	// 1. Autonomous Bootstrapping
	servers, err := c.api.FetchRemoteBootstrap(c.config)
	if err != nil || len(servers) == 0 {
		c.state = StateDisconnected
		return errors.New("failed to bootstrap psiphon network")
	}

	// 2. Establish Meek Connection (Domain Fronting)
	c.meekConn = NewMeekConn(c.config, c.dialer)
	
	// 3. Send initial handshake payload
	_, err = c.meekConn.Write([]byte("PSIPHON-HANDSHAKE-INIT"))
	if err != nil {
		c.state = StateDisconnected
		return err
	}

	c.state = StateConnected
	return nil
}

func (c *Controller) StopTunnel() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == StateDisconnected {
		return nil
	}

	c.cancel()
	if c.meekConn != nil {
		c.meekConn.Close()
	}
	
	c.state = StateDisconnected
	return nil
}
