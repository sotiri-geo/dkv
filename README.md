# Distributed Key-Value Store Learning Project

A Go implementation of a distributed key-value store with consensus, built step-by-step to learn distributed systems concepts.

## Progress Overview

### ✅ Exercise 1.1: Single Node KV Store (COMPLETED)

**What we built:**
- Thread-safe in-memory key-value store
- HTTP API with GET, PUT, DELETE operations
- Comprehensive tests using `httptest`

**Key files:**
- `store/store.go` - KVStore with mutex-protected map
- `server/server.go` - HTTP handlers for CRUD operations
- Tests for both store and server layers

**What we learned:**
- Go concurrency primitives (`sync.RWMutex`)
- HTTP server implementation with `ServeMux`
- Testing HTTP handlers through the router
- Thread-safety and race conditions

---

### ✅ Exercise 1.2: Command Pattern & State Machine (COMPLETED)

**What we built:**

#### Architecture Changes:
- **Command Pattern**: Separated operations into command objects
- **CQRS**: Split commands (writes) from queries (reads)
  - `Command` interface - stateful operations (Put, Delete)
  - `Query` interface - stateless operations (Get)
- **Command Log**: Store maintains history of all state-changing operations
- **Replay Functionality**: Rebuild state by replaying command log

#### Implementation Details:

**Files:**
```
store/
├── store.go         # KVStore with Apply(), Execute(), Replay(), GetLog()
├── commands.go      # Command/Query interfaces and implementations
├── store_test.go    # Integration tests (logging, replay, determinism)
└── commands_test.go # Unit tests (command logic)

server/
├── server.go        # Updated HTTP handlers using commands
└── server_test.go   # Server tests
```

