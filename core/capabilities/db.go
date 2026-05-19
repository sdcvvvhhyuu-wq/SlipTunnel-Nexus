package capabilities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type EncryptedDB struct {
	mu    sync.RWMutex
	db    *sql.DB
	gcm   cipher.AEAD
}

// NewEncryptedDB initializes a low-overhead SQLite database and sets up AES-GCM for payload encryption.
func NewEncryptedDB(dbPath string, masterKey string) (*EncryptedDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables if they do not exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS network_assets (
		id TEXT PRIMARY KEY,
		category TEXT NOT NULL,
		encrypted_payload TEXT NOT NULL,
		latency_profile INTEGER DEFAULT 0,
		last_verified DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	// Derive a 32-byte key using SHA-256 for AES-256
	hash := sha256.Sum256([]byte(masterKey))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &EncryptedDB{
		db:  db,
		gcm: gcm,
	}, nil
}

func (e *EncryptedDB) encrypt(data []byte) (string, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := e.gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

func (e *EncryptedDB) decrypt(hexStr string) ([]byte, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	nonceSize := e.gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return e.gcm.Open(nil, nonce, ciphertext, nil)
}

// StoreAsset securely encrypts and persists discovered IP pools, CDN domains, or Tor bridges.
func (e *EncryptedDB) StoreAsset(id string, category string, payload []byte, latency int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	encryptedPayload, err := e.encrypt(payload)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO network_assets (id, category, encrypted_payload, latency_profile, last_verified) 
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(id) DO UPDATE SET 
		encrypted_payload=excluded.encrypted_payload,
		latency_profile=excluded.latency_profile,
		last_verified=CURRENT_TIMESTAMP;`

	_, err = e.db.Exec(query, id, category, encryptedPayload, latency)
	return err
}

// RetrieveAsset decrypts and returns a cached network asset.
func (e *EncryptedDB) RetrieveAsset(id string) ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var encryptedPayload string
	query := `SELECT encrypted_payload FROM network_assets WHERE id = ?`
	err := e.db.QueryRow(query, id).Scan(&encryptedPayload)
	if err != nil {
		return nil, err
	}

	return e.decrypt(encryptedPayload)
}

// Close safely shuts down the database connection.
func (e *EncryptedDB) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.db.Close()
}
