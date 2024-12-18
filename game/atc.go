package game

import (
	"math/rand"
	"time"

	"github.com/Kleinish/dascr-board/player"
	"github.com/Kleinish/dascr-board/podium"
	"github.com/Kleinish/dascr-board/score"
	"github.com/Kleinish/dascr-board/throw"
	"github.com/Kleinish/dascr-board/undo"
	"github.com/Kleinish/dascr-board/utils"
	"github.com/Kleinish/dascr-board/ws"
)

// ATCGame will hold the ATC game information
type ATCGame struct {
	Base BaseGame
}

// StartGame will satisfy interface Game for ATC
func (g *ATCGame) StartGame() error {
	// CreateScore for each player
	// and init empty thtrow splice
	for i := range g.Base.Player {
		score := score.BaseScore{
			CurrentNumber: 1,
		}
		g.Base.Player[i].Score = score
		g.Base.Player[i].ThrowRounds = make([]throw.Round, 0)
		g.Base.Player[i].LastThrows = make([]throw.Throw, 3)
	}

	g.Base.Podium = &podium.Podium{}

	g.Base.UndoLog = &undo.Log{}
	sequence := g.Base.UndoLog.CreateSequence()
	sequence.AddActionToSequence(undo.Action{
		Action: "CREATEGAME",
	})
	g.Base.SoundToPlay = "nextplayer"

	return nil
}

// GetStatus will satisfy interface Game for ATC
func (g *ATCGame) GetStatus() BaseGame {
	return g.Base
}

// GetStatusDisplay will satisfy interface Game for game ATC
func (g *ATCGame) GetStatusDisplay() BaseGame {
	return stripDisplay(&g.Base)
}

// NextPlayer will satisfy interface Game for game ATC
func (g *ATCGame) NextPlayer(h *ws.Hub) {
	switchToNextPlayer(&g.Base, h)
}

// RequestThrow will satisfy inteface Game for ATC
func (g *ATCGame) RequestThrow(number, modifier int, h *ws.Hub) error {
	sequence := g.Base.UndoLog.CreateSequence()
	activePlayer := &g.Base.Player[g.Base.ActivePlayer]
	previousMessage := g.Base.Message
	previousState := g.Base.GameState

	// Check game state
	if g.Base.GameState == "THROW" {
		// check if ongoing round else create
		checkOngoingElseCreate(activePlayer, &g.Base, sequence)

		// Add Throw to last round
		throwRound := addThrowToCurrentRound(activePlayer, &g.Base, sequence, number, modifier)

		// Score logic
		if activePlayer.Score.CurrentNumber == number {
			switch g.Base.Variant {
			case "normal":
				increaseNumberByOne(g, throwRound, activePlayer, sequence)
				break
			case "fast":
				increaseNumberByMod(g, throwRound, activePlayer, sequence, modifier)
				break
			default:
				break
			}
		}

		if number == 0 || modifier == 0 {
			g.Base.SoundToPlay = "miss"
		}
		// Check if 3 throws in round and close round
		// Also set gameState and perhaps increase game.Base.ThrowRound
		// if everyone has already thrown to this round
		if len(throwRound.Throws) == 3 {
			closePlayerRound(&g.Base, activePlayer, throwRound, sequence, []undo.Action{}, previousState, previousMessage)
		}
	}

	// Set assets for Frontend
	setFrontendAssets(activePlayer, &g.Base)

	// Update scoreboard
	utils.WSSendUpdate(g.Base.UID, h)

	return nil
}

// Undo will satisfy interface Game for game ATC
func (g *ATCGame) Undo(h *ws.Hub) error {
	if err := triggerUndo(&g.Base, h); err != nil {
		return err
	}

	return nil
}

// Rematch will satisfy interface Game for game ATC
func (g *ATCGame) Rematch(h *ws.Hub) error {
	// Init random number generator
	s := rand.NewSource(time.Now().Unix())
	rg := rand.New(s)

	// Reset game state
	g.Base.Message = ""
	g.Base.GameState = "THROW"
	g.Base.Podium.ResetPodium()
	g.Base.UndoLog.ClearLog()
	g.Base.ActivePlayer = rg.Intn(len(g.Base.Player))
	g.Base.ThrowRound = 1
	g.Base.SoundToPlay = "nextplayer"

	for i := range g.Base.Player {
		score := score.BaseScore{
			CurrentNumber: 1,
		}
		g.Base.Player[i].Score = score
		g.Base.Player[i].ThrowRounds = make([]throw.Round, 0)
		g.Base.Player[i].LastThrows = make([]throw.Throw, 3)
	}

	sequence := g.Base.UndoLog.CreateSequence()
	sequence.AddActionToSequence(undo.Action{
		Action: "CREATEGAME",
	})

	// Update scoreboard
	utils.WSSendUpdate(g.Base.UID, h)

	return nil
}

// increases active players number by one
// will go from 20 to 25
// will go from 25 to win
func increaseNumberByOne(game *ATCGame, throwRound *throw.Round, p *player.Player, sequence *undo.Sequence) {
	previousNumberToHit := p.Score.CurrentNumber
	if p.Score.CurrentNumber == 25 {
		atcwin(game, throwRound, p, sequence)
	} else if p.Score.CurrentNumber != 20 {
		p.Score.CurrentNumber++
	} else {
		p.Score.CurrentNumber = 25
	}
	sequence.AddActionToSequence(undo.Action{
		Action:              "ATCINCREASENUMBER",
		GameID:              game.Base.UID,
		Player:              p,
		PreviousNumberToHit: previousNumberToHit,
	})
	game.Base.SoundToPlay = "1plib"
}

// increase active player number += modifier
// will go from 20 to 25
// will go from 25 to win
func increaseNumberByMod(game *ATCGame, throwRound *throw.Round, p *player.Player, sequence *undo.Sequence, modifier int) {
	previousNumberToHit := p.Score.CurrentNumber
	for i := 0; i < modifier; i++ {
		if p.Score.CurrentNumber == 25 {
			atcwin(game, throwRound, p, sequence)
		} else if p.Score.CurrentNumber != 20 {
			p.Score.CurrentNumber++
		} else {
			p.Score.CurrentNumber = 25
		}
	}
	sequence.AddActionToSequence(undo.Action{
		Action:              "ATCINCREASENUMBER",
		GameID:              game.Base.UID,
		Player:              p,
		PreviousNumberToHit: previousNumberToHit,
	})
	switch modifier {
	case 1:
		game.Base.SoundToPlay = "1plib"
	case 2:
		game.Base.SoundToPlay = "2plibs"
	case 3:
		game.Base.SoundToPlay = "3plibs"
	default:
	}
}

// atcwin is for winning the atc game
func atcwin(game *ATCGame, throwRound *throw.Round, activePlayer *player.Player, sequence *undo.Sequence) {
	previousMessage := game.Base.Message
	previousState := game.Base.GameState

	if game.Base.Settings.Podium {
		// Do podium and continue game
		doPodium(&game.Base, activePlayer, sequence)
		return
	}

	doWin(&game.Base)
	throwRound.Done = true
	sequence.AddActionToSequence(undo.Action{
		Action:            "DOWIN",
		GameID:            game.Base.UID,
		PreviousGameState: previousState,
		PreviousMessage:   previousMessage,
		Player:            activePlayer,
	})
}
