package main

import (
	"fmt"

	ant "ants/ant"
	"ants/world"
)

const (
	width     = 100
	height    = 40
	tileMapFp = "world/tilemap.json"
	antTypeFp = "ant/anttypes.json"
)

var w *world.World

func initWorld() {
	w = world.New(width, height)
	w.Generate()

	tileSet := world.LoadTileSet(tileMapFp)

	// Print the world
	for _, row := range w.Tiles {
		for _, tile := range row {
			fmt.Print(tileSet.ToASCII(tile))
		}
		fmt.Println()
	}
}

func main() {
	initWorld()
	queen := ant.Ant{Tile: w.Tiles[1][1], Type: ant.Queen, MoveSpeed: 0}
	worker := queen.Spawn(ant.Worker)

	//e := worker.Move(w.Tiles[2][2])
	fmt.Println(worker.Tile.X, worker.Tile.Y)
}
