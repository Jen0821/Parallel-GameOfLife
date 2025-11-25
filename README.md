# 🚀 Go-Parallel-GameOfLife

This project implements Conway's Game of Life cellular automaton using Go (Golang).  
It focuses on achieving **high-performance parallel simulation** using Goroutines & Channels, and extends into a **fully distributed architecture** using a Controller/Broker/Worker model with RPC networking.

---

## 🌟 Key Features & Highlights

### ⚡ High-Performance Parallel Core
- Distributor/Worker parallel pattern using Goroutines
- Efficient stripe-based board division for multi-core CPUs

### 🌐 Scalable Distributed System
- Local Controller ↔ Broker ↔ Worker Nodes  
- Remote workers (simulated AWS nodes) compute board slices via RPC

### 🖥 Real-Time Visualization
- SDL2-based live rendering of the simulation

### 🔄 Robust State Management
- Toroidal (closed) grid domain  
- Efficient change reporting via channels/events

### 🎮 User Interaction
- Pause / Resume  
- Save intermediate state  
- Graceful termination

---

## 💡 Overview: Conway's Game of Life

Conway’s Game of Life is a zero-player automaton where each generation evolves based solely on the previous state.

### Rules
1. Any live cell with < 2 live neighbours → **dies** (underpopulation)  
2. Any live cell with 2 or 3 neighbours → **lives**  
3. Any live cell with > 3 neighbours → **dies** (overpopulation)  
4. Any dead cell with exactly 3 neighbours → **becomes alive**

Our implementation parallelizes and distributes these computations for maximum throughput.

---

## 🧱 Key Technologies and Architecture

The system has **two layers**:

### **1. Local Parallel Layer**
- Distributor/Worker goroutine model
- Multi-core optimized turn computation

### **2. Distributed Network Layer**
- Controller/Broker/Worker nodes over RPC
- Sliced board distribution
- Halo exchange communication for neighbor consistency

| Feature | Technology / Concept | Description |
|--------|----------------------|-------------|
| Primary Language | Go (Golang) | Full project implemented in Go |
| Concurrency | Goroutines & Channels | Local parallel Distributor/Worker model |
| Distribution | RPC networking | Controller/Broker/Worker communication |
| System Model | Controller/Broker/Worker | Scales across machines/nodes |
| Visualization | SDL2 | Real-time graphics |
| Data I/O | PGM images | Load/save board states |
| Domain | Toroidal Grid | Wrap-around edges |

---

## 🏗 Distributed System Architecture

### **Local Controller**
- Handles I/O  
- Captures keypresses (`s`, `p`, `q`)  
- Displays board via SDL  
- Initiates game by connecting to Broker

### **Broker**
- Central orchestrator  
- Receives commands from Controller  
- Slices the global board  
- Dispatches tasks to Worker Nodes  
- Aggregates turn results

### **GOL Worker (Compute Node)**
- Computes assigned board slice  
- Uses goroutine parallelism internally  
- Performs Halo Exchange with neighbor workers  
- Returns computed slice via RPC

---

## 📈 Performance and Scalability

### Parallel Efficiency (Local)
- Distributor/Worker model uses all CPU cores
- Significant speedup vs serial baseline

### Distributed Scalability
- Horizontal scaling via remote worker nodes
- Turn-processing time decreases as N workers increases
- Reduced network cost via halo-only communication

### Fault Tolerance (Considered)
- System should maintain state if a new Controller reconnects  
- Worker failure handling out of current scope but considered

---

## ⚙️ Implemented Features (Development Stages)

### 1. **Parallel Core Logic**
- Serial → parallel migration  
- Board divided into stripes for worker goroutines  

### 2. **State Reporting & Events**
- AliveCellsCount sent every 2 seconds  
- RPC communication between Controller/Broker/Workers  

### 3. **Image Output & Persistent State**
- Load initial PGM  
- Save intermediate/final PGM images  

### 4. **User Control (Interactive)**
- `s` → save state  
- `q` → complete turn & save final image, then exit  
- `p` → pause / resume  

### 5. **Real-Time Visualization (SDL)**
- CellFlipped events for single-pixel changes  
- TurnComplete to refresh entire grid

---

## ▶️ Setup and Running

### **Prerequisites**
Install Go and SDL development libraries.

````markdown
### macOS (Homebrew)
```bash
brew install sdl2
```

### Linux (Ubuntu/Debian)
```bash
sudo apt-get install libsdl2-dev
```

### Windows
```bash
# Requires MinGW installation and manual SDL2 linking
# Refer to the Go SDL documentation for detailed platform-specific setup
```

## 🚀 Running

### Run the program (Controller/Broker/Workers assumed running or mocked)
```bash
go run .
```

### Test visualization + keyboard controls
```bash
go test ./tests -v -run TestKeyboard -sdl
```

### Test parallel core with race detector
```bash
go test ./tests -v -race
```
