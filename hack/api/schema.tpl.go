{{- define "schema.func" }}
{{- if .Top }}
// Schema{{ .Type | splitList "." | last | snakecase | camelcase }} exported as top level field
{{- end }}
func {{ if .Top }}S{{ else }}s{{ end }}chema{{ .Type | splitList "." | last | snakecase | camelcase -}}() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
{{- if .Required }}
		Required: true,
{{- else }}
		Optional: true,
{{- end }}
{{- if not .IsSlice }}
		MaxItems: 1,
{{- end }}
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
{{- range $n, $f := .Fields }}
	{{- if eq $f.Kind "map" }}
		{{- if contains "map[string]" $f.Type }}
				{{ $f.Name | quote }}: schemaStringMap(),
		{{- end }}
	{{- else if eq $f.Type "v1.Time" }}
				{{ $f.Name | quote }}: schemaStringComputed(),
	{{- else if (or (or (contains "int" $f.Kind) (contains "Quantity" $f.Type)) (contains "Duration" $f.Type)) }}
				{{ $f.Name | quote }}: schemaIntOptional(),
	{{- else if (contains "float" $f.Kind) }}
				{{ $f.Name | quote }}: schemaFloatOptional(),
	{{- else if contains "cidr" $f.Name }}
				{{ $f.Name | quote }}: schemaCIDRString{{- if $f.Required }}Required{{ else }}Optional{{ end }}(),
	{{- else if $f.Fields }}
				{{ $f.Name | quote }}: schema{{ $f.Type  | splitList "." | last | snakecase | camelcase }}(),
	{{- else if (and .IsSlice (eq .Kind "string")) }}
				{{ $f.Name | quote }}: schemaStringSlice{{- if $f.Required }}Required{{ else }}Optional{{ end }}(),
	{{- else if eq .Kind "struct" }}
				{{ $f.Name | quote }}: schema{{ $f.Type  | splitList "." | last | snakecase | camelcase }}(),
	{{- else }}
				{{ $f.Name | quote }}: schema{{ $f.Type  | splitList "." | last | snakecase | camelcase -}}
					{{- if $f.Required }}Required{{ else }}Optional{{ end }}(),
	{{- end }}
{{- end }}
			},
		},
	}
}
{{- range $_, $f := .Fields }}
	{{- if (and (eq .Kind "struct") (not .Seen)) }}
		{{- if (and (not (contains "Duration" .Type)) (not (contains "Quantity" .Type))) }}
{{ template "schema.func" $f }}
		{{- end }}
	{{- end }}
{{- end }}
{{- end }}

package api

import (
	"github.com/hashicorp/terraform/helper/schema"
)
{{ template "schema.func" .Schema }}
