package world

import "fmt"

const (
	tileMapFp = "world/tilemap.json"
	antTypeFp = "ant/anttypes.json"
	Width     = 100
	Height    = 40
)

var OverWorld *World
var OWInitialized bool = false

func Init() {
	OverWorld = New(Width, Height)
	OverWorld.Generate()
	OWInitialized = true

	tileSet := LoadTileSet(tileMapFp)

	// Print the world
	for _, row := range OverWorld.Tiles {
		for _, tile := range row {
			fmt.Print(tileSet.ToASCII(tile))
		}
		fmt.Println()
	}
}
