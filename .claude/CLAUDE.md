# simc-cloud

Cloud application that runs [SimulationCraft](https://github.com/simulationcraft/simc) simulations: first, comparing a World of Warcraft character's gear and talent combinations; long term, group simulations along a dungeon route (route building à la Mythic Dungeon Tools). It's a learning project for system design and distributed systems, so aim suggestions at the target architecture.

## Target architecture

**This section describes the desired end state, not the current code.** Before assuming any piece exists, cross-reference `backend/cmd/*`, `backend/go.mod`, and `.cluster/addons/*.yaml` against `TODO.md`.

Go services on Kubernetes (local: k3d) in an Istio ambient-mode mesh, using NATS JetStream for async work and ScyllaDB for persistence. Each service uses [fx](https://github.com/uber-go/fx) and composes `backend/platform` for its admin server (`/health`, `/ready`, `/metrics`).

Pipeline:

1. Frontend: user picks gear options and constraints, POSTs to create-job.
2. `create-job`: validates the request, creates the job, hands the item set to the gearset generator, and returns the job ID for polling.
3. Gearset generator: expands per-slot gear options into combinations that satisfy constraints (e.g. minimum 2- or 4-piece set bonus, at most 2 crafted embellishments, an upgrade-currency budget), then publishes each one as a simc input string to a JetStream stream.
4. `run-sim`: consumes gearsets, runs simc, and publishes results to a results stream.
5. `result-writer`: consumes results and writes result records plus job completion counters to the database. It's a separate service from `run-sim` so a DB write failure can't lose a finished simulation. Unprocessable messages go to a dead-letter stream.
6. Frontend: navigates to `/{job-id}` page, listens to SSE for job status, and displays final results when processing completes.

Planned deployment tooling: Helm chart for the app, ArgoCD + GitOps, Tilt for the local dev loop.

## Commands

- Validation tasks may be obtained from the repo root by running `task -l`
- Local cluster: `task install-deps` once, then `task cluster:create` / `task cluster:delete`.
- `task run-simc -- profiles/sample.simc` runs upstream simc in Docker, which is useful for checking a simc input string by hand.
