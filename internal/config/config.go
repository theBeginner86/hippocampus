//  Copyright 2024 Pranav Singh

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ServerMode defines the server operating mode
type ServerMode string

const (
	// SingleMode - Single node, no clustering
	SingleMode ServerMode = "single"
	// ClusterMode - Multi-node RAFT cluster
	ClusterMode ServerMode = "cluster"
)

// RaftConfig holds RAFT-specific configuration
type RaftConfig struct {
	Enabled            bool
	RpcPort            string
	MinElectionTimeout int // milliseconds
	MaxElectionTimeout int
	HeartbeatInterval  int
	SnapshotInterval   int64
	MaxLogSize         int
}

// ClusterConfig holds cluster-specific configuration
type ClusterConfig struct {
	Enabled     bool
	NodeName    string // k8s service for peer discovery
	ServiceName string
	Peers       []PeerConfig
}

// PeerConfig represents a peer in the cluster
type PeerConfig struct {
	NodeID   string
	Host     string
	RpcPort  string
	DataPort string
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	DataDir        string
	AofFile        string
	SnapshotDir    string
	EnableAof      bool
	AofFsyncPolicy string // always, everysec, no
}

// Config holds all server configuration
type Config struct {
	// Server mode
	Mode ServerMode
	// Server basics
	Port        string
	NodeID      string
	DatabaseDir string
	// RAFT configuration
	RAFT RaftConfig
	// Cluster configuration
	Cluster ClusterConfig
	// Storage configuration
	Storage StorageConfig
}

// Load loads configuration from environment variables and flags
func Load() (*Config, error) {
	cfg := &Config{
		Port:        getEnv("HIPPOCAMPUS_PORT", "6379"),
		NodeID:      getEnv("HIPPOCAMPUS_NODE_ID", "node-1"),
		DatabaseDir: getEnv("HIPPOCAMPUS_DATA_DIR", "/data"),
		Mode:        ServerMode(getEnv("HIPPOCAMPUS_MODE", "single")),
	}

	// Parse flags
	flag.StringVar(&cfg.Port, "port", cfg.Port, "Server port")
	flag.StringVar(&cfg.NodeID, "node-id", cfg.NodeID, "Node ID")
	flag.StringVar(&cfg.DatabaseDir, "data-dir", cfg.DatabaseDir, "Data directory")
	flag.Var((*modeFlag)(cfg), "mode", "Server mode: single|cluster")
	flag.Parse()

	// Configure RAFT
	cfg.RAFT = RaftConfig{
		Enabled:            getEnvBool("RAFT_ENABLED", true),
		RpcPort:            getEnv("RAFT_RPC_PORT", "7000"),
		MinElectionTimeout: getEnvInt("RAFT_MIN_ELECTION_TIMEOUT", 150),
		MaxElectionTimeout: getEnvInt("RAFT_MAX_ELECTION_TIMEOUT", 300),
		HeartbeatInterval:  getEnvInt("RAFT_HEARTBEAT_INTERVAL", 50),
		SnapshotInterval:   getEnvInt64("RAFT_SNAPSHOT_INTERVAL", 100000),
		MaxLogSize:         getEnvInt("RAFT_MAX_LOG_SIZE", 1000000),
	}

	// Configure storage
	cfg.Storage = StorageConfig{
		DataDir:        cfg.DatabaseDir,
		AofFile:        getEnv("STORAGE_AOF_FILE", "/data/database.aof"),
		SnapshotDir:    getEnv("STORAGE_SNAPSHOT_DIR", "/data/snapshots"),
		EnableAof:      getEnvBool("STORAGE_ENABLE_AOF", true),
		AofFsyncPolicy: getEnv("STORAGE_AOF_FSYNC", "everysec"),
	}

	// Configure cluster if mode is cluster
	if cfg.Mode == ClusterMode {
		cfg.Cluster = ClusterConfig{
			Enabled:     true,
			NodeName:    getEnv("CLUSTER_NODE_NAME", cfg.NodeID),
			ServiceName: getEnv("CLUSTER_SERVICE_NAME", "hippocampus-cluster"),
		}

		// Parse peers from environment
		peersStr := getEnv("CLUSTER_PEERS", "")
		if peersStr != "" {
			cfg.Cluster.Peers = parsePeers(peersStr)
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.NodeID == "" {
		return fmt.Errorf("node-id is required")
	}

	if c.Mode != SingleMode && c.Mode != ClusterMode {
		return fmt.Errorf("invalid mode: %s (must be 'single' or 'cluster')", c.Mode)
	}

	if c.Mode == ClusterMode && len(c.Cluster.Peers) == 0 {
		return fmt.Errorf("cluster mode requires at least one peer")
	}

	return nil
}

// IsClusterMode returns true if running in cluster mode
func (c *Config) IsClusterMode() bool {
	return c.Mode == ClusterMode
}

// IsSingleMode returns true if running in single mode
func (c *Config) IsSingleMode() bool {
	return c.Mode == SingleMode
}

// String returns a string representation of config
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{Mode: %s, NodeID: %s, Port: %s, RaftEnabled: %v, Peers: %d}",
		c.Mode, c.NodeID, c.Port, c.RAFT.Enabled, len(c.Cluster.Peers),
	)
}

// Helper functions

func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	val := getEnv(key, "")
	if val == "" {
		return defaultVal
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}

func getEnvInt(key string, defaultVal int) int {
	val := getEnv(key, "")
	if val == "" {
		return defaultVal
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func getEnvInt64(key string, defaultVal int64) int64 {
	val := getEnv(key, "")
	if val == "" {
		return defaultVal
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return i
}

// parsePeers parses peer configuration from a string
// Format: "node1:host1:7000:6379,node2:host2:7000:6379"
func parsePeers(peersStr string) []PeerConfig {
	entries := strings.Split(peersStr, ",")
	var peers []PeerConfig

	for _, entry := range entries {
		parts := strings.Split(strings.TrimSpace(entry), ":")
		if len(parts) >= 4 {
			peers = append(peers, PeerConfig{
				NodeID:   parts[0],
				Host:     parts[1],
				RpcPort:  parts[2],
				DataPort: parts[3],
			})
		}
	}
	return peers
}

// modeFlag implements flag.Value for ServerMode
type modeFlag Config

func (m *modeFlag) String() string {
	return string(m.Mode)
}

func (m *modeFlag) Set(value string) error {
	mode := ServerMode(strings.ToLower(value))
	if mode != SingleMode && mode != ClusterMode {
		return fmt.Errorf("invalid mode: %s", value)
	}
	m.Mode = mode
	return nil
}
