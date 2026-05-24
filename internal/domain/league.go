package domain

// League - all league state
type League struct {
	Teams   []*Team
	Matches []*Match
	CurrentWeek int
}

// Prediction - championship probability for a team
type Prediction struct {
	Team            *Team
	ChampionshipPct float64 // 0-100
}

// LeagueSimulator defines the simulation behavior
type LeagueSimulator interface {
	PlayWeek() ([]*Match, error)
	PlayAll() (map[int][]*Match, error) // week - matches
	GetTable() ([]*Team, error)
	GetPredictions() ([]*Prediction, error)
}