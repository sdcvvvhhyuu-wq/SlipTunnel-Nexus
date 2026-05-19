package capabilities

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

// =====================================================================
// DECENTRALIZED LOCAL P2P MESH SYNCHRONIZATION LAYER (NON-ROOT)
// =====================================================================
// Utilizes UDP Multicast to discover trusted peers locally and sync updated
// bridge metadata and routing tokens during total international connectivity blackouts.

const MulticastAddress = "224.0.0.251:9999"
const MaxDatagramSize = 8192

type MeshPayload struct {
	NodeID        string `json:"node_id"`
	AssetCategory string `json:"asset_category"`
	EncryptedData string `json:"encrypted_data"` // AES-GCM encrypted via EncryptedDB
	Timestamp     int64  `json:"timestamp"`
}

type MeshSynchronizer struct {
	mu         sync.RWMutex
	db         *EncryptedDB
	nodeID     string
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	conn       *net.UDPConn
}

func NewMeshSynchronizer(db *EncryptedDB, nodeID string) *MeshSynchronizer {
	ctx, cancel := context.WithCancel(context.Background())
	return &MeshSynchronizer{
		db:     db,
		nodeID: nodeID,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (m *MeshSynchronizer) Start() error {
	addr, err := net.ResolveUDPAddr("udp4", MulticastAddress)
	if err != nil {
		return err
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.conn = conn
	m.mu.Unlock()

	m.wg.Add(2)
	go m.listenLoop()
	go m.broadcastLoop(addr)

	return nil
}

func (m *MeshSynchronizer) Stop() {
	m.cancel()
	m.mu.Lock()
	if m.conn != nil {
		m.conn.Close()
	}
	m.mu.Unlock()
	m.wg.Wait()
}

func (m *MeshSynchronizer) listenLoop() {
	defer m.wg.Done()
	buffer := make([]byte, MaxDatagramSize)

	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			m.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			n, _, err := m.conn.ReadFromUDP(buffer)
			if err != nil {
				continue
			}

			var payload MeshPayload
			if err := json.Unmarshal(buffer[:n], &payload); err != nil {
				continue
			}

			// Ignore our own broadcasts
			if payload.NodeID == m.nodeID {
				continue
			}

			// Store the synced asset into the local EncryptedDB
			// The payload.EncryptedData is already AES-GCM encrypted by the peer.
			// In a real scenario, peers share a localized mesh key.
			err = m.db.StoreAsset(payload.NodeID+"_"+payload.AssetCategory, payload.AssetCategory, []byte(payload.EncryptedData), 0)
			if err != nil {
				log.Printf("MeshSync: Failed to store asset: %v", err)
			}
		}
	}
}

func (m *MeshSynchronizer) broadcastLoop(targetAddr *net.UDPAddr) {
	defer m.wg.Done()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			// Retrieve a healthy bridge asset to share with the local mesh
			assetData, err := m.db.RetrieveAsset("local_best_bridge")
			if err != nil {
				continue
			}

			payload := MeshPayload{
				NodeID:        m.nodeID,
				AssetCategory: "TOR_BRIDGE",
				EncryptedData: string(assetData), // Pre-encrypted by DB
				Timestamp:     time.Now().Unix(),
			}

			data, err := json.Marshal(payload)
			if err == nil {
				m.mu.RLock()
				if m.conn != nil {
					m.conn.WriteToUDP(data, targetAddr)
				}
				m.mu.RUnlock()
			}
		}
	}
}
