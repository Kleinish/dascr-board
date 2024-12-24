package player

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Kleinish/dascr-board/logger"
	"github.com/go-chi/chi"
)

type PlayerStats struct {
	GamesPlayed   int     `json:"gamesPlayed"`
	GamesWon      int     `json:"gamesWon"`
	GamesLost     int     `json:"gamesLost"`
	WinPercentage float64 `json:"winPercentage"`
	HighestScore  int     `json:"highestScore"`
	BestFinish    int     `json:"bestFinish"`
	CheckoutRate  float64 `json:"checkoutRate"`
	ThreeDartAvg  float64 `json:"threeDartAvg"`
	FirstNinePPR  float64 `json:"firstNinePPR"`
	OneEighties   int     `json:"oneEighties"`
	OneForties    int     `json:"oneForties"`
	TonPlus       int     `json:"tonPlus"`
}

func UpdatePlayerStats(db *sql.DB, playerID string, stats PlayerStats) error {
	upsert := `
		INSERT OR REPLACE INTO player_stats (
			player_id, games_played, games_won, games_lost,
			highest_score, best_finish, checkout_rate,
			three_dart_avg, first_nine_ppr,
			one_eighties, one_forties, ton_plus
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := db.Prepare(upsert)
	if err != nil {
		return fmt.Errorf("preparing update stats statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		playerID,
		stats.GamesPlayed,
		stats.GamesWon,
		stats.GamesLost,
		stats.HighestScore,
		stats.BestFinish,
		stats.CheckoutRate,
		stats.ThreeDartAvg,
		stats.FirstNinePPR,
		stats.OneEighties,
		stats.OneForties,
		stats.TonPlus,
	)
	return err
}

func GetPlayerStats(db *sql.DB, playerID string) (PlayerStats, error) {
	var stats PlayerStats

	query := `
		SELECT COALESCE(games_played, 0),
			   COALESCE(games_won, 0),
			   COALESCE(games_lost, 0),
			   COALESCE(highest_score, 0),
			   COALESCE(best_finish, 0),
			   COALESCE(checkout_rate, 0.0),
			   COALESCE(three_dart_avg, 0.0),
			   COALESCE(first_nine_ppr, 0.0),
			   COALESCE(one_eighties, 0),
			   COALESCE(one_forties, 0),
			   COALESCE(ton_plus, 0)
		FROM player_stats
		WHERE player_id = ?
	`

	row := db.QueryRow(query, playerID)
	err := row.Scan(
		&stats.GamesPlayed,
		&stats.GamesWon,
		&stats.GamesLost,
		&stats.HighestScore,
		&stats.BestFinish,
		&stats.CheckoutRate,
		&stats.ThreeDartAvg,
		&stats.FirstNinePPR,
		&stats.OneEighties,
		&stats.OneForties,
		&stats.TonPlus,
	)

	if err != nil && err != sql.ErrNoRows {
		return stats, fmt.Errorf("fetching player stats: %w", err)
	}

	if stats.GamesPlayed > 0 {
		stats.WinPercentage = float64(stats.GamesWon) / float64(stats.GamesPlayed) * 100
	}

	return stats, nil
}

func GetPlayerStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := chi.URLParam(r, "id")

		stats, err := GetPlayerStats(db, playerID)
		if err != nil {
			logger.Errorf("Getting player stats: %+v", err)
			http.Error(w, "Error fetching stats", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			logger.Errorf("Error writing response: %+v", err)
		}
	}
}
