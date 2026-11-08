{{- define "simc-chart.app.scaledobject" }}
{{- $_ := required "simc-chart.app.scaledobject: autoscaling requires a consumer" .app.consumer }}
---
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: {{ include "simc-chart.app.fullname" $ }}
  labels:
    {{- include "simc-chart.app.labels" $ | nindent 4 }}
spec:
  scaleTargetRef:
    name: {{ include "simc-chart.app.fullname" $ }}
  minReplicaCount: {{ .app.autoscaling.minReplicas }}
  maxReplicaCount: {{ .app.autoscaling.maxReplicas }}
  triggers:
    - type: nats-jetstream
      metadata:
        natsServerMonitoringEndpoint: {{ .root.Values.autoscaling.natsMonitoringEndpoint | quote }}
        account: "$G"
        stream: {{ .app.consumer.streamName }}
        consumer: {{ .app.consumer.durableName }}
        lagThreshold: {{ .app.autoscaling.lagThreshold | quote }}
{{- end }}
