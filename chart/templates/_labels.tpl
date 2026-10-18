
{{/*
Common labels
{{ include "simc-chart.chart.labels" $ }}
*/}}
{{- define "simc-chart.chart.labels" -}}
helm.sh/chart: {{ include "simc-chart.chart.chartLabel" . }}
app.kubernetes.io/part-of: {{ include "simc-chart.chart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
App common labels
{{ include "simc-chart.app.version" $ctx }}
*/}}
{{- define "simc-chart.app.version" -}}
{{- $_ := required "simc-chart.app.labels: expected app context" .app -}}
{{ .app.image.tag | replace "+" "_" | trunc 63 | trimSuffix "-" | quote }}
{{- end -}}

{{/*
App common labels
{{ include "simc-chart.app.labels" $ctx }}
*/}}
{{- define "simc-chart.app.labels" -}}
{{- $_ := required "simc-chart.app.labels: expected app context" .root }}
{{- include "simc-chart.chart.labels" .root }}
app.kubernetes.io/name: {{ .name }}
app.kubernetes.io/version: {{ include "simc-chart.app.version" . }}
{{- if .app.component }}
app.kubernetes.io/component: {{ .app.component | quote }}
{{- end }}
{{- end -}}

{{/*
App selector labels
{{ include "simc-chart.app.selectorLabels" $ctx }}
*/}}
{{- define "simc-chart.app.selectorLabels" -}}
{{- $_ := required "simc-chart.app.selectorLabels: expected app context" .root -}}
app.kubernetes.io/name: {{ .name }}
app.kubernetes.io/instance: {{ .root.Release.Name }}
{{- end -}}

{{/*
App pod labels
{{ include "simc-chart.app.podLabels" $ctx }}
*/}}
{{- define "simc-chart.app.podLabels" -}}
{{- $_ := required "simc-chart.app.podLabels: expected app context (dict root/name/app)" .root -}}
app.kubernetes.io/part-of: {{ include "simc-chart.chart.name" .root }}
{{ include "simc-chart.app.selectorLabels" . }}
app.kubernetes.io/version: {{ include "simc-chart.app.version" . }}
{{- if .app.component }}
app.kubernetes.io/component: {{ .app.component | quote }}
{{- end }}
{{- end -}}
