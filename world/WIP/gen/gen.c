// gcc gen.c perlin.c -lm -o gen -I /usr/local/include/hiredis -lhiredis
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "perlin.h"
#include <hiredis/hiredis.h>

typedef struct {
    uint8_t altitude;
    uint8_t moisture;
    uint8_t temperature;
    uint8_t fertility;
    uint8_t foliage_density;
} TileAttributes;


TileAttributes generate_tile_attributes(double x, double y) {
    TileAttributes attributes;
    
    double altitudeNoise = Perlin_Get2d(x, y, 0.1, 4);
    double moistureNoise = Perlin_Get2d(x, y, 0.1, 4);
    
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


int main(int argc, char *argv[]) {
    // Connect to Redis
    redisContext *c = redisConnect("127.0.0.1", 6379);
    if (c == NULL || c->err) {
        if (c) {
            printf("Redis error: %s\n", c->errstr);
        } else {
            printf("Can't allocate redis context\n");
        }
        exit(1);
    }

    int seed = 1997;
    int chunk_size = 0xf;
    initialize_hash(seed);
    
    // Terrain generation loop
    for(int y = 0; y < chunk_size; y++) { // Example: 256x256 terrain
        for(int x = 0; x < chunk_size; x++) {
            TileAttributes attributes = generate_tile_attributes(x, y);
            store_tile_attributes(c, x, y, attributes);
        }
    }
    
    // Clean up
    redisFree(c);
    return 0;
}

