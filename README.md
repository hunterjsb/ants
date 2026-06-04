# Ants

A simulation of life, labor, and expansion.

## Overview

`Ants` is a distributed colony simulation built in Go. It manifests a procedurally generated world where queens establish their reign, and colonies thrive through the coordinated efforts of workers, scouts, and soldiers.

## Architecture

- **The World**: A grid of tiles generated via Perlin noise, defining terrain from abyssal waters to snow-capped peaks.
- **The Colony**: A hierarchical structure centered around a Queen.
- **The API**: A RESTful gateway to the hive, facilitating colony initiation and ant gestation.
- **Persistence**: Powered by Redis, ensuring the state of the swarm is preserved.

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
