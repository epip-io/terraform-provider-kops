package provider

import (
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"reflect"

	"time"
)

func expandTime(d interface{}) (v1.Time, bool) {
	var t v1.Time

	c := d.([]interface{})

	if len(c) > 0 {
		stamp, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", c[0].(string))

		t.Time = stamp

		return t, false
	}

	return t, true

}

func expandDuration(in interface{}) (v1.Duration, bool) {
	t := v1.Duration{}

	d := in.([]interface{})

	if len(d) > 0 {
		pd, _ := time.ParseDuration(d[0].(string))

		t.Duration = pd

		return t, false
	}

	return t, true
}

func expandQuantity(in interface{}) (resource.Quantity, bool) {
	q := resource.Quantity{}

	d := in.([]interface{})

	if len(d) > 0 {
		q.Set(d[0].(int64))

		return q, false
	}

	return q, true
}

func expandInt64(d interface{}) (int64, bool) {
	out := int64(d.(float64))

	if out == 0 {
		return 0, true
	}

	return out, false
}

func expandBool(in interface{}) (bool, bool) {
	out := in.(bool)

	log.Printf("Boolean: %v\n", out)
	return out, !out
}

func expandInt32(in interface{}) (int32, bool) {
	d := in.(int)
	out := int32(d)

	if out == 0 {
		return out, true
	}

	return out, false
}

func expandFloat32(in interface{}) (float32, bool) {
	d := in.(float64)
	out := float32(d)

	if out == 0 {
		return out, true
	}

	return out, false
}

func isEmpty(in interface{}) bool {
	if reflect.ValueOf(in).IsZero() {
		return true
	}

	return false
}
