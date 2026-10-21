{{- define "simc-chart.app.consumer" }}
{{- $_ := required "simc-chart.app.consumer: expected app context" .app -}}
{{- include "simc-chart.app.streamExists" (dict "streams" .root.Values.streams "streamName" .app.consumer.streamName "consumerName" .app.consumer.durableName) }}
---
apiVersion: jetstream.nats.io/v1beta2
kind: Consumer
metadata:
  name: {{ .app.consumer.durableName }}
  labels:
    {{- include "simc-chart.app.labels" $ | nindent 4}}
spec:
  streamName: {{ .app.consumer.streamName }}
  durableName: {{ .app.consumer.durableName }}
  deliverPolicy: {{ default "all" .app.consumer.deliverPolicy }}
  {{- with .app.consumer.filterSubject }}
  filterSubject: {{ . }}
  {{- end }}
  {{- with .app.consumer.filterSubjects }}
  filterSubjects:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- with .app.consumer.maxDeliver }}
  maxDeliver: {{ . }}
  {{- end }}
  ackPolicy: {{ default "explicit" .app.consumer.ackPolicy }}
  ackWait: {{ default "30s" .app.consumer.ackWait }}
{{- end -}}

{{- define "simc-chart.app.consumerData" }}
{{- $_ := required "simc-chart.app.consumerData: expected app context" .app -}}
CONSUMER_STREAM_NAME: {{ .app.consumer.streamName }}
CONSUMER_NAME: {{ .app.consumer.durableName }}
{{- end }}