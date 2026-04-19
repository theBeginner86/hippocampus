{{/*
Expand the name of the chart.
*/}}
{{- define "hippocampus.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "hippocampus.fullname" -}}
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
*/}}
{{- define "hippocampus.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "hippocampus.labels" -}}
helm.sh/chart: {{ include "hippocampus.chart" . }}
{{ include "hippocampus.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "hippocampus.selectorLabels" -}}
app.kubernetes.io/name: {{ include "hippocampus.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "hippocampus.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "hippocampus.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Generate cluster peers list for RAFT
Format: "node-0:service-0.headless:7000:6379,node-1:service-1.headless:7000:6379"
*/}}
{{- define "hippocampus.clusterPeers" -}}
{{- $replicas := .Values.cluster.replicas | int }}
{{- $name := include "hippocampus.fullname" . }}
{{- $headless := .Values.headlessService.name }}
{{- $rpcPort := .Values.raft.rpcPort }}
{{- $dataPort := .Values.service.port }}
{{- range $i := until $replicas }}
{{- if gt $i 0 }},{{- end }}
{{- printf "%s-%d:%s-%d.%s:%s:%s" $name $i $name $i $headless $rpcPort $dataPort }}
{{- end }}
{{- end }}

