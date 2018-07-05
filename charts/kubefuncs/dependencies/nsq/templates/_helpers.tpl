{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "admin.name" -}}
{{ template "name" . }}-admin
{{- end -}}

{{- define "admin.fullname" -}}
{{ template "fullname" . }}-admin
{{- end -}}

{{- define "nsqd.name" -}}
{{ template "name" . }}-nsq
{{- end -}}

{{- define "nsqd.fullname" -}}
{{ template "fullname" . }}-nsq
{{- end -}}

{{- define "lookupd.name" -}}
{{ template "name" . }}-lookupd
{{- end -}}

{{- define "lookupd.fullname" -}}
{{ template "fullname" . }}-lookupd
{{- end -}}

{{/*
DNS lookupd address.
*/}}
{{- define "dns.lookupd" -}}
{{ template "lookupd.fullname" . }}.{{ .Release.Namespace }}.svc.cluster.local
{{- end -}}
