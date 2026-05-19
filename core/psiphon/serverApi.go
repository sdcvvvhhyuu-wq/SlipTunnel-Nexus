package psiphon

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"com.argotunnel/core/capabilities"
)

type ServerEntry struct {
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	PublicKey string `json:"public_key"`
}

type ServerAPI struct {
	db *capabilities.EncryptedDB
}

func NewServerAPI(db *capabilities.EncryptedDB) *ServerAPI {
	return &ServerAPI{
		db: db,
	}
}

// FetchRemoteBootstrap autonomously retrieves fresh server lists via Domain Fronting
func (api *ServerAPI) FetchRemoteBootstrap(cfg *Config) ([]ServerEntry, error) {
	frontDomain, hostHeader, _ := cfg.GetActiveFrontingSpec()

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", "https://"+frontDomain+"/server-list", nil)
	if err != nil {
		return nil, err
	}
	req.Host = hostHeader

	resp, err := client.Do(req)
	if err != nil {
		// Fallback to local encrypted cache if remote fetch fails (Intranet Lockdown)
		return api.loadLocalBootstrap()
	}
	defer resp.Body.Close()

	var entries []ServerEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return api.loadLocalBootstrap()
	}

	// Persist fresh entries to EncryptedDB
	data, _ := json.Marshal(entries)
	_ = api.db.StoreAsset("psiphon_server_list", "BOOTSTRAP", data, 0)

	return entries, nil
}

func (api *ServerAPI) loadLocalBootstrap() ([]ServerEntry, error) {
	data, err := api.db.RetrieveAsset("psiphon_server_list")
	if err != nil {
		return nil, errors.New("total blackout: remote fetch failed and local cache empty")
	}

	var entries []ServerEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
