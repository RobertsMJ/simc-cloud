#!/bin/bash
for te in 1.0 0.5 0.2 0.1 0.05 0.01; do for i in 1 2 3 4 5 6 7 8 9 10; do echo "$te $i"; done; done | shuf |
while read te i; do
  jq -n --rawfile p profiles/sample.simc --arg te "$te" --arg i "$i" '{
    job_id: "load-\($te)", gearset_id: "g\($i)",
    metadata: {target_error: ($te|tonumber)},
    input: ($p + "\ntarget_error=\($te)\n")
  }' | nats -s localhost:4222 pub simc.job
done
