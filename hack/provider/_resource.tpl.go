{{- define "resource.schemas" }}
	{{- range $_, $s := . }}
		{{- template "resource.schema" $s }}
	{{- end }}
{{- end }}

{{- define "resource.schema" }}
func schema{{ .Name }}() *schema.Schema {
	return &schema.Schema{
		Type: schema.Type{{ .Type }},
	{{- if (and (not (hasSuffix "Slice" .Name)) (eq .Type "List")) }}
		MaxItems: 1,
	{{- end }}
	{{- if .Required }}
		Required: true,
	{{- else }}
		Optional: true,
	{{- end }}
	{{- if .Elems }}
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
		{{- range $_, $f := .Elems }}
			{{- if $f.Schema }}
				{{ splitList "," $f.Tag | first | snakecase | quote }}: schema{{ $f.Schema.Name }}(),
			{{- end }}
		{{- end }}
			},
		},
	{{- else if .SubType }}
		Elem: &schema.Schema{
			Type : schema.Type{{ .SubType }},
		},
	{{- end }}
	}
}
{{ end }}
