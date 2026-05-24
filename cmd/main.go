package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/isilunveren/league_simulation_case_study/internal/db"
	"github.com/isilunveren/league_simulation_case_study/internal/domain"
	"github.com/isilunveren/league_simulation_case_study/internal/handler"
	"github.com/isilunveren/league_simulation_case_study/internal/repository"
	"github.com/isilunveren/league_simulation_case_study/internal/simulation"
)

func main() {
	// load .env
	godotenv.Load()

	// db connection
	database, err := db.Connect()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer database.Close()

	// repos
	teamRepo := repository.NewTeamRepo(database)
	matchRepo := repository.NewMatchRepo(database)

	// load teams and matches from db
	teams, err := teamRepo.GetAll()
	if err != nil {
		log.Fatal("Failed to load teams:", err)
	}

	matches, err := matchRepo.GetAll()
	if err != nil {
		log.Fatal("Failed to load matches:", err)
	}

	// build league
	currentWeek := 1
	for _, m := range matches {
		if m.IsPlayed && m.Week >= currentWeek {
			currentWeek = m.Week + 1
		}
	}

	league := &domain.League{
		Teams:       teams,
		Matches:     matches,
		CurrentWeek: currentWeek,
	}

	// if no matches exist yet, generate the schedule
	if len(matches) == 0 {
		generateSchedule(league, matchRepo)
	}

	// simulator
	sim := simulation.NewSimulator(league, teamRepo, matchRepo)

	// handlers
	h := handler.NewLeagueHandler(sim, matchRepo)

	// routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /league/table", h.GetTable)
	mux.HandleFunc("POST /league/next-week", h.PlayNextWeek)
	mux.HandleFunc("POST /league/play-all", h.PlayAll)
	mux.HandleFunc("GET /league/matches", h.GetMatches)
	mux.HandleFunc("GET /league/predictions", h.GetPredictions)
	mux.HandleFunc("PUT /matches/{id}", h.UpdateMatch)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// generateSchedule creates round-robin match schedule
func generateSchedule(league *domain.League, matchRepo domain.MatchRepository) {
	teams := league.Teams
	// 4 team round-robin: weeks 1-3 first half, weeks 4-6 return fixtures
	fixtures := [][2]int{
		{0, 3}, {1, 2}, // week 1
		{0, 2}, {3, 1}, // week 2
		{0, 1}, {2, 3}, // week 3
	}

	for week, pair := range fixtures {
		actualWeek := (week/2) + 1

		m := &domain.Match{
			Week:     actualWeek,
			HomeTeam: teams[pair[0]],
			AwayTeam: teams[pair[1]],
		}
		matchRepo.Create(m)
		league.Matches = append(league.Matches, m)

		// return fixture in second half
		m2 := &domain.Match{
			Week:     actualWeek + 3,
			HomeTeam: teams[pair[1]],
			AwayTeam: teams[pair[0]],
		}
		matchRepo.Create(m2)
		league.Matches = append(league.Matches, m2)
	}
}