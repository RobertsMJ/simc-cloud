{{- define "simc-chart.app.deployment" }}
{{- $_ := required "simc-chart.app.deployment: expected app context (dict root/name/app)" .root }}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "simc-chart.app.fullname" $ }}
  labels:
    {{- include "simc-chart.app.labels" $ | nindent 4 }}
spec:
  replicas: {{ default 1 .app.replicas }}
  selector:
    matchLabels:
      {{- include "simc-chart.app.selectorLabels" $ | nindent 6 }}
  template:
    metadata:
      {{- with include "simc-chart.app.podAnnotations" $ }}
      annotations:
        {{- . | nindent 8 }}
      {{- end }}
      labels:
        {{- include "simc-chart.app.podLabels" $ | nindent 8 }}
        {{- with .app.podLabels }}
        {{- toYaml . | nindent 8 }}
        {{- end }}
    spec:
      serviceAccountName: {{ include "simc-chart.app.serviceAccountName" $ }}
      {{- with .app.podSecurityContext }}
      securityContext:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      terminationGracePeriodSeconds: {{ include "simc-chart.app.terminationGraceSeconds" $ }}
      automountServiceAccountToken: false
      containers:
        - name: {{ .name }}
          {{- with .app.securityContext }}
          securityContext:
            {{- toYaml . | nindent 12 }}
          {{- end }}
          image: "{{.app.image.registry}}/{{ .app.image.repository }}:{{ .app.image.tag }}"
          {{- with .app.image.pullPolicy }}
          imagePullPolicy: {{ . }}
          {{- end }}
          {{- with include "simc-chart.app.envData" $ }}
          envFrom:
            - configMapRef:
                name: {{ include "simc-chart.app.envName" $ }}
          {{- end }}
          ports:
            - name: admin
              containerPort: {{ .app.adminPort }}
              protocol: TCP
            {{- with .app.service }}
            - name: http
              containerPort: {{ .port }}
              protocol: TCP
            {{- end }}
          readinessProbe:
            httpGet:
              path: /ready
              port: admin
            {{- with .app.probes.readiness }}
            {{- toYaml . | nindent 12 }}
            {{- end }}
          livenessProbe:
            httpGet:
              path: /health
              port: admin
            {{- with .app.probes.liveness }}
            {{- toYaml . | nindent 12 }}
            {{- end }}
          startupProbe:
            httpGet:
              path: /health
              port: admin
            {{- with .app.probes.startup }}
            {{- toYaml . | nindent 12 }}
            {{- end }}
          {{- with .app.resources }}
          resources:
            {{- toYaml . | nindent 12 }}
          {{- end }}
          {{- with .app.volumeMounts }}
          volumeMounts:
            {{- toYaml . | nindent 12 }}
          {{- end }}
      {{- with .app.volumes }}
      volumes:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .app.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .app.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .app.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
{{- end -}}