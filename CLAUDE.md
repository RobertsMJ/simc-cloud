# simc-cloud

Serverless AWS Lambda application that runs [SimulationCraft](https://github.com/simulationcraft/simc) simulations in the cloud.

## Architecture

(TODO) A frontend interface for building gear combinations to simulate
    - POSTs a request to the gearset generator with constraints
(TODO) Gearset generator
    - Microservice that, given a list of gear options for each gear slot, generates gearset combinations that satisfy the given constraints
        - Constraints could be similar to:
            - Minimum 4-piece set bonus
            - Minimum 2-piece set bonus
            - (built-in) Limit 2 embellishments
    - Puts the resulting valid gearsets as valid simulationcraft input strings to a queue to be simulated
Simulation Runner
    - Handles a single simulation request
    - Runs the simulation with the given configuration
    - Reports the result
(TODO) Results persistence and reporting back to the user

## Project Structure

```
backend/                    # Monolithic microservices pattern to share models and re-usable code across microservices, written in Go
  cmd/<service-name>/       # Each microservice has its own entry point here, sharing models and transport packages from the same Go module
profiles/                   # Request templates and sample simc profiles for local testing
template.yaml               # AWS SAM stack definition
Taskfile.yml                # Local dev tasks
```

## Key Patterns

### Transport handlers

Each transport package exposes a generic `NewRequestHandler[Req, Resp]` that adapts the transport-specific event format to the core handler signature:

```go
func NewRequestHandler[Req any, Resp any](
    callback func(ctx context.Context, req Req) (Resp, error),
) func(context.Context, TransportRequest) TransportResponse
```

When adding a new transport, follow this pattern — the core business logic should never depend on transport concerns.

### Models

`SimulationRequest` and `SimulationResponse` are the canonical types shared across all transports:

```go
type SimulationRequest struct {
    RequestID string          `json:"request_id"`
    GearsetID string          `json:"gearset_id"`
    Metadata  *map[string]any `json:"metadata,omitempty"`
    Input     string          `json:"input"`  // raw simc profile text
}
```

The `Input` field is the full simc profile as a plain text string (newline-separated key=value pairs).

### Simulation service

`sim/service.go` parses the `Input` string into simc CLI args (one per line, comments and empty lines filtered), then executes `/app/simc` with `json2=stdout` to capture structured JSON output directly on stdout.

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
