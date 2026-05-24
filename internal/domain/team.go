package domain
// team represents a football team with its statistics and performance metrics
type Team struct {
	ID             int
	Name           string
	Strength       int 
	PlayedMatches  int 
	Won            int 
	Drawn          int 
	Lost           int 
	GoalsFor       int 
	GoalsAgainst   int 
	GoalDifference int 
	Points         int 
}

// TeamRepository defines the interface for data access operations 
type TeamRepository interface {
	GetAll() ([]*Team, error)
	GetByID(id int) (*Team, error)
	Update(team *Team) error
	Create(team *Team) error
}