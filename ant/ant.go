package ant

import (
	"ants/world"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

var antTypeConfig AntTypeConfig = loadAntTypeConfig("ant/anttypes.json")

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

func (a *Ant) Move(t *world.Tile) *AntError {
	dist := math.Sqrt(float64((a.Tile.X-t.X)*(a.Tile.X-t.X) + (a.Tile.Y-t.Y)*(a.Tile.Y-t.Y)))
	if dist > float64(a.MoveSpeed) {
		return throw(400, "move too far")
	}
	a.Tile = t
	return nil
}

func (a *Ant) Adjacent() []*world.Tile {
	return a.Tile.Adjacent()
}

func (a *Ant) GetType() string {
	return string(a.Type)
}

// QUEEN
func (queen *Ant) Spawn(t AntType) (*Ant, *AntError) {
	// Check if the queen is of the correct type and the spawn type is not a queen
	if queen.Type != Queen || t == Queen {
		return nil, throw(http.StatusPreconditionFailed, "Only queens can spawn, and they cannot spawn other queens")
	}

	// Check if enough time has passed since the last spawn
	currentTime := time.Now()
	if currentTime.Sub(queen.Colony.LastSpawnTime) < queen.Colony.SpawnRate {
		return nil, throw(http.StatusTooEarly, "Not enough time has passed")
	}

	// Attempt to retrieve the properties for the ant type to be spawned
	props, exists := antTypeConfig[string(t)]
	if !exists {
		throw(http.StatusNotFound, fmt.Sprintf("Ant type '%s' not found in configuration", t))
	}

	// Create the new ant
	newAnt := &Ant{
		Colony:    queen.Colony,
		Tile:      queen.Tile,
		Type:      t,
		MoveSpeed: props.MoveSpeed,
		Attack:    props.Attack,
		Defense:   props.Defense,
		HP:        props.HP,
	}

	// Update the colony's last spawn time to now
	queen.Colony.LastSpawnTime = currentTime

	// Add the new ant to the colony's list of ants
	queen.Colony.Ants = append(queen.Colony.Ants, newAnt)

	return newAnt, nil
}

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

// ANT CONFIG
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

// ERRORS
type AntError struct {
	StatusCode int    // HTTP status code
	Message    string // A human-readable message describing the error
}

// New creates a new error with a specific HTTP status code and message
func throw(statusCode int, message string) *AntError {
	return &AntError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Error implements the error interface
func (e *AntError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}
