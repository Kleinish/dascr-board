package game

import (
	"github.com/Kleinish/dascr-board/player"
	"github.com/Kleinish/dascr-board/podium"
	"github.com/Kleinish/dascr-board/settings"
	"github.com/Kleinish/dascr-board/undo"
	"github.com/Kleinish/dascr-board/ws"
)

// BaseGame will hold all the game specific stuff
// for all games combined
type BaseGame struct {
	UID               string
	Game              string
	Player            []player.Player
	Variant           string
	In                string
	Out               string
	ActivePlayer      int
	ThrowRound        int
	GameState         string
	Message           string
	Settings          *settings.Settings
	UndoLog           *undo.Log
	Podium            *podium.Podium
	CricketController *CricketGameController
	SoundToPlay       string
}

// Game will be the interface for different games
type Game interface {
	StartGame() error
	GetStatus() BaseGame
	GetStatusDisplay() BaseGame
	NextPlayer(h *ws.Hub)
	RequestThrow(number, modifier int, h *ws.Hub) error
	Undo(h *ws.Hub) error
	Rematch(h *ws.Hub) error
}
