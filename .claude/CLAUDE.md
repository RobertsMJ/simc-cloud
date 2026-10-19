# simc-cloud

Cloud-native application that runs [SimulationCraft](https://github.com/simulationcraft/simc) simulations: first, comparing a World of Warcraft character's gear and talent combinations; long term, group simulations along a dungeon route (route building à la Mythic Dungeon Tools). It's a learning project for system design and distributed systems, so aim suggestions at the target architecture.

## Target architecture

**This section describes the desired end state, not the current code.** Before assuming any piece exists, reference the existin implementation.

Go services on Kubernetes (local: k3d) in an Istio ambient-mode mesh, using NATS JetStream for async work and ScyllaDB for persistence. Each service uses [fx](https://github.com/uber-go/fx) and composes the appropriate `backend/platform` for platform implementation details. The `backend/transport` package provides a layer of abstraction above the platform implementation details.

Pipeline:

1. Frontend: user picks gear options and constraints, POSTs to create-job.
2. `create-job`: validates the request, creates the job, hands the job configuration to the gearset generator, and returns the job ID to the frontend for finding status and results.
3. Gearset generator: expands per-slot gear options into combinations that satisfy constraints (e.g. minimum 2- or 4-piece set bonus, at most 2 crafted embellishments, an upgrade-currency budget), then publishes each one as a simc input string to an event stream.
4. `run-sim`: Durable consumer that consumes gearsets, runs simc, and publishes results to a results stream. Autoscaled via keda based on consumer lag.
5. `result-writer`: Idempotent on (job_id, gearset_id). Consumes results and writes result records plus job completion counters to the database. It's a separate service from `run-sim` so a DB write failure can't lose a finished simulation. Unprocessable messages go to a dead-letter stream.
6. Frontend: navigates to `/{job-id}` page, listens to SSE for job status, and displays final results when processing completes.

Planned deployment tooling: Helm chart for the app, ArgoCD + GitOps, Tilt for the local dev loop.

## Commands

- Repo tooling is implemented in Task. Use `task -l` to obtain a list of existing tasks.
