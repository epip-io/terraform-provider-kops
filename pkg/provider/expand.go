package provider

import (
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"time"
)

func expandTime(d interface{}) v1.Time {
	var t v1.Time

	c := d.(string)

	stamp, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", c)

	t.Time = stamp

	return t
}

func expandDuration(d interface{}) *v1.Duration {
	t := v1.Duration{}

	pd, _ := time.ParseDuration(d.(string))

	t.Duration = pd

	return &t
}

func expandQuantity(d interface{}) *resource.Quantity {
	q := resource.Quantity{}

	q.Set(d.(int64))

	return &q
}
