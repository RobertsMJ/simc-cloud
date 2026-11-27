CREATE TABLE jobs (
  id                   uuid PRIMARY KEY,
  status               text NOT NULL CHECK (status IN ('pending', 'done', 'error')),
  spec                 jsonb NOT NULL,
  created_at           timestamptz NOT NULL DEFAULT now(),
  completed_at         timestamptz,
  simc_version         text,
  baseline_dps         double precision,
  baseline_dps_std_err double precision,
  error                text
);

CREATE TABLE gearset_results (
  job_id      uuid NOT NULL REFERENCES jobs ON DELETE CASCADE,
  gearset_id  text NOT NULL,
  loadout     jsonb NOT NULL,
  dps         double precision NOT NULL,
  dps_std_err double precision NOT NULL,
  PRIMARY KEY (job_id, gearset_id)
);
