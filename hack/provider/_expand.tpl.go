{{- define "expand.functions" }}
	{{- $skipFuncs := list "Quantity" "Duration" "Time" -}}
	{{- range $n, $f := . }}
		{{- if not (mustHas $n $skipFuncs) }}
			{{- template "expand.function" $f }}
		{{- end }}
	{{- end }}
{{- end }}

{{- define "expand.function" }}
func expand{{ .Name }}(in interface{}) {{ .Type }} {
	{{ if or .Elems (hasSuffix "Slice" .Name) -}}
	d
	{{- else -}}
	out
	{{- end }} := in
	{{- if (or (hasSuffix "Slice" .Name) (.Elems)) }}.([]interface{}){{ end -}}
	{{- if (and (.Elems) (not (hasSuffix "Slice" .Name))) }}[0].(map[string]interface{})
	{{- else if not (or .Elems (hasSuffix "Slice" .Name)) -}}.({{ .Type }})
	{{- end -}}
	{{- if hasSuffix "Slice" .Name }}
	out := make({{ .Type }} , len(d))

	for i := 0; i < len(d); i++ {
		{{- if .Elems }}
		out[i] = {{ if (hasPrefix "*" (trimPrefix "[]" .Type)) }}&{{ end }}{{ trimPrefix "*" (trimPrefix "[]" .Type) }}{}
			{{- range $n, $e := .Elems }}
				{{- if  $e.Tag }}

		if v, ok := d[i].(map[string]interface{})[{{ splitList "," $e.Tag | first | snakecase | quote }}]; ok {
			{{ if (and $e.Function (not (eq $e.Type $e.Function.Type))) }}value :=
			{{- else }}out[i].{{ $n }} =
			{{- end }} {{ template "expand.value" $e }}
			{{- if (and $e.Function (not (eq $e.Type $e.Function.Type))) }}

			out[i].{{ $n }} = {{ if eq (trimSuffix $e.Function.Type $e.Type) "*" }}&{{ else }}(*{{ end }}value{{ if not (eq (trimSuffix $e.Function.Type $e.Type) "*") }}){{ end }}
			{{- end }}
		}
				{{- else }}

		{{ if (and $e.Function (not (eq $e.Type $e.Function.Type))) }}value :=
		{{- else }}out[i].{{ $n }} =
		{{- end }} {{ template "expand.value" $e }}
		{{- if (and $e.Function (not (eq $e.Type $e.Function.Type))) }}

		out[i].{{ $n }} = {{ if eq (trimSuffix $e.Function.Type $e.Type) "*" }}&{{ else }}(*{{ end }}value{{ if not (eq (trimSuffix $e.Function.Type $e.Type) "*") }}){{ end }}
		{{- end }}
				{{- end }}
			{{- end }}
		{{- else }}
		out[i] = d[i].({{ trimPrefix "[]" .Type }})
		{{- end }}
	}
	{{- else if .Elems }}
	out := {{ if hasPrefix "*" .Type }}&{{ end }}{{ trimPrefix "*" .Type }}{}
		{{- range $n, $e := .Elems }}
			{{- if contains "omitempty" $e.Tag }}

	if v, ok := d[{{ splitList "," $e.Tag | first | snakecase | quote }}]; ok {
		{{ if (and $e.Function (not (eq $e.Type $e.Function.Type))) }}value :=
		{{- else }}out.{{ $n }} =
		{{- end }} {{ template "expand.value" $e }}
		{{- if (and $e.Function (not (eq $e.Type $e.Function.Type))) }}

		out.{{ $n }} = {{ if eq (trimSuffix $e.Function.Type $e.Type) "*" }}&{{ else }}(*{{ end }}value{{ if not (eq (trimSuffix $e.Function.Type $e.Type) "*") }}){{ end }}
		{{- end }}
	}
	{{- else }}

	{
		{{ if (and $e.Function (not (eq $e.Type $e.Function.Type))) }}value :=
		{{- else }}out.{{ $n }} =
		{{- end }} {{ if eq .Name "TypeMeta" }}expandTypeMeta(in){{ else }}{{ template "expand.value" $e }}{{ end }}
		{{- if (and $e.Function (not (eq $e.Type $e.Function.Type))) }}

		out.{{ $n }} = {{ if eq (trimSuffix $e.Function.Type $e.Type) "*" }}&{{ else }}(*{{ end }}value{{ if not (eq (trimSuffix $e.Function.Type $e.Type) "*") }}){{ end }}
		{{- end }}
	}
			{{- end }}
		{{- end }}
	{{- end }}

	return out
}
{{ end }}

{{- define "expand.value" -}}
	{{- if .Function -}}
		expand{{ .Function.Name }}(
			{{- if eq .Function.Name "FieldsMap" }}d
			{{- else if (or (contains "omitempty" .Tag) (and (not .Schema) (not (contains "inline" .Tag)))) }}v
			{{- else if (not (contains "inline" .Tag)) }}d[{{ splitList "," .Tag | first | snakecase | quote }}]
			{{- end }})
	{{- end -}}
{{- end -}}
