package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/epip-io/terraform-provider-kops/pkg/provider"
)

//go:generate go mod vendor
//go:generate cp ./hack/bindata.go ./vendor/k8s.io/kops/upup/models/
//go:generate go run ./hack/gen-go.go -mod vendor

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
