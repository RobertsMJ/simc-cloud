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
