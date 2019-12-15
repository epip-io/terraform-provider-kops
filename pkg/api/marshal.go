package api

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"time"
)

func marshalTime(d []interface{}) v1.Time {
	if len(d) > 0 {
		var t v1.Time

		c := d[0].(string)

		stamp, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", c)

		t.Time = stamp
		
		return t
	}

	return v1.Time{}
}

func marshalDuration(d []interface{}) v1.Duration {
	if len(d) > 0 {
		var t v1.Duration

		pd, _ := time.ParseDuration(d[0].(string))
		
		t.Duration = pd

		return t
	}

	return v1.Duration{}
}

func marshalQuantity(d []interface{}) resource.Quantity {
	q := resource.Quantity{}

	if (len(d) > 0) {
		q.Set(d[0].(int64))
	}
	
	return q
}