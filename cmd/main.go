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

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/thebeginner86/hippocampus/internal/config"
	"github.com/thebeginner86/hippocampus/internal/consensus/cluster"
	"github.com/thebeginner86/hippocampus/internal/handlers"
	"github.com/thebeginner86/hippocampus/resp"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting Hippocampus: %s\n", cfg)

	// Create necessary directories
	if err := os.MkdirAll(cfg.Storage.DataDir, 0755); err != nil {
		fmt.Printf("Failed to create data directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize handler with AOF storage
	handlerInst, err := handlers.NewHandler(cfg.Storage.AofFile)
	if err != nil {
		fmt.Printf("Failed to initialize handler: %v\n", err)
		os.Exit(1)
	}
	defer handlerInst.AofH.Close()

	// Replay AOF
	err = handlerInst.AofH.Read(func(value *resp.Value) {
		resp := handlerInst.ExecuteCmd(value, true)
		if resp != nil && resp.Type == "error" {
			return
		}
	})
	if err != nil {
		fmt.Printf("Failed to replay AOF: %v\n", err)
		os.Exit(1)
	}

	var raftCluster *cluster.Cluster

	// Initialize RAFT cluster if enabled
	if cfg.IsClusterMode() && cfg.RAFT.Enabled {
		fmt.Printf("Initializing RAFT cluster mode with %d peers\n", len(cfg.Cluster.Peers))
		raftCluster = initializeRaftCluster(cfg)
		if err := raftCluster.Start(); err != nil {
			fmt.Printf("Failed to start RAFT cluster: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("[RAFT] Cluster started for node %s\n", cfg.NodeID)
	} else {
		fmt.Println("Running in single-node mode (no RAFT clustering)")
	}

	// Setup signal handler for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start TCP server
	fmt.Printf("Listening on port :%s\n", cfg.Port)
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		fmt.Printf("Failed to listen: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	// Graceful shutdown handling
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		listener.Close()
		if raftCluster != nil {
			if err := raftCluster.Stop(); err != nil {
				fmt.Printf("Error stopping RAFT cluster: %v\n", err)
			}
		}
		handlerInst.AofH.Close()
		os.Exit(0)
	}()

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}

		go handleConnection(conn, handlerInst, raftCluster, cfg)
	}
}

// handleConnection handles a single client connection
func handleConnection(conn net.Conn, handlerInst *handlers.Handler, raftCluster *cluster.Cluster, cfg *config.Config) {
	defer conn.Close()

	for {
		respClient := resp.NewResp(conn)
		req, err := respClient.Read()
		if err != nil {
			return
		}

		if req.Type != "array" || len(req.Array) == 0 {
			continue
		}

		// Check if operation should be replicated (RAFT cluster)
		isWriteOp := isWriteCommand(req.Array[0].Bulk)
		if isWriteOp && raftCluster != nil && cfg.IsClusterMode() {
			handleClusteredWrite(conn, req, handlerInst, raftCluster)
		} else {
			handleLocalWrite(conn, req, handlerInst)
		}
	}
}

// handleLocalWrite handles write in single-node mode
func handleLocalWrite(conn net.Conn, req *resp.Value, handlerInst *handlers.Handler) {
	writer := resp.NewWriter(conn)
	response := handlerInst.ExecuteCmd(req, false)
	writer.Write(response)
}

// handleClusteredWrite handles write in RAFT cluster mode
func handleClusteredWrite(conn net.Conn, req *resp.Value, handlerInst *handlers.Handler, raftCluster *cluster.Cluster) {
	writer := resp.NewWriter(conn)

	// If not leader, return error to redirect to leader
	if !raftCluster.IsLeader() {
		leader := raftCluster.GetLeader()
		if leader == "" {
			response := &resp.Value{Type: "error", Bulk: "ERR no leader elected"}
			writer.Write(response)
		} else {
			response := &resp.Value{Type: "error", Bulk: fmt.Sprintf("ERR redirect to leader: %s", leader)}
			writer.Write(response)
		}
		return
	}

	// Execute locally and it will be replicated via RAFT
	response := handlerInst.ExecuteCmd(req, false)
	writer.Write(response)
}

// isWriteCommand checks if a command is a write operation
func isWriteCommand(cmd string) bool {
	writeCommands := map[string]bool{
		"SET":      true,
		"HSET":     true,
		"DEL":      true,
		"FLUSHDB":  true,
		"FLUSHALL": true,
	}
	return writeCommands[cmd]
}

// initializeRaftCluster initializes the RAFT cluster
func initializeRaftCluster(cfg *config.Config) *cluster.Cluster {
	// Convert PeerConfig to cluster.Peer
	peers := make(map[string]cluster.Peer)
	for _, p := range cfg.Cluster.Peers {
		peers[p.NodeID] = cluster.Peer{
			ID:      p.NodeID,
			Address: p.Host,
		}
	}

	// Create cluster instance
	raftCluster := cluster.NewCluster(cfg.NodeID, peers, cfg.RAFT.RpcPort)
	return raftCluster
}
