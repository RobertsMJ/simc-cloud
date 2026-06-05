# simc-cloud

Serverless AWS Lambda application that runs [SimulationCraft](https://github.com/simulationcraft/simc) simulations in the cloud. The Go backend is designed to be platform-agnostic: AWS infrastructure (Lambda, SQS, DynamoDB) is isolated behind repository interfaces and transport adapters so the core business logic can be ported to other runtimes (e.g. Kubernetes + Kafka + Postgres) without modification.

## Architecture

(TODO) A frontend interface for building gear combinations to simulate

- POSTs a request to the Job Creation service with a set of gear and gear constraints

Job Creation Service

- Creates a job and invokes the Gearset Generator with the provided set of items
- Validates the request is valid
- Returns the job ID to the frontend for status monitoring and results

(TODO) Gearset generator

- Given a list of gear options for each gear slot, generates gear combinations that satisfy the given constraints
    - Constraints could be similar to:
        - Minimum 4-piece set bonus
        - Minimum 2-piece set bonus
        - Limit 2 crafted embellishments
        - Upgrade currency (i.e. user has 60 crests, how many items can be upgraded with this currency)
- Puts the resulting valid gearsets as valid simulationcraft input strings to a queue to be simulated

Simulation Runner

- Consumes gearset simulation requests from SQS
- Runs the simulation with the given configuration
- Publishes the result to a results SQS queue for persistence

(TODO) Persistence Service

- Consumes simulation results from the results SQS queue
- Writes results to DynamoDB (result record + atomic job completion counter)
- Decoupled from the Simulation Runner so that DB write failures don't lose simulation results — unprocessable messages go to a dead-letter queue
- DynamoDB table: single-table design with `PK=job_id`, `SK=RESULT#<gearset_id>`
    - Each result item stores `statistics` (promoted from simc output for leaderboard queries without deserializing the full result blob) and the full `result` blob
- Access patterns:
    - Get all sim results for a Job by Job ID (leaderboard, sharing, bookmarking)
    - Get a single result by Job ID + Gearset ID (drilldown)

Long term:

- Group simulations with a dungeon route
- Dungeon route building, a la Mythic Dungeon Tools

## Database Data Model

The database uses a single-table design keyed on `(job_id, record_type)`. All items for a job are co-located under the same `job_id` partition, enabling efficient queries without cross-table joins.

### Access Patterns

| Pattern | Key condition |
|---|---|
| Get job status / progress counters | `job_id = <id>`, `record_type = JOB` |
| Get all results for a job (leaderboard) | `job_id = <id>`, `record_type begins_with RESULT#` |
| Get a single result (drilldown) | `job_id = <id>`, `record_type = RESULT#<gearset_id>` |

### Record: Job (`record_type = JOB`)

Created by the Job Creation service when a job is submitted.

| Field | Type | Notes |
|---|---|---|
| `job_id` | string | UUID, partition key |
| `record_type` | string | `"JOB"` |
| `status` | string | `in_progress` \| `completed` \| `error` |
| `total_count` | number | Total number of gearsets to simulate |
| `completed_count` | number | Incremented atomically per successful result |
| `failed_count` | number | Incremented atomically per failed result |
| `created_at` | string | ISO 8601 timestamp |

`completed_count` and `failed_count` are updated atomically in the same transaction as each result write so the job record always reflects current progress.

### Record: Sim Result (`record_type = RESULT#<gearset_id>`)

Written by the Persistence service after a simulation completes. The `statistics` fields are promoted out of the raw simc output so leaderboard queries can sort/filter without deserializing the full result blob.

| Field | Type | Notes |
|---|---|---|
| `id` | string | Partition key |
| `record_type` | string | `"RESULT#<gearset_id>"` |
| `gearset_id` | string | |
| `status` | string | `completed` \| `error` |
| `error_message` | string? | Present only when `status=error` |
| `is_baseline` | bool | Whether this gearset is the comparison baseline |
| `statistics` | object | Promoted from simc output — see `SimStatistics` in [models/sim.go](backend/models/sim.go) |
| `metadata` | map? | Arbitrary key/value pairs forwarded from the original request |
| `result` | object? | Full raw simc JSON output (`SimcOutput = map[string]any`); omitted on error |

`statistics` sub-fields (all numeric): `elapsed_cpu_seconds`, `elapsed_time_seconds`, `init_time_seconds`, `merge_time_seconds`, `analyze_time_seconds`, `total_events_processed`, and stat samples (`sum`, `count`, `mean`, `min`, `max`, `median`, `variance`, `std_dev`, `mean_variance`, `mean_std_dev`) for `simulation_length`, `raid_dps`, and `total_dmg`.

## Project Structure

```
backend/                    # Monolithic microservices pattern to share models and re-usable code across microservices, written in Go
  cmd/<service-name>/       # Each microservice has its own entry point here, sharing models and transport packages from the same Go module
profiles/                   # Request templates and sample simc profiles for local testing
template.yaml               # AWS SAM stack definition
Taskfile.yml                # Local dev tasks
```

## Key Patterns

### Portability constraint

Repository interfaces and transport handlers must be defined in terms of domain types only (`models.*`, `context.Context`). AWS-specific types (`events.SQSEvent`, `dynamodbav` tags, etc.) must never appear in interfaces or domain packages — only in their concrete implementations.

### Transport handlers

Each transport package exposes a generic `NewRequestHandler[Req, Resp]` that adapts the transport-specific event format to the core handler signature:

```go
func NewRequestHandler[Req any, Resp any](
    callback func(ctx context.Context, req Req) (Resp, error),
) func(context.Context, TransportRequest) TransportResponse
```

When adding a new transport, follow this pattern — the core business logic should never depend on transport concerns.

### Models

`SimRequest` and `SimResult` are the canonical types shared across all transports:

```go
type SimRequest struct {
    JobID     string          `json:"job_id"`
    GearsetID string          `json:"gearset_id"`
    Metadata  *map[string]any `json:"metadata,omitempty"`
    Input     string          `json:"input"`  // raw simc profile text
}
```

The `Input` field is the full simc profile as a plain text string (newline-separated key=value pairs).

`SimResult` is the domain result type, carrying `Status`, `Statistics` (typed, promoted from `sim.statistics` in the simc JSON), and `Result` (the full raw simc output as `SimcOutput`).

### Simulation service

`sim/service.go` parses the `Input` string into simc CLI args (one per line, comments and empty lines filtered), executes `/app/simc` writing JSON output to a temp file, then unmarshals the output twice: once into a typed envelope to extract `sim.statistics`, and once into `SimcOutput` (a `map[string]any`) for the full result blob.

## Local Development

### Prerequisites

- Go, Docker, AWS CLI, AWS SAM CLI, [Task](https://taskfile.dev)

### Tasks

```bash
# Run simc directly via Docker (writes JSON to profiles/results/result.json) for verifying output formats and simulation validity
task run-simc -- profiles/sample.simc

# Invoke Lambda locally via SAM with an API Gateway event
task sim-profile-apigw -- profiles/sample.simc

# Invoke Lambda locally via SAM with an SQS event
task sim-profile-sqs -- profiles/sample.simc
```

The `apigw` and `sqs` tasks build a `SimulationRequest` JSON body from `profiles/templates/request.json` (injecting the simc file as `input`) and wrap it in the appropriate event envelope. Edit `request.json` to change `request_id`, `gearset_id`, or other defaults for local testing.

## Testing

Tests use [testify](https://github.com/stretchr/testify) suites. Test fixture files live in `test-data/` subdirectories alongside the package under test.

```bash
cd backend && go test ./...
```
