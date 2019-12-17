{{- define "expand.func" }}
func expand{{ .Type | splitList "." | last | snakecase | camelcase -}}
	(d []interface{}) {{ if .IsSlice -}}
	[]
		{{- if .IsPtr -}}
		*
		{{- end -}}
	{{- end -}}{{ .Type }} {
{{- if .IsSlice }}
	var data {{ if .IsSlice -}}
	[]
	{{- end -}}
	{{- if .IsPtr -}}
	*
	{{- end -}}{{ $.Type }}

	for _, i := range d {
		item := {{ .Type }}{}
		c := i.(map[string]interface{})
{{ range $n, $f := .Fields }}
		{{- if or (eq $f.Kind "struct") $f.IsPtr }}
		if top, ok := c[{{ $f.Name | quote }}]; ok {
		{{- if eq $f.Kind "struct" }}
			v := expand{{ $f.Type | splitList "." | last | snakecase | camelcase -}}
				(top.([]interface{}))
		{{- else }}
			v := top.({{ $f.Type }})

		{{- end }}

			item.{{ $n }} = {{ if (and $f.IsPtr (not $f.IsSlice)) }}&{{ end }}v
		}
		{{- else if $f.IsSlice }}
		{{- if eq $f.Kind "struct" }}
		item.{{ $n }} = expand{{ $f.Type | splitList "." | last | snakecase | camelcase -}}
			(c[{{ $f.Name | quote }}].([]interface{}))
		{{- else if not (eq $f.Kind $f.Type) }}
		item.{{ $n }} = c[{{ $f.Name | quote }}].([]{{ $f.Type }})
		{{- else }}
		item.{{ $n }} = strings.Split(c[{{ $f.Name | quote }}].({{ $f.Type }}), ",")
		{{- end }}
		{{- else }}
		item.{{ $n }} = c[{{ $f.Name | quote }}].({{ $f.Type }})
		{{- end }}
{{- end }}

		data = append(data, {{ if .IsPtr -}}&{{- end -}}item)
	}
{{- else }}
	data := {{ .Type }}{}

	if len(d) > 0 {
		c := d[0].(map[string]interface{})
{{ range $n, $f := .Fields }}
	{{- if or (eq $f.Kind "struct") $f.IsPtr }}
	if {{ if (and (and (kindIs "map" .Fields) (not .Fields)) (not (contains "timestamp" $f.Name))) }}_{{ else }}top{{ end }}, ok := c[{{ .Name | quote }}]; ok {
	{{- if eq $f.Kind "struct" }}
		{{- if (and (and (kindIs "map" .Fields) (not .Fields)) (not (contains "timestamp" $f.Name))) }}
		v := {{ $f.Type }}{}

		{{- else }}
		v := expand{{ $f.Type | splitList "." | last | snakecase | camelcase -}}
			(top.([]interface{}))

		{{- end }}

		data.{{ $n }} = {{ if (and $f.IsPtr (not $f.IsSlice)) }}&{{ end }}v
	{{- else }}
		v := top.({{ if $f.IsSlice }}[]{{ end }}{{ $f.Type }})

		data.{{ $n }} = {{ if (and $f.IsPtr (eq $f.First "ptr")) }}&{{ end }}v
	{{- end }}
	}
	{{- else if $f.IsSlice }}
	{{- if eq $f.Kind "struct" }}
	data.{{ $n }} = expand{{ $f.Type | splitList "." | last | snakecase | camelcase -}}
		(c[{{ $f.Name | quote }}].([]interface{}))
	{{- else if not (eq $f.Kind $f.Type) }}
	data.{{ $n }} = c[{{ $f.Name | quote }}].([]{{ $f.Type }})
	{{- else if eq .Kind "uint8" }}
	data.{{ $n }} = []byte{c[{{ $f.Name | quote }}].({{ $f.Type }})}
	{{- else }}
	data.{{ $n }} = strings.Split(c[{{ $f.Name | quote }}].({{ $f.Type }}), ",")
	{{- end }}
	{{- else }}
	data.{{ $n }} = c[{{ $f.Name | quote }}].({{ $f.Type }})
	{{- end }}
{{- end }}
	}
{{- end }}

	return data
}
{{- range $_, $f := .Fields }}
	{{- if $f.Fields }}
{{ template "expand.func" $f }}
	{{- end }}
{{- end }}
{{- end }}

package convert

import (
{{- if (not (contains "SSH" .Type)) }}
	"strings"
{{- end }}
{{- if eq .Type "ObjectMeta"}}
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
{{ else }}
	"k8s.io/kops/pkg/apis/kops"
{{- end }}
)
func Expand{{ .Schema.Type | splitList "." | last | snakecase | camelcase -}} (d map[string]interface{}) {{ .Schema.Type }} {
	data := {{ .Schema.Type }}{}
{{ range $n, $f := .Schema.Fields }}
	{{- if mustHas $n (list "ManagedFields") }}
	{{- else if or (eq $f.Kind "struct") $f.IsPtr }}
	if {{ if (and (and (kindIs "map" .Fields) (not .Fields)) (not (contains "timestamp" $f.Name))) }}_{{ else }}top{{ end }}, ok := d[{{ .Name | quote }}]; ok {
		{{- if eq $f.Kind "struct" }}
			{{- if (and (and (kindIs "map" .Fields) (not .Fields)) (not (contains "timestamp" $f.Name))) }}
		v := {{ $f.Type }}{}

			{{- else }}
		v := expand{{ $f.Type | splitList "." | last | snakecase | camelcase -}}
			(top.([]interface{}))

			{{- end }}

		data.{{ $n }} = {{ if (and $f.IsPtr (not $f.IsSlice)) }}&{{ end }}v
		{{- else }}
		v := top.({{ if $f.IsSlice }}[]{{ end }}{{ $f.Type }})

		data.{{ $n }} = {{ if (and $f.IsPtr (eq $f.First "ptr")) }}&{{ end }}v
		{{- end }}
	}
	{{- else if $f.IsSlice }}
		{{- if eq $f.Kind "struct" }}
	data.{{ $n }} = expand{{ $f.Type | splitList "." | last | snakecase | camelcase -}}
		(d[{{ $f.Name | quote }}].([]interface{}))
		{{- else if not (eq $f.Kind $f.Type) }}
	data.{{ $n }} = d[{{ $f.Name | quote }}].([]{{ $f.Type }})
		{{- else if eq .Kind "uint8" }}
	data.{{ $n }} = []byte{d[{{ $f.Name | quote }}].({{ $f.Type }})}
		{{- else }}
	data.{{ $n }} = strings.Split(d[{{ $f.Name | quote }}].({{ $f.Type }}), ",")
		{{- end }}
	{{- else }}
	data.{{ $n }} = d[{{ $f.Name | quote }}].({{ $f.Type }})
	{{- end }}
{{- end }}

	return data
}
{{- range $n, $f := .Schema.Fields }}
	{{- if (and $f.Fields (not (mustHas $n (list "ManagedFields")))) }}
{{ template "expand.func" $f }}
	{{- end }}
{{- end }}
