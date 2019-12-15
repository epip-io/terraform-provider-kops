{{- define "unmarshal.func" }}
{{- $_ := set . "Funcs" (append .Funcs .Type )}}
{{- if .Top }}
// Unmarshal{{ .Type | splitList "." | last | snakecase | camelcase }} exported as top level field
{{- end }}
func {{ if .Top }}U{{ else }}u{{ end }}nmarshal{{ .Type | splitList "." | last | snakecase | camelcase -}}(d {{ if .IsSlice -}}
	[]
	{{- if .IsPtr -}}
	*
	{{- end -}}
	{{- end -}}{{ .Type }}) {{ if not .Top }}[]{{ end }}map[string]interface{} {
{{- if $.IsSlice }}
	data := make([]map[string]interface{}, 0)

	for _, i := range d {
		item := make(map[string]interface{})
{{ range $n, $f := .Fields }}
	{{- if contains "HTTPProxy" $n }}
		item[{{ .Name | quote }}] = i.{{ $n }}
	{{- else if or (eq .Kind "struct") .IsPtr }}
		if i.{{ $n }} != nil {
			item[{{ .Name | quote }}] = {{ if .Fields -}}
				unmarshal{{ .Type | splitList "." | last | snakecase | camelcase }}({{- if (and $f.IsPtr (not $f.IsSlice)) }}*{{ end }}i.{{ $n }})
		{{- else -}}
			*i.{{ $n }}
		{{- end }}
		}
	{{- else }}
		item[{{ .Name | quote }}] = i.{{ $n }}
	{{- end }}
{{- end }}

		data = append(data, item)
	}

	return data
{{- else }}
	data := make(map[string]interface{})
	{{ range $n, $f := .Fields }}
		{{- if contains "ListMeta" $f.Type }}
		{{- else if eq "HTTPProxy" $n }}
	data[{{ .Name | quote }}] = unmarshalHttpProxy(d.{{ $n }})
		{{- else if (or (eq .Kind "struct") .IsPtr) }}
	if d.{{ $n }}{{ if (contains "CreationTimestamp" $n) }}.String() != ""{{ else }} != nil{{ end }} {
		data[{{ .Name | quote }}] = {{ if (eq .Kind "struct") -}}
			unmarshal{{ .Type | splitList "." | last | snakecase | camelcase }}({{- if (and $f.IsPtr (not $f.IsSlice)) -}}*{{ end }}d.{{ $n }})
		{{- else -}}
		*d.{{ $n }}
		{{- end }}
	}
		{{- else }}
	data[{{ .Name | quote }}] = d.{{ $n }}
		{{- end }}
	{{- end }}

	return {{ if .Top }}data{{ else }}[]map[string]interface{}{data}{{ end }}
{{- end }}
}
{{- range $_, $f := .Fields }}
	{{- if (and (and (eq .Kind "struct") (not .Seen)) (not (contains "ListMeta" $f.Type))) }}
		{{- if (and (not (contains "Duration" $f.Type)) (not (contains "Quantity" $f.Type))) }}
			{{- if (and (not (contains "Time" $f.Type)) (not (eq $.Type $f.Type))) }}
				{{- $_ := set $f "Funcs" $.Funcs }}
{{ template "unmarshal.func" $f }}
				{{- $_ := set $ "Funcs" $f.Funcs }}
			{{- end }}
		{{- end }}
	{{- end }}
{{- end }}
{{- end }}

package api

import (
	{{- if eq .Type "ObjectMeta"}}
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	{{- else }}
	"k8s.io/kops/pkg/apis/kops"
	{{- end }}
)
{{- $_ := set $.Schema "Funcs" list }}
{{ template "unmarshal.func" .Schema }}
