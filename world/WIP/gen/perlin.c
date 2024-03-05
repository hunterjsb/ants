#include "perlin.h"
#include "hashes.h"

#include <math.h>
#include <stdio.h>
#include <stdlib.h>

static int noise2(int x, int y, unsigned char* hash)
{
    int yindex = y % HASH_SIZE;
    if (yindex < 0)
        yindex += HASH_SIZE; // Correct negative indices
    int xindex = (hash[yindex] + x) % HASH_SIZE; // Use HASH_SIZE for wrapping
    if (xindex < 0)
        xindex += HASH_SIZE; // Ensure positive indices
    const int result = hash[xindex];
    return result;
}

static double lin_inter(double x, double y, double s)
{
    return x + s * (y-x);
}

static double smooth_inter(double x, double y, double s)
{
    return lin_inter( x, y, s * s * (3-2*s) );
}

static double noise2d(double x, double y, unsigned char* hash)
{
    const int  x_int = floor( x );
    const int  y_int = floor( y );
    const double  x_frac = x - x_int;
    const double  y_frac = y - y_int;
    const int  s = noise2( x_int, y_int, hash );
    const int  t = noise2( x_int+1, y_int, hash );
    const int  u = noise2( x_int, y_int+1, hash );
    const int  v = noise2( x_int+1, y_int+1, hash );
    const double  low = smooth_inter( s, t, x_frac );
    const double  high = smooth_inter( u, v, x_frac );
    const double  result = smooth_inter( low, high, y_frac );
    return result;
}

double Perlin_Get2d(double x, double y, double freq, int depth, unsigned char* hash)
{
    double  xa = x*freq;
    double  ya = y*freq;
    double  amp = 1.0;
    double  fin = 0;
    double  div = 0.0;
    for (int i=0; i<depth; i++)
    {
        div += 256 * amp;
        fin += noise2d( xa, ya, hash ) * amp;
        amp /= 2;
        xa *= 2;
        ya *= 2;
    }
    return fin/div;
}
