{{/*
Expand the name of the chart.
{{ include "simc-chart.chart.name" $ }}
*/}}
{{- define "simc-chart.chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
{{ include "simc-chart.chart.fullname" $ }}
*/}}
{{- define "simc-chart.chart.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
{{ include "simc-chart.chart.chartLabel" $ }}
*/}}
{{- define "simc-chart.chart.chartLabel" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create app fullname.
{{ include "simc-chart.app.fullname" $ctx }}
*/}}
{{- define "simc-chart.app.fullname" -}}
{{- $_ := required "simc-chart.app.fullname: expected app context" .name }}
{{- .name  | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create stream fullname
{{ include "simc-chart.stream.fullname" $ctx }}
*/}}
{{- define "simc-chart.stream.fullname" -}}
{{- $_ := required "simc-chart.stream.fullname: expected stream context" .stream }}
{{- .name  | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create app service account name
{{ include "simc-chart.app.serviceAccountName" $ctx }}
*/}}
{{- define "simc-chart.app.serviceAccountName" -}}
{{- $_ := required "simc-chart.app.serviceAccountName: expected app context" .app }}
{{- if and .app.serviceAccount .app.serviceAccount.existing }}
{{- .app.serviceAccount.existing | quote }}
{{- else }}
{{- include "simc-chart.app.fullname" . }}
{{- end }}
{{- end }}

{{/*
App pod annotations: config checksums merged with podAnnotations.
{{ include "simc-chart.app.podAnnotations" $ctx }}
*/}}
{{- define "simc-chart.app.podAnnotations" -}}
{{- $_ := required "simc-chart.app.podAnnotations: expected app context" .app -}}
{{- $annotations := dict -}}
{{- with .app.podAnnotations }}
{{- $annotations = merge $annotations . -}}
{{- end }}
{{- with $annotations -}}
{{- toYaml . -}}
{{- end -}}
{{- end }}

{{/*
App termination grace period
{{ include "simc-chart.app.terminationGraceSeconds" $ctx }}
*/}}
{{- define "simc-chart.app.terminationGraceSeconds" -}}
{{- $tail := int (default 15 (default .root.Values.default.shutdownTailSeconds .app.shutdownTailSeconds)) -}}
{{- $grace := int (default 30 .app.terminationGracePeriodSeconds) -}}
{{- if .app.consumer -}}
  {{-  $_ := required "consumer must define maxRuntimeSeconds" .app.maxRuntimeSeconds -}}
  {{- $min := add (int .app.maxRuntimeSeconds) $tail -}}
  {{- if .app.terminationGracePeriodSeconds -}}
    {{- if lt $grace $min -}}
      {{- fail (printf "%s: terminationGracePeriodSeconds (%d) < maxRuntimeSeconds + shutdownTailSeconds (%d)" .name $grace $min) -}}
    {{- end -}}
  {{- else -}}
    {{- $grace = $min -}}
  {{- end -}}
{{- end -}}
{{- $grace -}}
{{- end -}}
