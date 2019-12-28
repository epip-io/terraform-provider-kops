package provider

import (
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func flattenTime(d v1.Time) string {
	return d.String()
}

func flattenDuration(d *v1.Duration) string {
	return d.String()
}

func flattenQuantity(d *resource.Quantity) int64 {
	var q int64

	q = d.Value()

	return q
}

func flattenTypeMeta(in v1.TypeMeta) map[string]interface{} {
	out := make(map[string]interface{}, 1)

	if in.APIVersion != "" {
		out["api_version"] = flattenString(in.APIVersion)
	}

	if in.Kind != "" {
		out["kind"] = flattenString(in.Kind)
	}

	return out
}
