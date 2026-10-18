
{{- define "simc-chart.app.publisherEnvName" }}
{{- $_ := required "simc-chart.app.publisherEnvName: expected app context" .app -}}
{{- printf "%s-pub-env" .name }}
{{- end -}}

{{- define "simc-chart.app.publisherData" }}
{{- $_ := required "simc-chart.app.publisherData: expected app context" .app -}}
PUBLISHER_STREAM_NAME: {{ .app.publisher.streamName }}
PUBLISHER_SUBJECT: {{ .app.publisher.subject }}
{{- end }}