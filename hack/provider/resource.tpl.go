package provider
{{- if .Schemas }}

import (
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"
)
{{ template "resource.schemas" .Schemas }}
{{- else }}
	{{- range $n, $r := . }}

import (
	"github.com/hashicorp/terraform/helper/schema"
	{{- if not (contains "SSH" $n) }}
	// "github.com/hashicorp/terraform/helper/validation"
	{{- end }}
)

func resource{{ $n }}() *schema.Resource {
	return &schema.Resource{
		Create: resource{{ $n }}Create,
		Read:   resource{{ $n }}Read,
		Update: resource{{ $n }}Update,
		Delete: resource{{ $n }}Delete,
		Exists: resource{{ $n }}Exists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
		{{- range $_, $f := $r.Elems }}
			{{ splitList "," $f.Tag | first | snakecase | quote }}: schema{{ $f.Schema.Name }}(),
		{{- end }}
		},
	}
}

func datasource{{ $n }}() *schema.Resource {
	return &schema.Resource{
		Read: resource{{ $n }}Read,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
		{{- range $_, $f := $r.Elems }}
			{{ splitList "," $f.Tag | first | snakecase | quote }}: schema{{ $f.Schema.Name }}(),
		{{- end }}
		},
	}
}
		{{- template "resource.schemas" $r.Schemas }}
	{{- end }}
{{- end }}