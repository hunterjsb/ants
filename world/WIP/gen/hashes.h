// File: hashes.h

#ifndef HASHES_H
#define HASHES_H

#define HASH_SIZE 256

unsigned char** create_hash_array(int numHashes);
void free_hashes(unsigned char** hashes, int num_hashes);

#endif // HASHES_H
