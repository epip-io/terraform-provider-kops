package provider

{{- if .Functions }}
import (
	"k8s.io/kops/pkg/apis/kops"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

{{ template "flatten.functions" .Functions }}
{{- else }}
	{{- range $n, $r := . }}
import (
	"k8s.io/kops/pkg/apis/kops"
{{- if eq $n "Cluster" }}

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
{{- end }}
)

{{ template "flatten.functions" $r.Functions }}
	{{- end }}
{{- end }}