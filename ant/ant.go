package ant

import (
	"ants/world"
	"encoding/json"
	"errors"
	"log"
	"math"
	"os"
)

const antTypesFp string = "ant/anttypes.json"

type AntType string

const (
	Worker  AntType = "Worker"
	Soldier AntType = "Soldier"
	Scout   AntType = "Scout"
	Queen   AntType = "Queen"
)

type Ant struct {
	Tile      *world.Tile `json:"tile"`
	Type      AntType     `json:"type"`
	MoveSpeed int         `json:"moveSpeed"`
	Attack    int         `json:"attack"`
	Defense   int         `json:"defense"`
	HP        int         `json:"hp"`
	Colony    *Colony     `json:"-"`
}

var antTypeConfig AntTypeConfig = loadAntTypeConfig(antTypesFp)

func NewQueen(c *Colony, tile *world.Tile) *Ant {
	props, exists := antTypeConfig[string(Queen)]
	if !exists {
		log.Fatal("Queen not defined in config")
	}
	return &Ant{
		Colony:    c,
		Tile:      tile,
		Type:      Queen,
		MoveSpeed: props.MoveSpeed,
		Attack:    props.Attack,
		Defense:   props.Defense,
		HP:        props.HP,
	}
}

func (a *Ant) Spawn(t AntType) *Ant {
	if a.Type == Queen && t != Queen {
		props, exists := antTypeConfig[string(t)]
		if !exists {
			log.Fatalf("Ant type '%s' not found in configuration", t)
		}
		return &Ant{
			Colony:    a.Colony,
			Tile:      a.Tile,
			Type:      t,
			MoveSpeed: props.MoveSpeed,
			Attack:    props.Attack,
			Defense:   props.Defense,
			HP:        props.HP,
		}
	}
	return nil
}

func (a *Ant) Adjacent() []*world.Tile {
	return a.Tile.Adjacent()
}

func (a *Ant) Move(t *world.Tile) error {
	dist := math.Sqrt(float64((a.Tile.X-t.X)*(a.Tile.X-t.X) + (a.Tile.Y-t.Y)*(a.Tile.Y-t.Y)))
	if dist > float64(a.MoveSpeed) {
		return errors.New("move too far")
	}
	a.Tile = t
	return nil
}

func (a *Ant) GetType() string {
	return string(a.Type)
}

type AntTypeProperties struct {
	MoveSpeed int `json:"MoveSpeed"`
	Attack    int `json:"Attack"`
	Defense   int `json:"Defense"`
	HP        int `json:"HP"`
}

type AntTypeConfig map[string]AntTypeProperties

func loadAntTypeConfig(filename string) AntTypeConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	var config AntTypeConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %s", err)
	}

	return config
}
