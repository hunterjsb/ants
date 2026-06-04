# Ants

A simulation of life, labor, and expansion.

## Overview

`Ants` is a distributed colony simulation built in Go. It manifests a procedurally generated world where queens establish their reign, and colonies thrive through the coordinated efforts of workers, scouts, and soldiers.

## Architecture

- **The World**: A grid of tiles generated via Perlin noise, defining terrain from abyssal waters to snow-capped peaks.
- **The Colony**: A hierarchical structure centered around a Queen.
- **The API**: A RESTful gateway to the hive, facilitating colony initiation and ant gestation.
- **Persistence**: Powered by Redis, ensuring the state of the swarm is preserved.

## Technical Deep Dive

### Core Tenets

1.  **Performance-First Generation**: The world should be vast, yet generation must remain lightning-fast.
2.  **Noise as a Baseline**: Procedural noise (Perlin) serves as the "source of truth." Only deviations (deltas) from this noise need to be persisted, significantly reducing storage overhead.
3.  **Language Agnosticism**: Heavy compute (world generation) should live where it is fastest (C), while high-level logic and orchestration thrive in Go.

### Technical Goals

-   **High-Density Persistence**: Minimize the memory and disk footprint of millions of tiles.
-   **Seamless Scalability**: Enable the simulation of massive worlds through chunk-based loading and distributed architecture.
-   **Low-Latency Interactivity**: Ensure colony management and ant movement remain responsive even under heavy load.

### Approaches & Implementation

-   **CGO & C-Level Optimization**: Transitioning procedural generation to C (`world/WIP/gen`) to leverage raw performance and mature math libraries.
-   **Delta-Encoded Serialization**: A specialized binary format where a single byte of flags indicates which attributes (altitude, moisture, temperature, etc.) differ from the noise baseline. Only "mutated" data is stored.
-   **Bit-Packed Attributes**: Using compact data types and bit-fields to represent tile and entity states within the world chunks.
-   **Redis as a Sparse State Store**: Utilizing Redis not just for ephemeral cache, but as a robust store for the persistent deltas of the world state.

## Current State: The Great Transition

The project is currently in a state of evolution. While the core logic resides in Go, an ambitious migration towards C-based world generation and chunked serialization is underway in `world/WIP`. This shift aims for high-performance tile management and dense attribute storage.

## Getting Started

1. Ensure a Redis instance is accessible.
2. Configure environment variables in `.env`.
3. Run the swarm:
   ```bash
   go run main.go
   ```

---
*The hive never sleeps.*
