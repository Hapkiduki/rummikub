# Rummikub (Go)

A small terminal implementation of the Rummikub tile game written in Go. This is a hobby project intended to model core game mechanics (pieces, sets, runs, turns) and simple bot strategies.

## Project structure

- `main.go` - program entry point; initializes the game, creates players, and runs the main game loop.
- `player.go` - player-related logic: input handling for the human player, dealing, strategies for bots, and helper functions to manipulate hands.
- `types.go` - core types and constructors: `Pieza` (tile), `Jugador` (player), `Estrategia` interface and helper constructors (`crearMazo`, `crearJugadores`).
- `rules.go` - game rules and validation logic: checking valid sets (trios/quartets) and runs (escaleras), and helpers for adding tiles or scoring.
- `rules_test.go` - unit tests for rules (trio/run validation).

## Requirements

- Go 1.25+ (or newer). Ensure `go` is installed and on your PATH.

## How to run

From the repository root:

```bash
cd /Users/youruser/rummikub
go run .
```

This starts the interactive, terminal-based game. It will prompt for the number of players and then let a human player play against simple bot opponents.

## Tests

Run unit tests for rules with:

```bash
go test ./...
```

## Known issues & TODOs

- Some UI/UX improvements needed: ordering tiles on the table when adding, and ability for human players to pick tiles from table melds.
- Bot logic is intentionally simple; bots often prefer robbing if they cannot find a set.
- Handling of jokers (comodines) is basic; scoring and replacement logic can be improved.
- Some helper functions lack robust input validation (edge cases may cause panics if input is malformed).

## Suggested next steps

- Implement visualization improvements: sort tiles on the table when players add tiles.
- Add a CLI flag for deterministic seeding for reproducible game runs.
- Add more comprehensive tests for edge cases and bot behaviors.

## License

This project is hobby code — no specific license declared. Feel free to use and adapt it for learning and experimentation.
