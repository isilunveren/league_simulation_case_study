package simulation

import (
	"errors"
	"math/rand"
	"sort"
	"github.com/isilunveren/league_simulation_case_study/internal/domain"
)

type Simulator struct {
	League      *domain.League
	TeamRepo    domain.TeamRepository
	MatchRepo   domain.MatchRepository
}

func NewSimulator(league *domain.League, teamRepo domain.TeamRepository, matchRepo domain.MatchRepository) *Simulator {
	return &Simulator{
		League:    league,
		TeamRepo:  teamRepo,
		MatchRepo: matchRepo,
	}
}

// core match logic

// simulateScore generates a score based on team strengths
func simulateScore(home, away *domain.Team) (int, int) {
	homeAdvantage := 10
	homeGoals := rand.Intn((home.Strength + homeAdvantage) / 20)
	awayGoals := rand.Intn(away.Strength / 20)
	return homeGoals, awayGoals
}

// updateTeamStats updates team stats after a match
func updateTeamStats(team *domain.Team, goalsFor, goalsAgainst int) {
	team.PlayedMatches++
	team.GoalsFor += goalsFor
	team.GoalsAgainst += goalsAgainst
	team.GoalDifference = team.GoalsFor - team.GoalsAgainst

	if goalsFor > goalsAgainst {
		team.Won++
		team.Points += 3
	} else if goalsFor == goalsAgainst {
		team.Drawn++
		team.Points++
	} else {
		team.Lost++
	}
}

// LeagueSimulator interface implementation

func (s *Simulator) PlayWeek() ([]*domain.Match, error) {
	var weekMatches []*domain.Match

	for _, match := range s.League.Matches {
		if match.Week == s.League.CurrentWeek && !match.IsPlayed {
			homeGoals, awayGoals := simulateScore(match.HomeTeam, match.AwayTeam)
			match.HomeGoals = homeGoals
			match.AwayGoals = awayGoals
			match.IsPlayed = true

			updateTeamStats(match.HomeTeam, homeGoals, awayGoals)
			updateTeamStats(match.AwayTeam, awayGoals, homeGoals)

			if err := s.MatchRepo.Update(match); err != nil {
				return nil, err
			}
			if err := s.TeamRepo.Update(match.HomeTeam); err != nil {
				return nil, err
			}
			if err := s.TeamRepo.Update(match.AwayTeam); err != nil {
				return nil, err
			}

			weekMatches = append(weekMatches, match)
		}
	}

	if len(weekMatches) == 0 {
		return nil, errors.New("no matches to play this week")
	}

	s.League.CurrentWeek++
	return weekMatches, nil
}

func (s *Simulator) PlayAll() (map[int][]*domain.Match, error) {
	results := make(map[int][]*domain.Match)

	for _, match := range s.League.Matches {
		if !match.IsPlayed {
			weekMatches, err := s.PlayWeek()
			if err != nil {
				return nil, err
			}
			results[s.League.CurrentWeek-1] = weekMatches
		}
	}

	return results, nil
}

func (s *Simulator) GetTable() ([]*domain.Team, error) {
	teams, err := s.TeamRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// sort: points -> goal difference -> goals for
	sort.Slice(teams, func(i, j int) bool {
		if teams[i].Points != teams[j].Points {
			return teams[i].Points > teams[j].Points
		}
		if teams[i].GoalDifference != teams[j].GoalDifference {
			return teams[i].GoalDifference > teams[j].GoalDifference
		}
		return teams[i].GoalsFor > teams[j].GoalsFor
	})

	return teams, nil
}

func (s *Simulator) GetPredictions() ([]*domain.Prediction, error) {
	if s.League.CurrentWeek < 4 {
		return nil, errors.New("predictions available after week 4")
	}

	const iterations = 1000
	championCounts := make(map[int]int) // teamID -> win count

	for i := 0; i < iterations; i++ {
		// copy current team states
		teamCopies := make(map[int]*domain.Team)
		for _, t := range s.League.Teams {
			copy := *t
			teamCopies[t.ID] = &copy
		}

		// simulate remaining matches
		for _, match := range s.League.Matches {
			if !match.IsPlayed {
				home := teamCopies[match.HomeTeam.ID]
				away := teamCopies[match.AwayTeam.ID]
				hg, ag := simulateScore(home, away)
				updateTeamStats(home, hg, ag)
				updateTeamStats(away, ag, hg)
			}
		}

		// find champion (highest points)
		var champion *domain.Team
		for _, t := range teamCopies {
			if champion == nil || t.Points > champion.Points ||
				(t.Points == champion.Points && t.GoalDifference > champion.GoalDifference) {
				champion = t
			}
		}
		if champion != nil {
			championCounts[champion.ID]++
		}
	}

	var predictions []*domain.Prediction
	for _, team := range s.League.Teams {
		predictions = append(predictions, &domain.Prediction{
			Team:            team,
			ChampionshipPct: float64(championCounts[team.ID]) / float64(iterations) * 100,
		})
	}

	// sort by probability
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].ChampionshipPct > predictions[j].ChampionshipPct
	})

	return predictions, nil
}

// RecalculateStats resets and recalculates all team stats from match results
func (s *Simulator) RecalculateStats() error {
	// reset all teams
	for _, team := range s.League.Teams {
		team.PlayedMatches = 0
		team.Won = 0
		team.Drawn = 0
		team.Lost = 0
		team.GoalsFor = 0
		team.GoalsAgainst = 0
		team.GoalDifference = 0
		team.Points = 0
	}

	// rebuild from played matches
	teamMap := make(map[int]*domain.Team)
	for _, t := range s.League.Teams {
		teamMap[t.ID] = t
	}

	for _, match := range s.League.Matches {
		if !match.IsPlayed {
			continue
		}
		home := teamMap[match.HomeTeam.ID]
		away := teamMap[match.AwayTeam.ID]
		updateTeamStats(home, match.HomeGoals, match.AwayGoals)
		updateTeamStats(away, match.AwayGoals, match.HomeGoals)
	}

	// persist to db
	for _, team := range s.League.Teams {
		if err := s.TeamRepo.Update(team); err != nil {
			return err
		}
	}
	return nil
}