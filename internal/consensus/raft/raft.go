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

package raft

import (
	"math/rand"
	"sync"
	"time"
)

// RaftState represents the current state of a Raft node
type RaftState string

const (
	Follower  RaftState = "follower"
	Candidate RaftState = "candidate"
	Leader    RaftState = "leader"
)

// LogEntry represents a single entry in the Raft log
type LogEntry struct {
	Term    int64
	Index   int64
	Command []byte
	CmdType string
}

// RaftNode represents a node in the Raft cluster
type RaftNode struct {
	NodeID          string
	CurrentTerm     int64
	VotedFor        string
	State           RaftState
	Log             []LogEntry
	CommitIndex     int64
	LastApplied     int64
	NextIndex       map[string]int64
	MatchIndex      map[string]int64
	LastHeartbeat   time.Time
	ElectionTimeout time.Duration
	mu              sync.RWMutex
}

// NewRaftNode creates a new Raft node
func NewRaftNode(nodeID string) *RaftNode {
	return &RaftNode{
		NodeID:          nodeID,
		CurrentTerm:     0,
		VotedFor:        "",
		State:           Follower,
		Log:             make([]LogEntry, 0),
		CommitIndex:     0,
		LastApplied:     0,
		NextIndex:       make(map[string]int64),
		MatchIndex:      make(map[string]int64),
		LastHeartbeat:   time.Now(),
		ElectionTimeout: generateTimeout(),
	}
}

// GetState returns the current state
func (n *RaftNode) GetState() RaftState {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.State
}

// SetState sets the state
func (n *RaftNode) SetState(state RaftState) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.State = state
}

// GetCurrentTerm returns the current term
func (n *RaftNode) GetCurrentTerm() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.CurrentTerm
}

// IncrementTerm increments the term
func (n *RaftNode) IncrementTerm() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.CurrentTerm++
	n.VotedFor = ""
	return n.CurrentTerm
}

// SetTerm sets the term if higher
func (n *RaftNode) SetTerm(term int64) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	if term > n.CurrentTerm {
		n.CurrentTerm = term
		n.VotedFor = ""
		n.State = Follower
		return true
	}
	return false
}

// AppendLogEntry adds a log entry
func (n *RaftNode) AppendLogEntry(entry LogEntry) {
	n.mu.Lock()
	defer n.mu.Unlock()
	entry.Index = int64(len(n.Log))
	entry.Term = n.CurrentTerm
	n.Log = append(n.Log, entry)
}

// GetLastLogIndexAndTerm returns last log index and term
func (n *RaftNode) GetLastLogIndexAndTerm() (int64, int64) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if len(n.Log) == 0 {
		return 0, 0
	}
	last := n.Log[len(n.Log)-1]
	return last.Index, last.Term
}

// IsLeader checks if node is leader
func (n *RaftNode) IsLeader() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.State == Leader
}

// ResetElectionTimeout resets election timeout
func (n *RaftNode) ResetElectionTimeout() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.LastHeartbeat = time.Now()
	n.ElectionTimeout = generateTimeout()
}

// CheckElectionTimeout checks if timeout expired
func (n *RaftNode) CheckElectionTimeout() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return time.Since(n.LastHeartbeat) > n.ElectionTimeout
}

func generateTimeout() time.Duration {
	r := rand.Int63n(150) + 150
	return time.Duration(r) * time.Millisecond
}
