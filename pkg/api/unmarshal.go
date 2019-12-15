package api

import (
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func unmarshalTime(d v1.Time) string {
	return d.String()
}

func unmarshalDuration(d v1.Duration) string {
	return d.String()
}

func unmarshalQuantity(d resource.Quantity) int64 {
	var q int64

	q = d.Value()

	return q
}
