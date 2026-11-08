{{- define "simc-chart.app.service" -}}
{{- $_ := required "simc-chart.app.service: expected app context" .root }}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ include "simc-chart.app.fullname" . }}
  labels:
    {{- include "simc-chart.app.labels" . | nindent 4 }}
spec:
  selector:
    {{- include "simc-chart.app.selectorLabels" . | nindent 4 }}
  ports:
    {{- range $name, $port := .app.ports }}
    - name: {{ $name }}
      port: {{ $port }}
      targetPort: {{ $name }}
      protocol: TCP
    {{- end }}
{{- end }}