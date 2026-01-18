# Distributed Key-Value Store Learning Project

A Go implementation of a distributed key-value store with consensus, built step-by-step to learn distributed systems concepts.

## Progress Overview

### ✅ Exercise 1.1: Single Node KV Store (COMPLETED)

**What we built:**
- Thread-safe in-memory key-value store
- HTTP API with GET, PUT, DELETE operations
- Comprehensive tests using `httptest`

**Key files:**
- `store/store.go` - Original KVStore with mutex-protected map
- `server/server.go` - HTTP handlers for CRUD operations
- Tests for both store and server layers

**What we learned:**
- Go concurrency primitives (`sync.RWMutex`)
- HTTP server implementation with `ServeMux`
- Testing HTTP handlers through the router
- Thread-safety and race conditions

---

### 🚧 Exercise 1.2: Command Pattern & State Machine (IN PROGRESS)

**What we've built so far:**

#### Architecture Changes:
- **Command Pattern**: Separated operations into command objects
- **CQRS**: Split commands (writes) from queries (reads)
  - `Command` interface - stateful operations (Put, Delete)
  - `Query` interface - stateless operations (Get)

#### Current Implementation:

**Files:**
```
store/
├── store.go         # KVStore with Apply() and Execute() methods
├── commands.go      # Command/Query interfaces and implementations
├── store_test.go    # Integration tests (logging, replay)
└── commands_test.go # Unit tests (command logic)
```

**Core Concepts Implemented:**
- ✅ Commands are immutable (private fields)
- ✅ Commands directly modify store data (no proxying)
- ✅ Idempotent operations (safe to apply multiple times)
- ✅ Constructor functions (`NewPutCommand`, `NewGetQuery`, etc.)
- ✅ Separation of concerns (commands vs store orchestration)

**Store API:**
```go
// For writes (logged)
func (s *KVStore) Apply(cmd Command)

// For reads (not logged)
func (s *KVStore) Execute(query Query) (string, error)
```

**Commands Implemented:**
- `PutCommand` - Add/update key-value pairs
- `DeleteCommand` - Remove keys
- `GetQuery` - Retrieve values by key

---

### 📋 What's Left for Exercise 1.2

#### 1. **Implement Command Logging**
```go
type KVStore struct {
    mu   sync.RWMutex
    data map[string]string
    log  []Command  // Need to implement this
}

func (s *KVStore) Apply(cmd Command) {
    s.mu.Lock()
    defer s.mu.Unlock()

    cmd.Apply(s)
    s.log = append(s.log, cmd)  // ✅ Add this
}
```

#### 2. **Implement Replay Functionality**
```go
func (s *KVStore) Replay() {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Clear state
    s.data = make(map[string]string)

    // Replay all commands from log
    for _, cmd := range s.log {
        cmd.Apply(s)
    }
}
```

#### 3. **Add GetLog() Method**
```go
func (s *KVStore) GetLog() []Command {
    s.mu.RLock()
    defer s.mu.RUnlock()

    logCopy := make([]Command, len(s.log))
    copy(logCopy, s.log)
    return logCopy
}
```

#### 4. **Complete Testing**

**Store tests (`store_test.go`):**
- [ ] Test that `Apply()` appends to log
- [ ] Test that `Execute()` does NOT append to log
- [ ] Test `Replay()` rebuilds state correctly
- [ ] Test determinism (two stores with same commands = same state)

**Command tests (`commands_test.go`):**
- [x] Test PutCommand modifies data
- [x] Test GetQuery returns correct value
- [x] Test DeleteCommand removes key

#### 5. **Update HTTP Server**

Change handlers to use commands:
```go
// Before (Exercise 1.1)
func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
    s.store.Put(req.Key, req.Value)
}

// After (Exercise 1.2)
func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
    cmd := store.NewPutCommand(req.Key, req.Value)
    s.store.Apply(cmd)
}
```

---

### 🎯 Next Steps (Immediate)

1. **Add logging to `Apply()` method**
2. **Implement `Replay()` function**
3. **Write tests for logging and replay**
4. **Update HTTP handlers to use commands**
5. **Test the complete system end-to-end**

---

### 🔜 Upcoming Exercises

**Exercise 1.3:** (Not started)
- Understanding the exercise requirements

**Module 2:** Multi-Node Setup
- Node-to-node communication
- Migrate from HTTP to gRPC
- Naive replication (before consensus)

**Module 3:** Raft Consensus
- Leader election
- Log replication
- Network partition handling

---

### 📚 Key Learnings So Far

**Exercise 1.1:**
- Building thread-safe data structures
- HTTP server patterns in Go
- Testing strategies (unit vs integration)

**Exercise 1.2 (in progress):**
- Command Pattern for encapsulating operations
- State machine design
- Idempotency in distributed systems
- CQRS (Command Query Responsibility Segregation)
- Why commands shouldn't return errors

---

### 🏃 Quick Start

**Run tests:**
```bash
go test ./store/...
go test ./server/...
```

**Run server:**
```bash
go run main.go
```

**Test with curl:**
```bash
# PUT
curl -X PUT http://localhost:8080/kv \
  -H "Content-Type: application/json" \
  -d '{"key":"name","value":"alice"}'

# GET
curl http://localhost:8080/kv?key=name

# DELETE
curl -X DELETE http://localhost:8080/kv/name/
```

---

### 📝 Notes

- Package visibility in Go is at the package level (not type level)
- Commands access `store.data` directly (same package = full access)
- Queries don't modify state so they're not logged
- Commands must be deterministic for distributed consensus
