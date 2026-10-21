{{- define "simc-chart.app.podmonitor" }}
{{- $_ := required "simc-chart.app.podmonitor: expected app context" .root }}
---
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: {{ include "simc-chart.app.fullname" . }}-monitor
spec:
  selector:
    matchLabels:
      {{- include "simc-chart.app.selectorLabels" $ | nindent 6 }}
  podMetricsEndpoints:
    - port: admin
      path: /metrics
{{- end }}