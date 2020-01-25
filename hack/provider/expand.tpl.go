package provider

import (
	"k8s.io/kops/pkg/apis/kops"
{{- if (or .Functions (hasKey . "Cluster")) }}
	"k8s.io/apimachinery/pkg/apis/meta/v1"
{{- end }}
{{- if hasKey . "Cluster" }}
	"k8s.io/apimachinery/pkg/api/resource"
{{- end }}
{{- if .Functions }}
	"k8s.io/apimachinery/pkg/types"
{{- end }}
)

{{- if .Functions }}
{{ template "expand.functions" .Functions }}
{{- else }}
	{{- range $_, $r := . }}
{{ template "expand.functions" $r.Functions }}
	{{- end }}
{{- end }}