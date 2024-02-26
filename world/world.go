package world

import (
	"math/rand"
	"time"

	"github.com/aquilax/go-perlin"
)

type World struct {
	Tiles  [][]*Tile
	Width  int
	Height int
}

func New(width, height int) *World {
	w := &World{
		Tiles:  make([][]*Tile, height),
		Width:  width,
		Height: height,
	}
	for i := range w.Tiles {
		w.Tiles[i] = make([]*Tile, width)
	}
	return w
}

func (w *World) Generate() {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	pAltitude := perlin.NewPerlin(2, 2, 3, rng.Int63())
	pMoisture := perlin.NewPerlin(2, 2, 3, rng.Int63())

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			altitudeValue := pAltitude.Noise2D(float64(x)/10, float64(y)/10)
			moistureValue := pMoisture.Noise2D(float64(x)/10, float64(y)/10)

			altitude := scaleAltitude(altitudeValue)
			moisture := scaleMoisture(moistureValue)

			w.Tiles[y][x] = &Tile{
				World:    w,
				X:        x,
				Y:        y,
				Altitude: altitude,
				Moisture: moisture,
			}
		}
	}
}

// Adjacent returns a slice of pointers to all adjacent tiles.
func (w *World) Adjacent(t *Tile) []*Tile {
	var adjacentTiles []*Tile

	// Coordinates of the tile
	x, y := t.X, t.Y

	// Check each of the four directions, ensuring we don't go out of bounds
	if x > 0 {
		adjacentTiles = append(adjacentTiles, w.Tiles[y][x-1]) // Left
	}
	if x < w.Width-1 {
		adjacentTiles = append(adjacentTiles, w.Tiles[y][x+1]) // Right
	}
	if y > 0 {
		adjacentTiles = append(adjacentTiles, w.Tiles[y-1][x]) // Up
	}
	if y < w.Height-1 {
		adjacentTiles = append(adjacentTiles, w.Tiles[y+1][x]) // Down
	}

	return adjacentTiles
}

func scaleAltitude(noiseValue float64) float64 {
	return (noiseValue + 1) * 500
}

func scaleMoisture(noiseValue float64) float64 {
	return (noiseValue + 1) / 2
}
