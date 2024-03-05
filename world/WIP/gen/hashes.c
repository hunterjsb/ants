#include "hashes.h"

#include <stdio.h>
#include <stdlib.h>
#include <time.h>

// for the love of God do not forget to free the memory to allocate to the hash table

unsigned char* generate_random_hash_table() {
    unsigned char* hash_table = (unsigned char*)malloc(HASH_SIZE * sizeof(unsigned char));
    if (hash_table == NULL) {
        fprintf(stderr, "Memory allocation failed\n");
        exit(EXIT_FAILURE);
    }

    // Initialize hash table with values 0-255
    for (int i = 0; i < HASH_SIZE; ++i) {
        hash_table[i] = i;
    }

    // Shuffle using Fisher-Yates shuffle
    for (int i = HASH_SIZE - 1; i > 0; i--) {
        int j = rand() % (i + 1);

        // Swap
        unsigned char temp = hash_table[i];
        hash_table[i] = hash_table[j];
        hash_table[j] = temp;
    }

    return hash_table;
}

unsigned char** create_hash_array(int numHashes) {
    unsigned char** hashes = (unsigned char**)malloc(numHashes * sizeof(unsigned char*));
    if (hashes == NULL) {
        fprintf(stderr, "Memory allocation failed\n");
        exit(EXIT_FAILURE);
    }

    for (int i = 0; i < numHashes; ++i) {
        hashes[i] = generate_random_hash_table();  // Generate and assign a new hash table
    }

    return hashes;
}

void free_hashes(unsigned char** hashes, int num_hashes) {
    // Cleanup
    for (int i = 0; i < num_hashes; ++i) {
        free(hashes[i]);  // Free each individual hash table
    }
    free(hashes);  // Free the array of pointers
}