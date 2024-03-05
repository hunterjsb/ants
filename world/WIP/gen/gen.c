// gcc gen.c perlin.c hashes.c -lm -o gen -I /usr/local/include/hiredis -lhiredis
#include "perlin.h"
#include "hashes.h"

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <hiredis/hiredis.h>

#define ATTRIBUTE_COUNT 5

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

void store_tile_attributes(redisContext *c, int x, int y, TileAttributes attributes) {
    char key[256];
    sprintf(key, "tile:%d:%d", x, y);
    
    redisCommand(c, "HSET %s altitude %d moisture %d temperature %d fertility %d foliage_density %d",
                 key, attributes.altitude, attributes.moisture, attributes.temperature,
                 attributes.fertility, attributes.foliage_density);
}

uint8_t* serialize_tile_attributes(TileAttributes attributes, int* size) {
    uint8_t flags = 0b11111; // Assume all data must be saved initially
    *size = 1 + ATTRIBUTE_COUNT; // 1 byte for flags + 1 byte per attribute

    uint8_t* serializedData = (uint8_t*)malloc(*size);
    if (!serializedData) {
        perror("Failed to allocate memory");
        exit(EXIT_FAILURE);
    }

    serializedData[0] = flags; // The first byte is the flag byte
    serializedData[1] = attributes.altitude;
    serializedData[2] = attributes.moisture;
    serializedData[3] = attributes.temperature;
    serializedData[4] = attributes.fertility;
    serializedData[5] = attributes.foliage_density;

    return serializedData;
}

void write_chunk_to_file(int chunkX, int chunkY, TileAttributes* attributes, int chunkSize) {
    char filename[256];
    snprintf(filename, sizeof(filename), "%d-%d.chunk", chunkX, chunkY);

    FILE* file = fopen(filename, "wb");
    if (!file) {
        perror("Failed to open file");
        exit(EXIT_FAILURE);
    }

    for (int i = 0; i < chunkSize * chunkSize; i++) {
        int serializedSize;
        uint8_t* serializedData = serialize_tile_attributes(attributes[i], &serializedSize);
        fwrite(serializedData, sizeof(uint8_t), serializedSize, file);
        free(serializedData);
    }

    fclose(file);
}

int main(int argc, char *argv[]) {
    // Initialize Redis connection
    redisContext *c = redisConnect("127.0.0.1", 6379);
    if (c == NULL || c->err) {
        if (c) {
            printf("Redis error: %s\n", c->errstr);
        } else {
            printf("Can't allocate redis context\n");
        }
        exit(1);
    }

    // Setup
    int seed = 1997;
    unsigned char** hashes = create_hash_array(2);
    int CHUNK_SIZE = 0xf;

    TileAttributes* attributesArray = (TileAttributes*)malloc(CHUNK_SIZE * CHUNK_SIZE * sizeof(TileAttributes));
    if (!attributesArray) {
        printf("Failed to allocate memory for tile attributes\n");
        exit(1);
    }

    // Generate attributes for each tile in the chunk
    for(int y = 0; y < CHUNK_SIZE; y++) {
        for(int x = 0; x < CHUNK_SIZE; x++) {
            TileAttributes attributes = generate_tile_attributes(x, y, hashes);
            attributesArray[y * CHUNK_SIZE + x] = attributes; // Store attributes in array
            store_tile_attributes(c, x, y, attributes);
        }
    }

    // Write chunk to file
    write_chunk_to_file(0, 0, attributesArray, CHUNK_SIZE); // Using top-left coords (0,0) for filename

    // Cleanup
    free(attributesArray);
    free_hashes(hashes, 2);
    redisFree(c);

    return 0;
}
