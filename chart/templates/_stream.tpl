{{- define "simc-chart.stream" }}
{{- $_ := .stream.subjects | required "A stream must have at least one subject defined" -}}
---
apiVersion: jetstream.nats.io/v1beta2
kind: Stream
metadata:
  name: {{ include "simc-chart.stream.fullname" $ }}
  labels:
    {{- include "simc-chart.chart.labels" .root | nindent 4 }}
spec:
  name: {{ include "simc-chart.stream.fullname" $ }}
  subjects:
    {{- .stream.subjects | toYaml | nindent 4 }}
  storage: {{ .stream.storage | default "memory" }}
  maxAge: {{ .stream.maxAge | default "1h" }}
{{- end }}