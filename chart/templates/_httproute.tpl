{{- define "simc-chart.app.httproute" -}}
{{- $_ := required (printf "%s: route requires ports" .name) .app.ports -}}
{{- $svc := include "simc-chart.app.fullname" . -}}
{{- $rules := list -}}
{{- range (.app.route.rules | default (list (dict))) }}
  {{- $rule := omit . "port" -}}
  {{- if not $rule.backendRefs }}
    {{- $portName := .port -}}
    {{- if not $portName }}
      {{- if ne (len $.app.ports) 1 }}
        {{- fail (printf "%s: app has multiple ports; set `port` on each route rule" $.name) }}
      {{- end }}
      {{- $portName = keys $.app.ports | first -}}
    {{- end }}
    {{- $port := required (printf "%s: route references unknown port %q" $.name $portName) (get $.app.ports $portName) -}}
    {{- $_ := set $rule "backendRefs" (list (dict "name" $svc "port" $port)) -}}
  {{- end }}
  {{- $rules = append $rules $rule -}}
{{- end }}
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: {{ $svc }}
  labels:
    {{- include "simc-chart.app.labels" . | nindent 4 }}
spec:
  parentRefs:
    {{- toYaml (.app.route.parentRefs | default .root.Values.gateway.parentRefs) | nindent 4 }}
  hostnames:
    {{- toYaml (.app.route.hostnames | default .root.Values.gateway.hostnames) | nindent 4 }}
  rules:
    {{- toYaml $rules | nindent 4 }}
{{- end }}
