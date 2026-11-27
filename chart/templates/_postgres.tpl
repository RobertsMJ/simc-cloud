{{/*
Postgres connection settings as env data, for apps with a `db` block.
{{ include "simc-chart.postgres.connectionData" $ }}
*/}}
{{- define "simc-chart.postgres.connectionData" -}}
DATABASE_URL: {{ required "db.existingSecret is required" .app.db.existingSecret }}
{{- end }}
