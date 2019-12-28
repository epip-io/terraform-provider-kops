{{- define "flatten.functions" }}
	{{- $skipFuncs := list "Quantity" "Duration" "Time" "TypeMeta" -}}
	{{- range $n, $f := . }}
		{{- if (not (mustHas $n $skipFuncs)) }}
			{{- template "flatten.function" $f }}
		{{- end }}
	{{- end }}
{{- end }}

{{- define "flatten.function" }}
func flatten{{ .Name }}(in {{ template "flatten.type.arg" . }}) {{ template "flatten.type.return" . }} {
	out := {{ template "flatten.type.out" . }}
	{{- if hasPrefix "map" .Type }}

	for k, v := range in {
		out[k] = v
	}
	{{- else if (or .Elems (hasPrefix "[]" $.Type)) }}

	for i := 0; i < len(out); i++ {
		{{- if .Elems }}
		out[i] = make(map[string]interface{})
			{{- range $n, $e := .Elems }}
				{{- if contains "inline" $e.Tag }}
	
		{
			m := {{ template "flatten.elem.func" (dict "FuncName" $e.Function.Name "Name" $n "Type" $e.Type "ParentType" $.Type "Function" $e.Function) }}

			for k, v := range m {
				out[i][k] = v
			}
		}
				{{- else if contains "omitempty" $e.Tag }}

		if 	{{ if mustHas .Type (list "v1.Time" "*v1.Time") }}{{ else if (and (and (hasPrefix "*" $e.Type) (hasPrefix "*" $e.Function.Type)) (eq .Type "*bool")) }}*
		{{- else }}{{ template "flatten.ptr.if" (dict "Type" $e.Type "FuncType" $e.Function.Type) }}
		{{- end }}in{{ if hasPrefix "[]" $.Type }}[i]{{ end }}.{{ $n }}{{ if mustHas $e.Type (list "v1.Time" "*v1.Time") }}.String(){{ end }}
			{{- template "flatten.var.if" (dict "FuncName" $e.Function.Name "Type" $e.Type "Name" $n) }} {
			out[i][{{ template "flatten.var.out" $e.Tag }}] = {{ template "flatten.elem.func" (dict "FuncName" $e.Function.Name "Name" $n "Type" $e.Type "ParentType" $.Type "Function" $e.Function) }}
		}
				{{- else }}

		out[i][{{ template "flatten.var.out" $e.Tag }}] = {{ template "flatten.elem.func" (dict "FuncName" $e.Function.Name "Name" $n "Type" $e.Type "ParentType" $.Type "Function" $e.Function) }}
				{{- end }}
			{{- end }}
		{{- else }}
		out[i] = in[i]
		{{- end }}
	}
	{{- end }}

	return {{ template "flatten.var.return" . }}
}
{{ end }}

{{- define "flatten.elem.func" -}}
flatten{{ .FuncName }}(
	{{- if .Function }}
		{{- if (and (hasPrefix "*" .Function.Type) (not (hasPrefix "*" .Type))) -}}
			&
		{{- else -}}
			{{ template "flatten.ptr.if" (dict "Type" .Type "FuncType" .Function.Type) }}
		{{- end -}}
	{{- end -}}
	in{{ if hasPrefix "[]" .ParentType }}[i]{{ end }}.{{ .Name }})
{{- end -}}

{{- define "flatten.type.arg" -}}
{{ .Type }}
{{- end -}}

{{- define "flatten.type.out" -}}
	{{- if (or .Elems (or (hasPrefix "map" .Type) (hasPrefix "[]" $.Type)))  -}}
		make(
			{{- template "flatten.type.return" . }}
			{{- if hasPrefix "[]" $.Type }}, len(in)
			{{- else if (not (hasPrefix "map" .Type)) }}, 1
			{{- end -}}
		)
	{{- else -}}
		{{ if hasPrefix "*" .Type }}*{{ end }}in
	{{- end }}
{{- end -}}

{{- define "flatten.type.return" -}}
	{{- if (or (hasPrefix "map" .Type) (eq .Name "TypeMeta")) -}}
		map[string]interface{}
	{{- else if .Elems -}}
		[]map[string]interface{}
	{{- else if hasPrefix "[]" $.Type -}}
		[]interface{}
	{{- else -}}
		interface{}
	{{- end -}}
{{- end -}}

{{- define "flatten.var.if" -}}
	{{- $stringFunc := list "kops.InstanceGroupRole" "types.UID" "v1.StatusReason" "v1.Time" "*v1.Time" -}}
	{{- if (not (eq .FuncName "Bool")) }} != 
		{{- if mustHas .Type (list "v1.ListMeta") }} (v1.ListMeta{})
		{{- else if mustHas .Type (list "kops.HTTPProxy") }} (kops.HTTPProxy{})
		{{- else if (or (or (eq .FuncName "String") (mustHas .Type $stringFunc)) (contains "Type" .Type)) }} ""
		{{- else if (contains "int" .Type) }} 0
		{{- else }} nil
		{{- end }}
	{{- end }}
{{- end -}}

{{- define "flatten.var.out" -}}
{{ splitList "," . | first | snakecase | quote }}
{{- end -}}

{{- define "flatten.var.return" -}}
out
{{- end -}}

{{- define "flatten.ptr.if" }}
	{{- if (and (hasPrefix "*" .Type) (not (hasPrefix "*" .FuncType))) -}}
		*
	{{- end -}}
{{- end -}}
