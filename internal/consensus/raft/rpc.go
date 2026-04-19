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
	"net"
	"net/rpc"
	"sync"
	"time"
)

// RequestVoteArgs for RequestVote RPC
type RequestVoteArgs struct {
	Term         int64
	CandidateID  string
	LastLogIndex int64
	LastLogTerm  int64
}

// RequestVoteReply for RequestVote RPC
type RequestVoteReply struct {
	Term        int64
	VoteGranted bool
}

// AppendEntriesArgs for AppendEntries RPC
type AppendEntriesArgs struct {
	Term         int64
	LeaderID     string
	PrevLogIndex int64
	PrevLogTerm  int64
	Entries      []LogEntry
	LeaderCommit int64
}

// AppendEntriesReply for AppendEntries RPC
type AppendEntriesReply struct {
	Term    int64
	Success bool
}

// RPCHandler handles incoming RPC calls
type RPCHandler struct {
	node *RaftNode
	mu   sync.RWMutex
}

// NewRPCHandler creates RPC handler
func NewRPCHandler(node *RaftNode) *RPCHandler {
	return &RPCHandler{node: node}
}

// RequestVote handles RequestVote RPC
func (h *RPCHandler) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	currentTerm := h.node.GetCurrentTerm()

	if args.Term < currentTerm {
		reply.Term = currentTerm
		reply.VoteGranted = false
		return nil
	}

	if args.Term > currentTerm {
		h.node.SetTerm(args.Term)
		currentTerm = args.Term
	}

	reply.Term = currentTerm

	// Check if already voted
	if h.node.VotedFor != "" && h.node.VotedFor != args.CandidateID {
		reply.VoteGranted = false
		return nil
	}

	// Check candidate log is up-to-date
	lastLogIdx, lastLogTerm := h.node.GetLastLogIndexAndTerm()
	if args.LastLogTerm < lastLogTerm || (args.LastLogTerm == lastLogTerm && args.LastLogIndex < lastLogIdx) {
		reply.VoteGranted = false
		return nil
	}

	// Grant vote
	h.node.VotedFor = args.CandidateID
	h.node.ResetElectionTimeout()
	reply.VoteGranted = true

	return nil
}

// AppendEntries handles AppendEntries RPC
func (h *RPCHandler) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	currentTerm := h.node.GetCurrentTerm()

	if args.Term < currentTerm {
		reply.Term = currentTerm
		reply.Success = false
		return nil
	}

	if args.Term > currentTerm {
		h.node.SetTerm(args.Term)
	}

	reply.Term = h.node.CurrentTerm
	h.node.ResetElectionTimeout()
	h.node.SetState(Follower)

	// Check previous log entry
	if args.PrevLogIndex > 0 {
		if args.PrevLogIndex > int64(len(h.node.Log)-1) {
			reply.Success = false
			return nil
		}
		if h.node.Log[args.PrevLogIndex].Term != args.PrevLogTerm {
			reply.Success = false
			return nil
		}
	}

	// Append entries
	for _, entry := range args.Entries {
		h.node.AppendLogEntry(entry)
	}

	// Update commit index
	if args.LeaderCommit > h.node.CommitIndex {
		newCommit := args.LeaderCommit
		if args.LeaderCommit > int64(len(h.node.Log)-1) {
			newCommit = int64(len(h.node.Log) - 1)
		}
		h.node.CommitIndex = newCommit
	}

	reply.Success = true
	return nil
}

// RPCClient for outgoing RPC calls
type RPCClient struct {
	address string
	conn    *rpc.Client
	mu      sync.Mutex
}

// NewRPCClient creates RPC client
func NewRPCClient(address string) *RPCClient {
	return &RPCClient{address: address}
}

// Connect establishes connection
func (c *RPCClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	conn, err := net.DialTimeout("tcp", c.address, 5*time.Second)
	if err != nil {
		return err
	}

	c.conn = rpc.NewClient(conn)
	return nil
}

// Close closes connection
func (c *RPCClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// RequestVote sends RequestVote RPC
func (c *RPCClient) RequestVote(ctx context.Context, args *RequestVoteArgs) (*RequestVoteReply, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	reply := &RequestVoteReply{}
	done := make(chan *rpc.Call, 1)

	go func() {
		c.mu.Lock()
		if c.conn != nil {
			c.conn.Go("RPCHandler.RequestVote", args, reply, done)
		}
		c.mu.Unlock()
	}()

	select {
	case <-ctx.Done():
		c.Close()
		return nil, ctx.Err()
	case call := <-done:
		if call.Error != nil {
			c.Close()
			return nil, call.Error
		}
		return reply, nil
	case <-time.After(1 * time.Second):
		c.Close()
		return nil, fmt.Errorf("RPC timeout")
	}
}

// AppendEntries sends AppendEntries RPC
func (c *RPCClient) AppendEntries(ctx context.Context, args *AppendEntriesArgs) (*AppendEntriesReply, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	reply := &AppendEntriesReply{}
	done := make(chan *rpc.Call, 1)

	go func() {
		c.mu.Lock()
		if c.conn != nil {
			c.conn.Go("RPCHandler.AppendEntries", args, reply, done)
		}
		c.mu.Unlock()
	}()

	select {
	case <-ctx.Done():
		c.Close()
		return nil, ctx.Err()
	case call := <-done:
		if call.Error != nil {
			c.Close()
			return nil, call.Error
		}
		return reply, nil
	case <-time.After(1 * time.Second):
		c.Close()
		return nil, fmt.Errorf("RPC timeout")
	}
}

// RPCServer manages RPC server
type RPCServer struct {
	node     *RaftNode
	port     string
	listener net.Listener
	mu       sync.Mutex
	running  bool
}

// NewRPCServer creates RPC server
func NewRPCServer(node *RaftNode, port string) *RPCServer {
	return &RPCServer{
		node: node,
		port: port,
	}
}

// Start starts RPC server
func (s *RPCServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("server already running")
	}

	handler := NewRPCHandler(s.node)
	rpc.Register(handler)

	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}

	s.listener = listener
	s.running = true

	go func() {
		for s.running {
			conn, err := listener.Accept()
			if err != nil {
				if s.running {
					fmt.Printf("accept error: %v\n", err)
				}
				break
			}
			go rpc.ServeConn(conn)
		}
	}()

	return nil
}

// Stop stops RPC server
func (s *RPCServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.running = false
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
