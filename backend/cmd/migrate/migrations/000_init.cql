CREATE TABLE IF NOT EXISTS jobs (
    job_id        uuid,
    gearset_id    text,          -- hash of the rendered loadout

    -- job
    status        text static,   -- queued | done | failed
    export        text static,   -- raw addon export, verbatim
    spec          text static,   -- options + constraints, JSON
    simc_version  text static,
    baseline_dps  double static,
    baseline_err  double static,
    error         text static,
    created_at    timestamp static,
    completed_at  timestamp static,

    -- per-gearset
    loadout       text,          -- JSON, written by blacksmith
    run_id        uuid static,   -- the claim: null until farseer accepts a run
    dps           double,        -- written by farseer
    dps_err       double,

    PRIMARY KEY (job_id, gearset_id)
) WITH default_time_to_live = 604800 -- 7 days
  AND compaction = {'class': 'TimeWindowCompactionStrategy',
                    'compaction_window_unit': 'DAYS',
                    'compaction_window_size': 1};
