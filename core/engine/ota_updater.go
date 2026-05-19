package engine

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

// =====================================================================
// SECURE OVER-THE-AIR (OTA) UPDATE ENGINE
// =====================================================================

type OTAUpdater struct {
	UpdateURL string
	PublicKey ed25519.PublicKey
}

// NewOTAUpdater initializes the updater with a hardcoded Ed25519 public key for signature verification.
func NewOTAUpdater(url string, pubKeyHex string) (*OTAUpdater, error) {
	pubKey, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, err
	}
	return &OTAUpdater{
		UpdateURL: url,
		PublicKey: pubKey,
	}, nil
}

// CheckAndUpdate fetches the binary diff over the active tunnel, verifies the cryptographic signature,
// and performs an atomic replacement of the core routing assets.
func (ota *OTAUpdater) CheckAndUpdate() error {
	// Client routes through the active VPN interface automatically
	client := &http.Client{Timeout: 45 * time.Second}
	
	req, err := http.NewRequest("GET", ota.UpdateURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Nexus-OTA-Agent/1.4.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to fetch OTA payload, status: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Payload Structure: [64-byte Ed25519 Signature] + [Binary Data]
	if len(body) < 64 {
		return errors.New("invalid OTA payload size")
	}

	signature := body[:64]
	binaryData := body[64:]

	// Strict Signature Verification
	if !ed25519.Verify(ota.PublicKey, binaryData, signature) {
		return errors.New("CRITICAL: OTA signature verification failed. Payload rejected")
	}

	// Write to temporary file
	tmpFile := "/data/local/tmp/nexus_update.tmp"
	err = os.WriteFile(tmpFile, binaryData, 0755)
	if err != nil {
		return err
	}

	// Atomic rename to prevent corruption during crash/power-loss
	targetFile := "/data/local/tmp/nexus_core_assets.bin"
	err = os.Rename(tmpFile, targetFile)
	if err != nil {
		return err
	}

	return nil
}
