// Code generated by engine.go; DO NOT EDIT.

package provider

import (
	"k8s.io/kops/pkg/apis/kops"
)

func expandSSHCredentialSpec(in interface{}) (kops.SSHCredentialSpec, bool) {
	d := in.([]interface{})
	out := kops.SSHCredentialSpec{}

	if len(d) < 1 {
		return out, true
	}

	m := d[0].(map[string]interface{})

	if v, ok := m["public_key"]; ok {
		if value, e := expandString(v); !e {
      out.PublicKey = value
    }
	}

  if isEmpty(out) {
    return out, true
  }

	return out, false
}
