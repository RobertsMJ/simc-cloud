# AGENTS

Cloud-native application that runs [SimulationCraft](https://github.com/simulationcraft/simc) simulations to compare a World of Warcraft character's gear and talent combinations.

It's a learning project for system design and distributed systems. The service topology is the learning goal: don't propose merging services or replacing messaging with direct calls. Within that architecture, prefer the simplest design and the least code.

## Instructions

- Ground code in real API contracts: verify Go signatures with `go doc`, and SimC options against the simc version pinned in `backend/Dockerfile`, before using them.

## Pipeline

1. `api`: HTTP entry point. Parses addon exports, validates job specs against the gearset cap, creates the job, and hands it to `blacksmith`.
2. `blacksmith`: expands per-slot gear options into gearsets and renders them as profilesets in one SimC document per job.
3. `target-dummy`: runs simc on a job's document and publishes the result. Scaled by KEDA on consumer lag.
4. `farseer`: consumes results and records them against the job.

## Invariants

- One simc run per job: gearsets are profilesets in a single document, so simc shares setup across them and KEDA scales across jobs, not within one.
- Every JetStream consumer sees at-least-once delivery, so handlers must tolerate duplicates. `farseer` is the dedupe point, idempotent on `job_id`.
- Persistence stays out of `target-dummy`, so a DB write failure can't lose a finished simulation.
- Streams and consumers are declared in `chart/values.yaml`; services bind to them and never create them. The id passed to `natsx.NewPublisherModule` must match its `publisher:` key in the chart.
- Cluster infrastructure (NATS, KEDA, Flux, gateway, monitoring) is installed from the `homelab` repo, expected at `../homelab`. This repo only ships the app chart.

## Roadmap

Planned; parts may exist, so check the code.

- `blacksmith` filters combinations by constraints: minimum 2- or 4-piece set bonus, at most 2 crafted embellishments, an upgrade-currency budget.
- Fan-out: chunking a job's gearsets across multiple simc runs.
- ScyllaDB for persistence of results and job status.
- SSE event stream for live job progress.
- Dead-letter stream for messages `farseer` or `target-dummy` can't process.
- Frontend: picks gear options and constraints, then follows progress and results on a `/{job-id}` page.

## Commands

- Repo tooling is implemented in Task. Use `task -l` to list tasks.
