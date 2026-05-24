package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/isilunveren/league_simulation_case_study/internal/simulation"
	"github.com/isilunveren/league_simulation_case_study/internal/domain"
)

// LeagueHandler holds the simulator
type LeagueHandler struct {
	Simulator domain.LeagueSimulator
	MatchRepo domain.MatchRepository
}

func NewLeagueHandler(simulator domain.LeagueSimulator, matchRepo domain.MatchRepository) *LeagueHandler {
	return &LeagueHandler{Simulator: simulator, MatchRepo: matchRepo}
}

// helpers
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /league/table
func (h *LeagueHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	teams, err := h.Simulator.GetTable()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, teams)
}

// POST /league/next-week
func (h *LeagueHandler) PlayNextWeek(w http.ResponseWriter, r *http.Request) {
	matches, err := h.Simulator.PlayWeek()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, matches)
}

// POST /league/play-all
func (h *LeagueHandler) PlayAll(w http.ResponseWriter, r *http.Request) {
	results, err := h.Simulator.PlayAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

// GET /league/matches?week=3
func (h *LeagueHandler) GetMatches(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	if weekStr == "" {
		matches, err := h.MatchRepo.GetAll()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, matches)
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid week number")
		return
	}

	matches, err := h.MatchRepo.GetByWeek(week)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, matches)
}

// GET /league/predictions
func (h *LeagueHandler) GetPredictions(w http.ResponseWriter, r *http.Request) {
	predictions, err := h.Simulator.GetPredictions()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, predictions)
}

// PUT /matches/{id}
func (h *LeagueHandler) UpdateMatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid match id")
		return
	}

	match, err := h.MatchRepo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "match not found")
		return
	}

	var body struct {
		HomeGoals int `json:"home_goals"`
		AwayGoals int `json:"away_goals"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	match.HomeGoals = body.HomeGoals
	match.AwayGoals = body.AwayGoals

	if err := h.MatchRepo.Update(match); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, match)
	// recalculate team stats after match edit
	if sim, ok := h.Simulator.(*simulation.Simulator); ok {
		// update match in league memory
		for i, m := range sim.League.Matches {
			if m.ID == id {
				sim.League.Matches[i].HomeGoals = body.HomeGoals
				sim.League.Matches[i].AwayGoals = body.AwayGoals
				break
			}
		}
		sim.RecalculateStats()
	}

	writeJSON(w, http.StatusOK, match)
}