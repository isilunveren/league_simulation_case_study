# League Simulation Case Study

A football league simulation REST API built with Go. Simulates a 4-team league with match scheduling, results, standings, and championship predictions.

## Tech Stack

- **Go** — core language
- **PostgreSQL** — database
- **net/http** — HTTP server (no framework)

## Project Structure

```
cmd/            → entry point
internal/
  domain/       → data models and interfaces
  simulation/   → match simulation and prediction logic
  repository/   → database operations
  handler/      → HTTP handlers
  db/           → database connection
migrations/     → SQL schema
```

---

## Quick Start (Live Demo)

The API is deployed and ready to use. No setup required.

**Base URL:** `https://league-simulation-case-study.onrender.com`

Test it directly in Postman:

```
GET https://league-simulation-case-study.onrender.com/league/table
```

---

## Local Setup

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)

### Steps

**1. Clone the repo**

```bash
git clone https://github.com/isilunveren/league_simulation_case_study.git
cd league_simulation_case_study
```

**2. Create database and user**

```bash
psql postgres
```

```sql
CREATE DATABASE league_simulation;
CREATE USER league_user WITH PASSWORD 'league_pass';
GRANT ALL PRIVILEGES ON DATABASE league_simulation TO league_user;
\q
```

**3. Run migrations**

```bash
psql -U league_user -d league_simulation -f migrations/schema.sql
```

**4. Configure environment**

```bash
cp .env.example .env
```

`.env` file:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=league_user
DB_PASSWORD=league_pass
DB_NAME=league_simulation
```

**5. Install dependencies**

```bash
go mod tidy
```

**6. Run**

```bash
go run cmd/main.go
```

Server runs on `http://localhost:8080`

---

## API Endpoints

### GET /league/table

Returns current standings sorted by points, goal difference, goals scored.

**Response:**

```json
[
  {
    "ID": 1,
    "Name": "Chelsea",
    "Strength": 85,
    "PlayedMatches": 3,
    "Won": 2,
    "Drawn": 0,
    "Lost": 1,
    "GoalsFor": 6,
    "GoalsAgainst": 4,
    "GoalDifference": 2,
    "Points": 6
  }
]
```

---

### POST /league/next-week

Simulates the next week's matches. Match results are determined by team strength scores.

**Response:**

```json
[
  {
    "ID": 1,
    "Week": 1,
    "HomeTeam": { "..." },
    "AwayTeam": { "..." },
    "HomeGoals": 2,
    "AwayGoals": 1,
    "IsPlayed": true
  }
]
```

---

### POST /league/play-all

Simulates all remaining weeks at once. Returns results grouped by week.

**Response:**

```json
{
  "4": ["..."],
  "5": ["..."],
  "6": ["..."]
}
```

---

### GET /league/matches?week={n}

Returns matches for a specific week. Omit `week` param to get all matches.

**Example:** `GET /league/matches?week=3`

---

### GET /league/predictions

Returns championship probability for each team. Available after week 3.
Uses Monte Carlo simulation (1000 iterations) over remaining matches.

**Response:**

```json
[
  {
    "Team": { "..." },
    "ChampionshipPct": 55.0
  }
]
```

---

### PUT /matches/{id}

Edit a match result. Automatically recalculates all team statistics.

**Request Body:**

```json
{
  "home_goals": 3,
  "away_goals": 1
}
```

---

## League Rules

- 4 teams: Chelsea (85), Arsenal (80), Manchester City (90), Liverpool (75)
- Round-robin format: each team plays every other team home and away
- 6 weeks total, 2 matches per week
- Points: Win = 3, Draw = 1, Loss = 0
- Tiebreaker: goal difference → goals scored

## Simulation Logic

Match scores are generated based on team strength:

- Higher strength = higher scoring ceiling
- Home advantage: +10 to home team strength
- Results are probabilistic — upsets are possible

## Championship Prediction

Monte Carlo method: remaining matches are simulated 1000 times.
Championship probability = (times team won / 1000) \* 100
