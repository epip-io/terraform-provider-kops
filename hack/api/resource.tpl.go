package provider

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/epip-io/terraform-provider-kops/pkg/api"
)
{{ range .Resources }}
func resource{{ . }}() *schema.Resource {
	return &schema.Resource{
		Create: resource{{ . }}Create,
		Read: 	resource{{ . }}Read,
		Update: resource{{ . }}Update,
		Delete: resource{{ . }}Delete,
		Exists: resource{{ . }}Exists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"metadata": api.SchemaObjectMeta(),
			"spec": 	api.Schema{{ . }}Spec(),
		},

	}
}
{{ end }}