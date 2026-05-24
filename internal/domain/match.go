package domain

// Match - a single league match
type Match struct {
	ID         int
	Week       int
	HomeTeam   *Team
	AwayTeam   *Team
	HomeGoals  int
	AwayGoals  int
	IsPlayed   bool
}

// MatchRepository - db operations for matches
type MatchRepository interface {
	GetAll() ([]*Match, error)
	GetByID(id int) (*Match, error)
	GetByWeek(week int) ([]*Match, error)
	Update(match *Match) error
	Create(match *Match) error
}