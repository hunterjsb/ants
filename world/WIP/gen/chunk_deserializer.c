#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

typedef struct {
    int32_t x, y; // Chunk coordinates
    uint8_t* tiles; // Pointer to tiles data
    size_t tiles_count; // Number of tiles
} Chunk;

Chunk* deserialize_chunk(const uint8_t* data, size_t data_size) {
    if (data_size < 8) { // Minimum size check
        return NULL;
    }
    
    // Allocate memory for chunk
    Chunk* chunk = (Chunk*)malloc(sizeof(Chunk));
    if (!chunk) {
        return NULL; // Memory allocation failed
    }
    
    // Deserialize coordinates
    chunk->x = (data[0] << 24) | (data[1] << 16) | (data[2] << 8) | data[3];
    chunk->y = (data[4] << 24) | (data[5] << 16) | (data[6] << 8) | data[7];
    
    // Remaining data represents tiles
    chunk->tiles_count = data_size - 8;
    chunk->tiles = (uint8_t*)malloc(chunk->tiles_count);
    if (!chunk->tiles) {
        free(chunk); // Cleanup chunk allocation
        return NULL; // Memory allocation for tiles failed
    }
    
    // Copy tile data
    for (size_t i = 0; i < chunk->tiles_count; i++) {
        chunk->tiles[i] = data[i + 8];
    }
    
    return chunk;
}

void free_chunk(Chunk* chunk) {
    if (chunk) {
        if (chunk->tiles) {
            free(chunk->tiles); // Free tiles memory
        }
        free(chunk); // Free chunk memory
    }
}

int main() {
    // Example usage
    uint8_t example_chunk_data[] = {
        0x00, 0x00, 0x00, 0x01, // X coordinate
        0x00, 0x00, 0x00, 0x02, // Y coordinate
        // Tile data...
    };
    size_t data_size = sizeof(example_chunk_data) / sizeof(example_chunk_data[0]);
    
    Chunk* chunk = deserialize_chunk(example_chunk_data, data_size);
    if (chunk) {
        printf("Chunk (%d, %d) with %zu tiles\n", chunk->x, chunk->y, chunk->tiles_count);
        // Process chunk...
        
        free_chunk(chunk);
    }
    
    return 0;
}
