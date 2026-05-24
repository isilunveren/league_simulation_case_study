CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    strength INT NOT NULL,
    played_matches INT DEFAULT 0,
    won INT DEFAULT 0,
    drawn INT DEFAULT 0,
    lost INT DEFAULT 0,
    goals_for INT DEFAULT 0,
    goals_against INT DEFAULT 0,
    goal_difference INT DEFAULT 0,
    points INT DEFAULT 0
);
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    week INT NOT NULL,
    home_team_id INT REFERENCES teams(id),
    away_team_id INT REFERENCES teams(id),
    home_goals INT DEFAULT 0,
    away_goals INT DEFAULT 0,
    is_played BOOLEAN DEFAULT FALSE
);
-- initial teams
INSERT INTO teams (name, strength)
VALUES ('Chelsea', 85),
    ('Arsenal', 80),
    ('Manchester City', 90),
    ('Liverpool', 75);