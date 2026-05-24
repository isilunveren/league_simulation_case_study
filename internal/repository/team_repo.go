package repository

import (
	"database/sql"

	"github.com/isilunveren/league_simulation_case_study/internal/domain"
)

// TeamRepo implements domain.TeamRepository
type TeamRepo struct {
	DB *sql.DB
}

func NewTeamRepo(db *sql.DB) *TeamRepo {
	return &TeamRepo{DB: db}
}

func (r *TeamRepo) GetAll() ([]*domain.Team, error) {
	rows, err := r.DB.Query(`SELECT id, name, strength, played_matches, won, drawn, lost, goals_for, goals_against, goal_difference, points FROM teams`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*domain.Team
	for rows.Next() {
		t := &domain.Team{}
		err := rows.Scan(
			&t.ID, &t.Name, &t.Strength,
			&t.PlayedMatches, &t.Won, &t.Drawn, &t.Lost,
			&t.GoalsFor, &t.GoalsAgainst, &t.GoalDifference, &t.Points,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func (r *TeamRepo) GetByID(id int) (*domain.Team, error) {
	t := &domain.Team{}
	err := r.DB.QueryRow(`SELECT id, name, strength, played_matches, won, drawn, lost, goals_for, goals_against, goal_difference, points FROM teams WHERE id = $1`, id).Scan(
		&t.ID, &t.Name, &t.Strength,
		&t.PlayedMatches, &t.Won, &t.Drawn, &t.Lost,
		&t.GoalsFor, &t.GoalsAgainst, &t.GoalDifference, &t.Points,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TeamRepo) Update(team *domain.Team) error {
	_, err := r.DB.Exec(`
		UPDATE teams SET
			played_matches = $1, won = $2, drawn = $3, lost = $4,
			goals_for = $5, goals_against = $6, goal_difference = $7, points = $8
		WHERE id = $9`,
		team.PlayedMatches, team.Won, team.Drawn, team.Lost,
		team.GoalsFor, team.GoalsAgainst, team.GoalDifference, team.Points,
		team.ID,
	)
	return err
}

func (r *TeamRepo) Create(team *domain.Team) error {
	return r.DB.QueryRow(`
		INSERT INTO teams (name, strength) VALUES ($1, $2) RETURNING id`,
		team.Name, team.Strength,
	).Scan(&team.ID)
}