package world

import (
	"encoding/json"
	"log"
	"os"
)

type Tile struct {
	World    *World
	X        int
	Y        int
	Altitude float64
	Moisture float64
}

func (t *Tile) Adjacent() []*Tile {
	return t.World.Adjacent(t)
}

type TerrainType struct {
	MaxAltitude       float64 `json:"maxAltitude"`
	MoistureThreshold float64 `json:"moistureThreshold,omitempty"`
	TileDry           string  `json:"tileDry,omitempty"`
	TileWet           string  `json:"tileWet,omitempty"`
	Tile              string  `json:"tile"`
	DescriptionDry    string  `json:"descriptionDry,omitempty"`
	DescriptionWet    string  `json:"descriptionWet,omitempty"`
	Description       string  `json:"description"`
}

type TileSet struct {
	TerrainTypes []TerrainType `json:"terrainTypes"`
}

func LoadTileSet(path string) TileSet {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	var tileSet TileSet
	err = json.Unmarshal(file, &tileSet)
	if err != nil {
		log.Fatalf("Error during Unmarshal(): %s", err)
	}

	return tileSet
}

func (tileSet *TileSet) ToASCII(t *Tile) string {
	for _, terrain := range tileSet.TerrainTypes {
		if t.Altitude <= terrain.MaxAltitude {
			if terrain.MoistureThreshold != 0 {
				if t.Moisture < terrain.MoistureThreshold {
					return terrain.TileDry
				} else {
					return terrain.TileWet
				}
			}
			return terrain.Tile
		}
	}
	return "?" // Default tile if no match is found
}
