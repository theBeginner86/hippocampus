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
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ElectionManager handles leader election with optimizations
type ElectionManager struct {
	node        *RaftNode
	votesRcvd   atomic.Int32
	votesNeeded int32
	mu          sync.RWMutex
	// Track peer failures for exponential backoff
	peerFailures map[string]int
}

// NewElectionManager creates election manager
func NewElectionManager(node *RaftNode, peerCount int) *ElectionManager {
	return &ElectionManager{
		node:         node,
		votesNeeded:  int32((peerCount + 2) / 2),
		peerFailures: make(map[string]int),
	}
}

// StartElection starts optimized leader election
func (em *ElectionManager) StartElection(peers map[string]string) (bool, error) {
	// Increment term and vote for self
	_ = em.node.IncrementTerm()
	em.node.SetState(Candidate)
	em.votesRcvd.Store(1)

	_, _ = em.node.GetLastLogIndexAndTerm()

	// Request votes from all peers with optimized timeout
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	votes := make(chan bool, len(peers))

	for peerID := range peers {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			// Apply exponential backoff for failed peers
			em.mu.RLock()
			failures := em.peerFailures[id]
			em.mu.RUnlock()

			backoff := time.Duration(1<<uint(min(failures, 5))) * 10 * time.Millisecond
			select {
			case <-time.After(backoff):
				// Simulate vote request (RPC would go here)
				voteGranted := false // Will be updated with actual RPC

				if voteGranted {
					if em.votesRcvd.Add(1) >= em.votesNeeded {
						votes <- true
					}
				}
			case <-ctx.Done():
				return
			}
		}(peerID)
	}

	go func() {
		wg.Wait()
		close(votes)
	}()

	// Wait for election result
	for range votes {
		if em.votesRcvd.Load() >= em.votesNeeded {
			em.node.SetState(Leader)
			return true, nil
		}
	}

	if em.votesRcvd.Load() >= em.votesNeeded {
		em.node.SetState(Leader)
		return true, nil
	}

	em.node.SetState(Follower)
	return false, fmt.Errorf("election failed: votes %d < needed %d", em.votesRcvd.Load(), em.votesNeeded)
}

// RecordFailure records RPC failure for exponential backoff
func (em *ElectionManager) RecordFailure(peerID string) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.peerFailures[peerID]++
}

// RecordSuccess resets failure count
func (em *ElectionManager) RecordSuccess(peerID string) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.peerFailures[peerID] = 0
}

// GetStats returns election statistics
func (em *ElectionManager) GetStats() map[string]int {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.peerFailures
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
