// gcc -o tilemap tilemap.c perlin.c hashes.c -lm && ./tilemap
#include "perlin.h"
#include "hashes.h"
#include <stdio.h>

typedef enum {
    WATER,
    SAND,
    GRASS,
    FOREST,
    MOUNTAIN
} TileType;

TileType get_tile_type(float perlinValue) {
    if (perlinValue < 0.3) return WATER;
    if (perlinValue < 0.4) return SAND;
    if (perlinValue < 0.6) return GRASS;
    if (perlinValue < 0.8) return FOREST;
    return MOUNTAIN;
}

int main(int argc, char *argv[]) {
    int seed = 1997;
    unsigned char** hashes = create_hash_array(2);

    int chunk_dim = 0xf;

    int x, y;
    float per;
    for (y = 0; y < chunk_dim; y++) {
        for (x = 0; x < chunk_dim; x++) {
            per = Perlin_Get2d(x, y, 0.1, 4, hashes[1]);
            TileType tile = get_tile_type(per);
            // For demonstration, print the tile type as a character
            char tileChar;
            switch (tile) {
                case WATER: tileChar = '~'; break;
                case SAND: tileChar = '.'; break;
                case GRASS: tileChar = '"'; break;
                case FOREST: tileChar = 'T'; break;
                case MOUNTAIN: tileChar = '^'; break;
                default: tileChar = '?'; // Unknown tile type
            }
            // printf(" %.3f", per);
            printf("%c", tileChar);
        }
        puts("\n");
    }

    free_hashes(hashes, 2);

    return 0;
}
