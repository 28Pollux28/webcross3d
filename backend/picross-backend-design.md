
# Webcross 3D – Backend Design Document

## 1. Overview

Webcross 3D is a 3D puzzle game where players chisel blocks from a voxel grid to reveal a hidden shape, using numerical clues. This backend supports:

- Puzzle storage and retrieval
- Validating player actions against the real solution
- Tracking player lives and progress
- (Future) User accounts and authentication
- (Future) Puzzle progress saving
- (Future) Admin puzzle submission

## 2. Technology Stack

- **Language**: Go
- **Web Framework**: Fiber
- **API Schema**: OpenAPI 3.0 (with oapi-codegen)
- **Database**: PostgreSQL (future)
- **ORM**: (planned) Ent
- **Authentication**: JWT (future)
- **Project Layout**: Standard `cmd`, `internal`, `pkg` structure

## 3. Functional Scope (MVP)

| Feature                   | Description                                  |
|---------------------------|----------------------------------------------|
| Get Puzzle Metadata       | Basic info: name, size, description, author |
| Get Puzzle Clues          | Returns clues without the actual grid       |
| Validate Voxel Action     | Backend checks if removing voxel is valid   |
| Notify Puzzle Completion  | Client signals puzzle completion            |
| Health Check              | Simple service liveness endpoint            |

## 4. API Endpoints

| Method | Path                             | Description                         |
|--------|----------------------------------|-------------------------------------|
| GET    | `/api/healthz`                   | Health check                        |
| GET    | `/api/puzzles`                   | List all public puzzles             |
| GET    | `/api/puzzles/{id}`              | Get puzzle clues and metadata       |
| POST   | `/api/puzzles/{id}/actions`      | Validate a voxel removal            |
| POST   | `/api/puzzles/{id}/complete`     | Validate completion of the puzzle   |

## 5. Domain Models

### Puzzle

```go
type Puzzle struct {
ID          string
Name        string
Author      string
SizeX       int
SizeY       int
SizeZ       int
Clues       []ClueSet
}
```

### ClueSet

Each ClueSet refers to a single line in the 3D puzzle grid (along X, Y, or Z axis), identified by fixing two coordinates. The clue describes how many cubes remain in that line and how they are grouped.

```go
type ClueSplit int

const (
  NoSplit ClueSplit = iota // blocks form one group
  Split2                   // two groups (circled)
  Split3Plus               // three or more groups (squared)
)

type ClueSet struct {
  Axis     string    // "X", "Y", or "Z"
  Coord1   int       // fixed axis 1
  Coord2   int       // fixed axis 2
  Count    *int      // number of kept cubes (nil means no clue)
  Split    ClueSplit // how cubes are grouped (0 = one group, 1 = two, 2 = three+)
}
```

## 6. Folder Structure

```
backend/
├── cmd/server/main.go
├── internal/
│   ├── api/              # Generated handlers and types
│   ├── puzzle/           # Business logic
│   ├── db/               # DB access layer (future)
│   └── server/           # Fiber setup
├── api/openapi.yaml
├── api/specs.gen.yaml
├── api/stubs.gen.yaml
├── Makefile
└── go.mod
```

## 7. Security Plan

- Backend-only validation of voxel actions
- No solution grid sent to frontend
- Track player sessions and remaining lives
- JWT-based authentication (future)
- Rate limiting (planned)

## 8. Deployment

- Go binary
- Config via environment variables
- Docker support (future)
- CI/CD pipeline (planned)

## 9. Milestones

| Milestone                          | Scope                              |
|-----------------------------------|------------------------------------|
| ✅ MVP Puzzle API                 | List, fetch puzzles                |
| ⏳ Auth system                    | Register, login, JWT               |
| ⏳ Admin API                      | Submit new puzzles                 |
| ⏳ Save puzzle progress           | Per-user, optional checkpointing   |