
{{- define "simc-chart.app.publisherEnvName" }}
{{- $_ := required "simc-chart.app.publisherEnvName: expected app context" .app -}}
{{- printf "%s-pub-env" .name }}
{{- end }}

{{- define "simc-chart.app.publisherEnvPrefix" }}
{{- $_ := required "simc-chart.app.publisherEnvPrefix: expected publisher key" . -}}
{{- printf "%s_" (upper . | replace "-" "_") }}
{{- end }}

{{- define "simc-chart.app.publisherData" }}
{{- $_ := required "simc-chart.app.publisherData: expected app context" .app -}}
{{- range $key, $value := .app.publisher }}
{{- $_ := required "simc-chart.app.publisherData: expected stream name" $value.streamName -}}
{{- $_ := required "simc-chart.app.publisherData: expected subject" $value.subject -}}
{{- $prefix := include "simc-chart.app.publisherEnvPrefix" $key -}}
{{ printf "%sPUBLISHER_STREAM_NAME: %s" $prefix $value.streamName }}
{{ printf "%sPUBLISHER_SUBJECT: %s" $prefix $value.subject }}
{{- end }}
{{- end }}