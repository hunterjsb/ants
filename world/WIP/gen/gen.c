// gcc gen.c perlin.c hashes.c -lm -o gen -I /usr/local/include/hiredis -lhiredis
#include "perlin.h"
#include "hashes.h"

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <hiredis/hiredis.h>

#define ATTRIBUTE_COUNT 5
#define CHUNK_SIZE 0xF+1 

typedef struct {
    uint8_t altitude;
    uint8_t moisture;
    uint8_t temperature;
    uint8_t fertility;
    uint8_t foliage_density;
} TileAttributes;


TileAttributes generate_tile_attributes(double x, double y, unsigned char** hashes) {
    TileAttributes attributes;
    
    double altitudeNoise = Perlin_Get2d(x, y, 0.1, 4, hashes[0]);
    double moistureNoise = Perlin_Get2d(x, y, 0.1, 4, hashes[1]);
    
    // Normalize and scale noise values to 0-255
    attributes.altitude = (uint8_t)((altitudeNoise) * 256);
    attributes.moisture = (uint8_t)((moistureNoise) * 256);
    
    // Simplified calculations for demo purposes
    attributes.temperature = 255 - attributes.altitude; // Higher altitude, lower temperature
    attributes.fertility = attributes.moisture; // Higher moisture, higher fertility
    attributes.foliage_density = (attributes.moisture + attributes.temperature) / 2; // Depends on moisture and temp
    
    return attributes;
}

void write_tile_redis(redisContext *c, int x, int y, TileAttributes attributes) {
    char key[256];
    sprintf(key, "tile:%d:%d", x, y);
    
    redisCommand(c, "HSET %s altitude %d moisture %d temperature %d fertility %d foliage_density %d",
                 key, attributes.altitude, attributes.moisture, attributes.temperature,
                 attributes.fertility, attributes.foliage_density);
}

uint8_t* serialize_tile_attributes(TileAttributes attributes, TileAttributes perlin_attributes, int* size) {
    uint8_t flags = 0; // Start with all flags set to 0 (assume data matches Perlin noise)
    *size = 1; // Start with 1 byte for flags

    // Buffer to hold serialized data, initially allocate 1 byte for flags + 5 bytes for potential data
    uint8_t* buffer = (uint8_t*)malloc(6);
    if (!buffer) {
        perror("Failed to allocate memory");
        exit(EXIT_FAILURE);
    }

    // Flag and serialize each attribute only if it differs from its Perlin-generated counterpart
    if (attributes.altitude != perlin_attributes.altitude) {
        flags |= (1 << 4); // Set flag for altitude
        buffer[*size] = attributes.altitude;
        (*size)++;
    }
    if (attributes.moisture != perlin_attributes.moisture) {
        flags |= (1 << 3); // Set flag for moisture
        buffer[*size] = attributes.moisture;
        (*size)++;
    }
    if (attributes.temperature != perlin_attributes.temperature) {
        flags |= (1 << 2); // Set flag for temperature
        buffer[*size] = attributes.temperature;
        (*size)++;
    }
    if (attributes.fertility != perlin_attributes.fertility) {
        flags |= (1 << 1); // Set flag for fertility
        buffer[*size] = attributes.fertility;
        (*size)++;
    }
    if (attributes.foliage_density != perlin_attributes.foliage_density) {
        flags |= 1; // Set flag for foliage_density
        buffer[*size] = attributes.foliage_density;
        (*size)++;
    }

    // Reallocate buffer to match actual data size
    uint8_t* serialized_data = realloc(buffer, *size);
    if (!serialized_data) {
        perror("Failed to reallocate memory");
        free(buffer); // Ensure to free original buffer on failure
        exit(EXIT_FAILURE);
    }

    // Insert flags byte at the beginning
    serialized_data[0] = flags;

    return serialized_data;
}


void write_chunk_to_file(int chunk_x, int chunk_y, TileAttributes* attributes, unsigned char** hashes) {
    char filename[256];
    snprintf(filename, sizeof(filename), "%d-%d.chunk", chunk_x, chunk_y);

    FILE* file = fopen(filename, "wb");
    if (!file) {
        perror("Failed to open file");
        exit(EXIT_FAILURE);
    }

    for (int i = 0; i < CHUNK_SIZE * CHUNK_SIZE; i++) {
        // Calculate the coordinates of the tile within the chunk
        int tileX = chunk_x * CHUNK_SIZE + (i % CHUNK_SIZE);
        int tileY = chunk_y * CHUNK_SIZE + (i / CHUNK_SIZE);

        // Regenerate attributes from Perlin noise for comparison
        TileAttributes noiseAttributes = generate_tile_attributes((double)tileX, (double)tileY, hashes);

        int serializedSize;
        // Pass both actual attributes and noise-generated attributes for serialization
        uint8_t* serialized_data = serialize_tile_attributes(attributes[i], noiseAttributes, &serializedSize);
        fwrite(serialized_data, sizeof(uint8_t), serializedSize, file);
        free(serialized_data);
    }

    fclose(file);
}

void _new_chunk_write(int chunk_x, int chunk_y) {
    char filename[256];
    snprintf(filename, sizeof(filename), "%d-%d.chunk", chunk_x, chunk_y);

    FILE* file = fopen(filename, "wb");
    if (!file) {
        perror("Failed to open file");
        exit(EXIT_FAILURE);
    }

    // Create a buffer of 256 bytes set to zero
    unsigned char emptyBuffer[256];
    memset(emptyBuffer, 0, sizeof(emptyBuffer)); // Initialize buffer with zeros

    // Write the empty buffer to the file
    fwrite(emptyBuffer, sizeof(emptyBuffer), 1, file);
    fclose(file);
}

void new_chunk(redisContext *c, unsigned char** hashes) {
    // Setup
    int seed = 1997;

    // Generate attributes for each tile in the chunk
    for(int y = 0; y < CHUNK_SIZE; y++) {
        for(int x = 0; x < CHUNK_SIZE; x++) {
            TileAttributes attributes = generate_tile_attributes(x, y, hashes);
            write_tile_redis(c, x, y, attributes);
        }
    }

    // new empty chunk at coords
    _new_chunk_write(0, 0);
}

redisContext* redis_connect() {
    redisContext *c = redisConnect("127.0.0.1", 6379);
    if (c == NULL || c->err) {
        if (c) {
            printf("Redis error: %s\n", c->errstr);
        } else {
            printf("Can't allocate redis context\n");
        }
        exit(1);
    }
    return c;
}

int main(int argc, char *argv[]) {
    redisContext *c = redis_connect();
    unsigned char** hashes = create_hash_array(2);
    new_chunk(c, hashes);

    TileAttributes* attributes_array = (TileAttributes*)malloc(CHUNK_SIZE * CHUNK_SIZE * sizeof(TileAttributes));
    if (!attributes_array) {
        printf("Failed to allocate memory for tile attributes\n");
        exit(1);
    }

    // Cleanup
    free(attributes_array);
    free_hashes(hashes, 2);
    redisFree(c);

    return 0;
}
