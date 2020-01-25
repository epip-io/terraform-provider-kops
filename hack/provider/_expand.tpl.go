{{- define "expand.functions" }}
	{{- $skipFuncs := list "Quantity" "Duration" "Time" "Int64" "Float32" "Bool" "Int32" -}}
	{{- range $n, $f := . }}
		{{- if not (mustHas $n $skipFuncs) }}
			{{- template "expand.function" $f }}
		{{- end }}
	{{- end }}
{{- end }}

{{- define "expand.function" }}
func expand{{ .Name }}(in interface{}) ({{ trimPrefix "*" .Type }}, bool) {
	{{ if hasSuffix "Map" .Name -}}
	out := make({{ .Type }})
  d := in.(map[string]interface{})
  
  if len(d) < 1 {
    return out, true
  }

	for k, v := range d {
		{{- $value := trimPrefix "map[string]" .Type -}}
		{{- if hasPrefix "[]" $value }}
			l := v.([]interface{})
			out[k] = make({{ $value }}, len(l))
			{{- $value := trimPrefix "[]" $value }}
			for i := range l {
				 out[k][i] = l[i].({{ $value }})
			}
		{{- else }}
		out[k] = v.({{ $value }})
		{{- end }}
	}
	{{- else -}}
	{{- if not (and (kindIs "map" .Elems) (not .Elems)) }}d := in
	{{- if (or (hasSuffix "Slice" .Name) (.Elems)) }}.([]interface{})
	{{- else }}.({{ template "expand.type" .Type }})
  {{- end }}
  {{- end }}
	{{- if hasSuffix "Slice" .Name }}
	out := make({{ .Type }} , len(d))

  if len(d) < 1 {
    return out, true
  }

	for i := 0; i < len(d); i++ {
		{{- if .Elems }}
		out[i] = {{ if (hasPrefix "*" (trimPrefix "[]" .Type)) }}&{{ end }}{{ trimPrefix "*" (trimPrefix "[]" .Type) }}{}
			{{- range $n, $e := .Elems }}
				{{- if  $e.Tag }}

		if v, ok := d[i].(map[string]interface{})[{{ splitList "," $e.Tag | first | snakecase | quote }}]; ok {
      if value, e := {{ template "expand.value" $e }}; !e {
        out[i].{{ $n }} = {{ if hasPrefix "*" .Type }}&{{ end }}value
      }
		}
				{{- else }}

    if value, e := {{ template "expand.value" $e }}; !e {
      out[i].{{ $n }} = {{ if hasPrefix "*" .Type }}&{{ end }}value
    }
				{{- end }}
			{{- end }}
		{{- else }}
		out[i] = d[i].({{ trimPrefix "[]" .Type }})
		{{- end }}
	}
	{{- else if .Elems }}
	out := {{ trimPrefix "*" .Type }}{}

	if len(d) < 1 {
		return out, true
	}

	m := d[0].(map[string]interface{})

		{{- range $n, $e := .Elems }}
			{{- if contains "omitempty" $e.Tag }}

	if v, ok := m[{{ splitList "," $e.Tag | first | snakecase | quote }}]; ok {
		if value, e := {{ template "expand.value" $e }}; !e {
      out.{{ $n }} = {{ if hasPrefix "*" .Type }}&{{ end }}value
    }
	}
	{{- else }}

	{
		if value, e := {{ template "expand.value" $e }}; !e {
      out.{{ $n }} = {{ if hasPrefix "*" .Type }}&{{ end }}value
    }
	}
			{{- end }}
    {{- end }}

  if isEmpty(out) {
    return out, true
  }
  {{- else if kindIs "map" .Elems -}}
	out := {{ trimPrefix "*" .Type }}{}
	
	if isEmpty(out) {
		return out, true
	}
  {{- else }}
	r := {{ trimPrefix "*" .Type }}(d)
  out := r
  
  if out == {{ if or (contains "int" .Type) (contains "float" .Type) -}}0
            {{- else -}}"" 
            {{- end }} {
    return out, true
  }
	{{- end }}
	{{- end }}

	return out, false
}
{{ end }}

{{- define "expand.value" -}}
	{{- if .Function -}}
		expand{{ .Function.Name }}(
			{{- if eq .Function.Name "FieldsMap" }}m
			{{- else if (or (contains "omitempty" .Tag) (and (not .Schema) (not (contains "inline" .Tag)))) }}v
      {{- else if (not (contains "inline" .Tag)) }}m[{{ splitList "," .Tag | first | snakecase | quote }}]
      {{- else }}in
			{{- end }})
	{{- end -}}
{{- end -}}

{{- define "expand.type" -}}
  {{- $t := trimPrefix "*" . -}}
  {{- if or (contains "int" $t) (hasPrefix "float" $t) -}}
      {{- $t -}}
  {{- else -}}
      string
  {{- end -}}
{{- end -}}