package ant

import (
	"ants/world"
	"time"
)

type Colony struct {
	Ants          []*Ant    `json:"ants,omitempty"`
	Queen         *Ant      `json:"queen"`
	Owner         *User     `json:"owner"`
	LastSpawnTime time.Time `json:"lastSpawnTime"`
}

func NewColony(user *User, tile *world.Tile) *Colony {
	colony := &Colony{Owner: user, LastSpawnTime: time.Now()} // Set the last spawn time to now
	queen := NewQueen(colony, tile)
	colony.Queen = queen
	colony.Ants = append([]*Ant{}, queen) // Initialize with queen
	return colony
}
