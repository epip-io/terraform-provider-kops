package provider

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/epip-io/terraform-provider-kops/pkg/api"
)
{{ range .Resources }}
func datasource{{ . }}() *schema.Resource {
	return &schema.Resource{
		Read: 	resource{{ . }}Read,
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