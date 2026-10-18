
{{/*
Assert that the stream with the given name exists
{{- include "simc-chart.app.streamExists" (dict "streams" .root.Values.streams "streamName" .app.consumer.streamName) }}
*/}}
{{- define "simc-chart.app.streamExists" -}}
{{- if not (hasKey .streams .streamName) -}}
{{- printf "Consumer %s references stream %s which does not exist" .app.consumer.durableName .app.consumer.streamName | fail -}}
{{- end -}}
{{- end -}}