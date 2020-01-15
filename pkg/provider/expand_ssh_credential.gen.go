// Code generated by engine.go; DO NOT EDIT.

package provider

import (
	"k8s.io/kops/pkg/apis/kops"
)

func expandSSHCredentialSpec(in interface{}) kops.SSHCredentialSpec {
	d := in.([]interface{})[0].(map[string]interface{})
	out := kops.SSHCredentialSpec{}

	if v, ok := d["public_key"]; ok {
		out.PublicKey = expandString(v)
	}

	return out
}
