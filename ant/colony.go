package ant

import (
	"ants/world"
	"time"
)

const DEFAULT_SPAWN_RATE time.Duration = 10 * time.Second

type Colony struct {
	Ants          []*Ant        `json:"ants,omitempty"`
	Queen         *Ant          `json:"queen"`
	Owner         *User         `json:"owner"`
	LastSpawnTime time.Time     `json:"lastSpawnTime"`
	SpawnRate     time.Duration `json:"spawnRate"`
}

func NewColony(user *User, tile *world.Tile) *Colony {
	colony := &Colony{
		Owner:         user,
		LastSpawnTime: time.Now(), // Set the last spawn time to now
		SpawnRate:     DEFAULT_SPAWN_RATE,
	}
	// Create and assign the queen
	queen := NewQueen(colony, tile)
	colony.Queen = queen
	colony.Ants = append([]*Ant{}, queen) // Initialize with queen

	return colony
}

func (c *Colony) Spawn(t AntType) (*Ant, *AntError) {
	return c.Queen.Spawn(t)
}
