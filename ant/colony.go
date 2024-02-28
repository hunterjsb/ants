package ant

import (
	"ants/world"
)

type Colony struct {
	Ants  []*Ant `json:"ants,omitempty"`
	Queen *Ant   `json:"queen"`
	Owner *User  `json:"owner"`
}

func NewColony(user *User, tile *world.Tile) *Colony {
	colony := &Colony{Owner: user}
	queen := NewQueen(colony, tile)
	colony.Queen = queen
	colony.Ants = make([]*Ant, 1)
	colony.Ants = append(colony.Ants, queen)
	return colony
}
