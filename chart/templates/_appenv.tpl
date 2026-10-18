{{- define "simc-chart.app.env" }}
{{- $_ := required "simc-chart.app.env: expected app context" .app -}}
{{- $data := include "simc-chart.app.envData" . }}
{{- if $data }}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "simc-chart.app.envName" $ }}
  labels:
    {{- include "simc-chart.app.labels" $ | nindent 4}}
data:
  {{- $data | nindent 2 }}
{{- end }}
{{- end -}}

{{- define "simc-chart.app.envName" -}}
{{- $_ := required "simc-chart.app.envName: expected app context" .app -}}
{{- include "simc-chart.app.fullname" $ }}-env
{{- end }}

{{- define "simc-chart.app.envData" }}
{{- $_ := required "simc-chart.app.envData: expected app context" .app -}}
{{- $env := dict -}}
{{- if .app.env }}
{{- $env = merge $env .app.env -}}
{{- end -}}
{{- if .app.consumer -}}
{{- $env = merge $env (include "simc-chart.app.consumerData" . | fromYaml) -}}
{{- end -}}
{{- if .app.publisher -}}
{{- $env = merge $env (include "simc-chart.app.publisherData" . | fromYaml) -}}
{{- end -}}
{{- with $env -}}
{{- toYaml . }}
{{- end -}}
{{- end }}