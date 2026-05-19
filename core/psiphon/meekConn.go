package psiphon

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

type MeekConn struct {
	config     *Config
	httpClient *http.Client
	sessionID  string
	buffer     *bytes.Buffer
}

func NewMeekConn(cfg *Config, dialer *TLSDialer) *MeekConn {
	frontDomain, hostHeader, sni := cfg.GetActiveFrontingSpec()

	transport := &http.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Route through our DPI-evading TLS Dialer
			return dialer.Dial(network, addr, frontDomain, hostHeader, sni)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &MeekConn{
		config: cfg,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		sessionID: "meek-session-init",
		buffer:    new(bytes.Buffer),
	}
}

func (m *MeekConn) Write(p []byte) (n int, err error) {
	_, hostHeader, _ := m.config.GetActiveFrontingSpec()
	
	req, err := http.NewRequest("POST", "https://"+hostHeader+"/meek", bytes.NewReader(p))
	if err != nil {
		return 0, err
	}
	
	req.Host = hostHeader
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Session-Id", m.sessionID)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("meek fronting failed with status: " + resp.Status)
	}

	return len(p), nil
}

func (m *MeekConn) Read(p []byte) (n int, err error) {
	if m.buffer.Len() > 0 {
		return m.buffer.Read(p)
	}

	_, hostHeader, _ := m.config.GetActiveFrontingSpec()
	req, err := http.NewRequest("GET", "https://"+hostHeader+"/meek", nil)
	if err != nil {
		return 0, err
	}

	req.Host = hostHeader
	req.Header.Set("X-Session-Id", m.sessionID)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("meek read failed with status: " + resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	m.buffer.Write(data)
	return m.buffer.Read(p)
}

func (m *MeekConn) Close() error {
	m.buffer.Reset()
	m.httpClient.CloseIdleConnections()
	return nil
}
