{{/*
Scylla connection settings as env data, for apps with a `scylla` block.
{{ include "simc-chart.scylla.connectionData" $ }}
*/}}
{{- define "simc-chart.scylla.connectionData" -}}
SCYLLA_HOSTS: {{ required "scylla.hosts is required when an app uses scylla" .Values.scylla.hosts | quote }}
SCYLLA_KEYSPACE: {{ required "scylla.keyspace is required when an app uses scylla" .Values.scylla.keyspace | quote }}
{{- end }}

{{/*
Scylla login env vars for an app, read from its `scylla.existingSecret`: a Secret with `username` and `password` keys.
{{ include "simc-chart.scylla.credentialsEnv" $ctx }}
*/}}
{{- define "simc-chart.scylla.credentialsEnv" -}}
{{- $scylla := required (printf "%s: scylla is required" .name) .app.scylla -}}
{{- $secret := required (printf "%s: scylla.existingSecret is required" .name) $scylla.existingSecret -}}
- name: SCYLLA_USERNAME
  valueFrom:
    secretKeyRef:
      name: {{ $secret }}
      key: username
- name: SCYLLA_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ $secret }}
      key: password
{{- end }}
