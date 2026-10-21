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
{{- $appName := include "simc-chart.app.fullname" $ -}}
{{- $hash := include "simc-chart.app.envData" $ | sha256sum | trunc 8 -}}
{{- printf "%s-env-%s" $appName $hash }}
{{- end }}

{{- define "simc-chart.app.envData" }}
{{- $_ := required "simc-chart.app.envData: expected app context" .app -}}
{{- $env := dict -}}
{{- if .app.env }}
{{- $env = merge $env .app.env -}}
{{- end -}}
{{- with .app.maxRuntimeSeconds -}}
{{- $_ := set $env "MAX_RUNTIME" (printf "%ds" (int .)) -}}
{{- end -}}
{{- $_ := set $env "TERMINATION_PERIOD" (printf "%ss" (include "simc-chart.app.terminationGraceSeconds" .)) }}
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