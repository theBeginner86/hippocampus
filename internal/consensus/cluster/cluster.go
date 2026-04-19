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

package cluster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/thebeginner86/hippocampus/internal/consensus/raft"
)

// Peer represents a peer in the cluster
type Peer struct {
	ID      string
	Address string
}

// Cluster manages RAFT cluster
type Cluster struct {
	node          *raft.RaftNode
	peers         map[string]Peer
	clients       map[string]*raft.RPCClient
	server        *raft.RPCServer
	election      *raft.ElectionManager
	stopCh        chan struct{}
	stoppedCh     chan struct{}
	mu            sync.RWMutex
	lastHeartbeat time.Time
	stats         ClusterStats
}

// ClusterStats holds cluster statistics
type ClusterStats struct {
	mu                 sync.RWMutex
	TotalCommitted     int64
	FailedReplicates   int64
	ElectionsHeld      int64
	CurrentLeader      string
	LastLeaderElection time.Time
}

// NewCluster creates new cluster
func NewCluster(nodeID string, peers map[string]Peer, rpcPort string) *Cluster {
	node := raft.NewRaftNode(nodeID)

	cluster := &Cluster{
		node:          node,
		peers:         peers,
		clients:       make(map[string]*raft.RPCClient),
		stopCh:        make(chan struct{}),
		stoppedCh:     make(chan struct{}),
		lastHeartbeat: time.Now(),
	}

	// Create RPC clients
	for peerID, peer := range peers {
		addr := fmt.Sprintf("%s:7000", peer.Address)
		cluster.clients[peerID] = raft.NewRPCClient(addr)
	}

	// Create election manager
	cluster.election = raft.NewElectionManager(node, len(peers)+1)

	// Create RPC server
	cluster.server = raft.NewRPCServer(node, rpcPort)

	return cluster
}

// Start starts the cluster
func (c *Cluster) Start() error {
	if err := c.server.Start(); err != nil {
		return err
	}

	go c.run()
	return nil
}

// run implements main RAFT loop
func (c *Cluster) run() {
	defer close(c.stoppedCh)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.tick()
		}
	}
}

// tick processes one tick
func (c *Cluster) tick() {
	c.mu.RLock()
	state := c.node.GetState()
	c.mu.RUnlock()

	switch state {
	case raft.Follower:
		c.tickFollower()
	case raft.Candidate:
		c.tickCandidate()
	case raft.Leader:
		c.tickLeader()
	}
}

// tickFollower handles follower logic
func (c *Cluster) tickFollower() {
	if c.node.CheckElectionTimeout() {
		c.startElection()
	}
}

// tickCandidate handles candidate logic
func (c *Cluster) tickCandidate() {
	if c.node.CheckElectionTimeout() {
		c.startElection()
	}
}

// tickLeader handles leader logic
func (c *Cluster) tickLeader() {
	c.sendHeartbeats()
}

// startElection starts election
func (c *Cluster) startElection() {
	// Convert peers to string map for election manager
	peerMap := make(map[string]string)
	for id, peer := range c.peers {
		peerMap[id] = peer.Address
	}

	success, err := c.election.StartElection(peerMap)

	if success {
		c.stats.mu.Lock()
		c.stats.ElectionsHeld++
		c.stats.CurrentLeader = c.node.NodeID
		c.stats.LastLeaderElection = time.Now()
		c.stats.mu.Unlock()

		fmt.Printf("[%s] Became leader for term %d\n", c.node.NodeID, c.node.GetCurrentTerm())
	} else if err != nil {
		fmt.Printf("[%s] Election failed: %v\n", c.node.NodeID, err)
	}
}

// sendHeartbeats sends heartbeats to all peers
func (c *Cluster) sendHeartbeats() {
	if !c.node.IsLeader() {
		return
	}

	for peerID := range c.peers {
		go c.sendHeartbeatToPeer(peerID)
	}
}

// sendHeartbeatToPeer sends heartbeat to peer
func (c *Cluster) sendHeartbeatToPeer(peerID string) {
	client, ok := c.clients[peerID]
	if !ok {
		return
	}

	args := &raft.AppendEntriesArgs{
		Term:         c.node.GetCurrentTerm(),
		LeaderID:     c.node.NodeID,
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		Entries:      []raft.LogEntry{},
		LeaderCommit: c.node.CommitIndex,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	reply, err := client.AppendEntries(ctx, args)
	if err != nil {
		c.election.RecordFailure(peerID)
		return
	}

	if reply.Term > c.node.GetCurrentTerm() {
		c.node.SetTerm(reply.Term)
		c.node.SetState(raft.Follower)
	} else {
		c.election.RecordSuccess(peerID)
	}
}

// Stop stops the cluster
func (c *Cluster) Stop() error {
	close(c.stopCh)
	<-c.stoppedCh

	if c.server != nil {
		if err := c.server.Stop(); err != nil {
			return err
		}
	}

	for _, client := range c.clients {
		client.Close()
	}

	return nil
}

// GetStatus returns cluster status
func (c *Cluster) GetStatus() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.stats.mu.RLock()
	leader := c.stats.CurrentLeader
	elections := c.stats.ElectionsHeld
	c.stats.mu.RUnlock()

	return map[string]interface{}{
		"node_id":        c.node.NodeID,
		"state":          c.node.GetState(),
		"term":           c.node.GetCurrentTerm(),
		"current_leader": leader,
		"elections_held": elections,
		"log_size":       len(c.node.Log),
		"peer_failures":  c.election.GetStats(),
	}
}

// AddLogEntry adds entry to log (leader only)
func (c *Cluster) AddLogEntry(cmd []byte, cmdType string) error {
	if !c.node.IsLeader() {
		return fmt.Errorf("not a leader")
	}

	entry := raft.LogEntry{
		CmdType: cmdType,
		Command: cmd,
	}
	c.node.AppendLogEntry(entry)
	return nil
}

// GetLeader returns current leader
func (c *Cluster) GetLeader() string {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()
	return c.stats.CurrentLeader
}

// IsLeader checks if this node is leader
func (c *Cluster) IsLeader() bool {
	return c.node.IsLeader()
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
