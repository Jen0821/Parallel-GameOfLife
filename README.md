# 🚀 Go-Parallel-GameOfLife (Conway's Game of Life Parallel Implementation)

This project implements **Conway's Game of Life** cellular automaton using **Go (Golang)**. The core focus is on achieving **high-performance parallel simulation** using Go's built-in concurrency primitives: **Goroutines** and **Channels**. This approach enables efficient simulation and real-time visualization on large-scale grids.

---

## 🎯 Key Technologies and Implementation

The simulation utilizes a **Distributor/Worker model** to divide the grid (board) into sections, distributing the computation load across multiple Goroutines running on a single machine.

| Feature | Technology / Concept | Description |
| :--- | :--- | :--- |
| **Primary Language** | Go (Golang) | Used for entire application logic. |
| **Concurrency** | **Goroutines** & **Channels** | Implements the **Distributor/Worker** pattern for parallel updates. |
| **Visualization** | **SDL (Simple DirectMedia Layer)** | Handles the real-time graphical display of the grid state. |
| **Data I/O** | **PGM (Portable Graymap)** | Used for loading initial states and saving final/intermediate states. |
| **Domain** | **Closed Domain (Toroidal)** | Implements toric boundary conditions (pixels on opposite edges are connected). |

---

## Initial State Example

The simulation starts from an initial PGM image representing the live and dead cells on the grid.

Here's an example of an initial board state:


---

## ⚙️ Implemented Features (Step-by-Step)

The project was developed following the coursework guidelines, integrating serial implementation with progressive parallel features, I/O, and user control.

### 1. Parallel Core Logic (Steps 1 & 2)

* **Serial Baseline:** Initial single-threaded implementation of the Game of Life rules.
* **Parallelization:** Implementation of the **Distributor** and **Worker** model. The Distributor divides the board into stripes and assigns them to a pool of **Worker Goroutines** (as specified by `gol.Params.Threads`) to calculate the next state in parallel.

### 2. State Reporting and Events (Step 3)

* **Alive Count Ticker:** Uses a **Ticker** to report the total number of **Alive Cells** via the `AliveCellsCount` event every **2 seconds**. This provides real-time feedback on the simulation's activity.

### 3. Image Output (Step 4)

* Implements logic to output the final state of the board as a **PGM image** after all turns have been completed. This allows for persistent storage of simulation results.

### 4. User Control Rules (Step 5)

Implemented interactive keyboard controls processed by the main event loop:

* **`s` (Save):** Saves the current board state as a PGM image (`ImageOutputComplete` event).
* **`q` (Quit):** Completes the current turn, saves the final state as a PGM image, and terminates the program (`FinalTurnComplete` event).
* **`p` (Pause/Resume):** Toggles the simulation state between running and paused (`StateChange` event).

### 5. Real-Time Visualization (Step 6)

* Integration with **SDL** to display the simulation in real-time within a dedicated window.
* Utilizes **`CellFlipped`** and **`TurnComplete`** events to manage graphical updates efficiently.
    
    Here's an example of the real-time visualization:
    

---

## ▶️ Running and Testing

To run the implementation and tests (as suggested in the coursework):

```bash
# Run the program with SDL visualization (Step 6)
go run .

# Run tests with the SDL window to test visualization and keyboard control (Step 5 & 6)
# Ensure SDL development libraries are installed on your system.
go test ./tests -v -run TestKeyboard -sdl

# Run tests with race detector for thread safety checks (Step 6)
go test ./tests -v -race
