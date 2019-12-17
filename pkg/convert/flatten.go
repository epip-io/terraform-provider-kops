package convert

import (
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func flattenTime(d v1.Time) string {
	return d.String()
}

func flattenDuration(d v1.Duration) string {
	return d.String()
}

func flattenQuantity(d resource.Quantity) int64 {
	var q int64

	q = d.Value()

	return q
}
