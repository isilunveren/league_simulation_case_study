package repository

import (
	"database/sql"

	"github.com/isilunveren/league_simulation_case_study/internal/domain"
)

// MatchRepo implements domain.MatchRepository
type MatchRepo struct {
	DB *sql.DB
}

func NewMatchRepo(db *sql.DB) *MatchRepo {
	return &MatchRepo{DB: db}
}

func (r *MatchRepo) GetAll() ([]*domain.Match, error) {
	rows, err := r.DB.Query(`
		SELECT m.id, m.week, m.home_goals, m.away_goals, m.is_played,
			ht.id, ht.name, ht.strength, ht.played_matches, ht.won, ht.drawn, ht.lost, ht.goals_for, ht.goals_against, ht.goal_difference, ht.points,
			at.id, at.name, at.strength, at.played_matches, at.won, at.drawn, at.lost, at.goals_for, at.goals_against, at.goal_difference, at.points
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMatches(rows)
}

func (r *MatchRepo) GetByWeek(week int) ([]*domain.Match, error) {
	rows, err := r.DB.Query(`
		SELECT m.id, m.week, m.home_goals, m.away_goals, m.is_played,
			ht.id, ht.name, ht.strength, ht.played_matches, ht.won, ht.drawn, ht.lost, ht.goals_for, ht.goals_against, ht.goal_difference, ht.points,
			at.id, at.name, at.strength, at.played_matches, at.won, at.drawn, at.lost, at.goals_for, at.goals_against, at.goal_difference, at.points
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		WHERE m.week = $1
	`, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMatches(rows)
}

func (r *MatchRepo) GetByID(id int) (*domain.Match, error) {
	row := r.DB.QueryRow(`
		SELECT m.id, m.week, m.home_goals, m.away_goals, m.is_played,
			ht.id, ht.name, ht.strength, ht.played_matches, ht.won, ht.drawn, ht.lost, ht.goals_for, ht.goals_against, ht.goal_difference, ht.points,
			at.id, at.name, at.strength, at.played_matches, at.won, at.drawn, at.lost, at.goals_for, at.goals_against, at.goal_difference, at.points
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		WHERE m.id = $1
	`, id)

	m := &domain.Match{HomeTeam: &domain.Team{}, AwayTeam: &domain.Team{}}
	err := row.Scan(
		&m.ID, &m.Week, &m.HomeGoals, &m.AwayGoals, &m.IsPlayed,
		&m.HomeTeam.ID, &m.HomeTeam.Name, &m.HomeTeam.Strength,
		&m.HomeTeam.PlayedMatches, &m.HomeTeam.Won, &m.HomeTeam.Drawn, &m.HomeTeam.Lost,
		&m.HomeTeam.GoalsFor, &m.HomeTeam.GoalsAgainst, &m.HomeTeam.GoalDifference, &m.HomeTeam.Points,
		&m.AwayTeam.ID, &m.AwayTeam.Name, &m.AwayTeam.Strength,
		&m.AwayTeam.PlayedMatches, &m.AwayTeam.Won, &m.AwayTeam.Drawn, &m.AwayTeam.Lost,
		&m.AwayTeam.GoalsFor, &m.AwayTeam.GoalsAgainst, &m.AwayTeam.GoalDifference, &m.AwayTeam.Points,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MatchRepo) Update(match *domain.Match) error {
	_, err := r.DB.Exec(`
		UPDATE matches SET home_goals = $1, away_goals = $2, is_played = $3
		WHERE id = $4`,
		match.HomeGoals, match.AwayGoals, match.IsPlayed, match.ID,
	)
	return err
}

func (r *MatchRepo) Create(match *domain.Match) error {
	return r.DB.QueryRow(`
		INSERT INTO matches (week, home_team_id, away_team_id) VALUES ($1, $2, $3) RETURNING id`,
		match.Week, match.HomeTeam.ID, match.AwayTeam.ID,
	).Scan(&match.ID)
}

// scanMatches is a helper to avoid duplicate scan logic
func scanMatches(rows *sql.Rows) ([]*domain.Match, error) {
	var matches []*domain.Match
	for rows.Next() {
		m := &domain.Match{HomeTeam: &domain.Team{}, AwayTeam: &domain.Team{}}
		err := rows.Scan(
			&m.ID, &m.Week, &m.HomeGoals, &m.AwayGoals, &m.IsPlayed,
			&m.HomeTeam.ID, &m.HomeTeam.Name, &m.HomeTeam.Strength,
			&m.HomeTeam.PlayedMatches, &m.HomeTeam.Won, &m.HomeTeam.Drawn, &m.HomeTeam.Lost,
			&m.HomeTeam.GoalsFor, &m.HomeTeam.GoalsAgainst, &m.HomeTeam.GoalDifference, &m.HomeTeam.Points,
			&m.AwayTeam.ID, &m.AwayTeam.Name, &m.AwayTeam.Strength,
			&m.AwayTeam.PlayedMatches, &m.AwayTeam.Won, &m.AwayTeam.Drawn, &m.AwayTeam.Lost,
			&m.AwayTeam.GoalsFor, &m.AwayTeam.GoalsAgainst, &m.AwayTeam.GoalDifference, &m.AwayTeam.Points,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, nil
}