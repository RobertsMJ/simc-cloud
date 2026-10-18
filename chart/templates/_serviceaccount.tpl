{{- define "simc-chart.app.serviceaccount" }}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "simc-chart.app.serviceAccountName" $ }}
  labels:
    {{- include "simc-chart.app.labels" $ | nindent 4 -}}
  {{- if and .app.serviceAccount .app.serviceAccount.annotations }}
  annotations:
    {{- toYaml .app.serviceAccount.annotations | nindent 4 }}
  {{- end }}
automountServiceAccountToken: false
{{- end -}}