**Core Concepts Implemented:**
- ✅ Commands are immutable (private fields with constructor functions)
- ✅ Commands directly modify store data (no proxying through methods)
- ✅ Idempotent operations (safe to apply multiple times)
- ✅ Command logging (only for state-changing operations)
- ✅ Query execution (reads don't modify log)
- ✅ Replay mechanism (rebuild state from log)
- ✅ Thread-safe with proper locking
- ✅ Deterministic state machine (same commands = same state)

**Store API:**
```go
// For writes (logged)
func (s *KVStore) Apply(cmd Command)

// For reads (not logged)
func (s *KVStore) Execute(query Query) (string, error)

// Get command history
func (s *KVStore) GetLog() []Command

// Rebuild state from log
func (s *KVStore) Replay()
```

**Commands & Queries:**
- `PutCommand` - Add/update key-value pairs
- `DeleteCommand` - Remove keys
- `GetQuery` - Retrieve values by key

**HTTP API (Updated):**
```bash
# PUT
curl -X PUT http://localhost:8080/kv \
  -H "Content-Type: application/json" \
  -d '{"key":"name","value":"alice"}'

# GET
curl http://localhost:8080/kv/name

# DELETE
curl -X DELETE http://localhost:8080/kv/name
```

**What we learned:**
- Command Pattern for encapsulating operations
- CQRS (Command Query Responsibility Segregation)
- State machine design and determinism
- Idempotency in distributed systems
- Why commands shouldn't return errors
- Event sourcing concepts (log as source of truth)
- Package-level visibility in Go
- Avoiding deadlocks with proper lock management
- Testing strategies (unit vs integration)

---

## 🔜 Next Steps: Module 2 - Multi-Node Setup

Now that we have a solid single-node foundation with command logging, we're ready to move into distributed systems territory.

### Exercise 2.1: Node-to-Node Communication

**Goal:** Get multiple KVStore instances talking to each other.

**What we'll build:**
- Launch 3+ node instances
- Implement node discovery and peer management
- Create RPC mechanism for inter-node communication
- Implement health checks (Ping/Pong)
- **Transition from HTTP to gRPC** for structured communication

**Key concepts:**
- Service discovery
- RPC (Remote Procedure Call)
- Network programming
- gRPC and Protocol Buffers

**Why gRPC?**
- Type-safe communication (Protocol Buffers)
- Bi-directional streaming
- Better performance than HTTP/JSON
- Industry standard for microservices
- We'll need it for Raft anyway

### Exercise 2.2: Naive Replication (Without Consensus)

**Goal:** Broadcast writes to all nodes (simple replication, no conflict resolution).

**What we'll build:**
- Leader receives client writes
- Leader broadcasts commands to all peers
- All nodes apply the same commands
- Observe the problems this creates

**Problems we'll encounter (intentionally):**
- What if one node is down? Should writes fail?
- What if two clients write to different nodes simultaneously?
- Race conditions and inconsistent state
- Split-brain scenarios

**Why encounter these problems?**
This motivates why we need consensus algorithms like Raft!

---

## 🚀 Future Modules

### Module 3: Understanding Raft (Theory)
- Study the Raft paper
- Understand leader election
- Understand log replication
- Design Raft implementation

### Module 4: Implementing Raft
- Leader election
- Log replication with AppendEntries RPC
- Handling network partitions
- Log consistency

### Module 5: Advanced Features
- Snapshots and log compaction
- Membership changes
- Client interaction improvements
- Performance optimization

---

## 📚 Key Learnings So Far

### Technical Skills:
- ✅ Building thread-safe data structures in Go
- ✅ HTTP server patterns and REST APIs
- ✅ gRPC fundamentals (coming in Module 2)
- ✅ Testing strategies (unit, integration, concurrent)
- ✅ Command Pattern and CQRS
- ✅ State machine design

### Distributed Systems Concepts:
- ✅ State machine replication (same commands → same state)
- ✅ Idempotency (operations safe to retry)
- ✅ Determinism (no randomness or time-based state)
- ✅ Event sourcing (log as source of truth)
- 🔜 Consensus and leader election
- 🔜 Network partitions and fault tolerance
- 🔜 CAP theorem in practice

---

## 🏃 Quick Start

**Install dependencies:**
```bash
go mod tidy
```

**Run tests:**
```bash
# All tests
go test ./...

# With race detector
go test -race ./...

# Specific package
go test ./store/...
go test ./server/...
```

**Run server:**
```bash
go run main.go
```

**Test with curl:**
```bash
# PUT a key
curl -X PUT http://localhost:8080/kv \
  -H "Content-Type: application/json" \
  -d '{"key":"name","value":"alice"}'

# GET a key
curl http://localhost:8080/kv/name

# DELETE a key
curl -X DELETE http://localhost:8080/kv/name
```

---

## 📁 Project Structure

```
kvstore/
├── store/
│   ├── store.go         # KVStore with Apply, Execute, Replay
│   ├── commands.go      # Command/Query implementations
│   ├── store_test.go    # Integration tests
│   └── commands_test.go # Unit tests
├── server/
│   ├── server.go        # HTTP handlers
│   └── server_test.go   # Server tests
├── main.go              # Entry point
├── go.mod
└── README.md
```

---

## 🎓 Learning Resources

**Command Pattern:**
- [Refactoring Guru - Command Pattern](https://refactoring.guru/design-patterns/command)
- [Go Patterns - Command](https://github.com/tmrts/go-patterns/blob/master/behavioral/command.md)

**State Machine Replication:**
- [Raft Paper (Section 5.3)](https://raft.github.io/raft.pdf)
- [Martin Fowler - Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)

**Upcoming (Module 2):**
- [gRPC Basics](https://grpc.io/docs/what-is-grpc/introduction/)
- [Protocol Buffers](https://protobuf.dev/)

**Upcoming (Module 3+):**
- [Raft Visualization](http://thesecretlivesofdata.com/raft/)
- [Designing Data-Intensive Applications (Chapter 9)](https://dataintensive.net/)

---

## 📝 Important Notes

**Design Decisions:**
- Commands don't return errors (idempotent by design)
- Queries can return errors (they don't modify state)
- Validation happens in the API layer (HTTP handlers)
- Package visibility in Go is package-level, not type-level
- Commands access store internals directly (same package)

**Thread Safety:**
- All store operations are mutex-protected
- `Apply()` uses write lock (modifies state and log)
- `Execute()` uses read lock (read-only)
- `GetLog()` returns a copy to prevent external modification
- `Replay()` holds write lock and calls `cmd.Apply()` directly

**Testing Strategy:**
- Command tests: Verify logic in isolation (package `store`)
- Store tests: Verify orchestration (logging, replay, determinism)
- Server tests: Verify HTTP API (end-to-end through router)
- Always run with `-race` flag to catch concurrency issues

---

## 🎯 Immediate Next Steps

1. **Review gRPC basics** (links above)
2. **Define Protocol Buffer schema** for node communication
3. **Implement peer-to-peer RPC** between nodes
4. **Set up multi-node testing infrastructure**

Ready to dive into distributed systems! 🚀
